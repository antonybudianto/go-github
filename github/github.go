package github

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

const ghGqlURL = "https://api.github.com/graphql"

// RepoPayload = response payload from github
type RepoPayload struct {
	StarCount       int       `json:"stargazers_count"`
	ForkCount       int       `json:"forks_count"`
	WatcherCount    int       `json:"watchers_count"`
	SubscriberCount int       `json:"subscribers_count"`
	NetworkCount    int       `json:"network_count"`
	FullName        string    `json:"full_name"`
	Description     string    `json:"description"`
	Language        string    `json:"language"`
	Owner           RepoOwner `json:"owner"`
}

// RepoOwner = owner of the repo
type RepoOwner struct {
	AvatarURL string `json:"avatar_url"`
}

// RepoData = generated summary from raw data
type RepoData struct {
	StarCount   int
	RepoCount   int
	ForkCount   int
	LanguageMap map[string]int32
	AvatarURL   string
}

// UserRepositoryEdge = user's single repo
type UserRepositoryEdge struct {
	Node struct {
		Name            string `json:"name"`
		ForkCount       int    `json:"forkCount"`
		PrimaryLanguage *struct {
			Name string `json:"name"`
		} `json:"primaryLanguage"`
		Stargazers struct {
			TotalCount int `json:"totalCount"`
		} `json:"stargazers"`
	} `json:"node"`
}

// UserRepositoryResponse = response from user gql for repo
type UserRepositoryResponse struct {
	Data struct {
		User struct {
			AvatarURL    string `json:"avatarUrl"`
			Repositories struct {
				TotalCount int `json:"totalCount"`
				PageInfo   struct {
					EndCursor   string `json:"endCursor"`
					HasNextPage bool   `json:"hasNextPage"`
				} `json:"pageInfo"`
				Edges []UserRepositoryEdge `json:"edges"`
			} `json:"repositories"`
		} `json:"user"`
	} `json:"data"`
}

// FetchGhGql = generic fetch for github gql
func FetchGhGql(query, variables string) (map[string]interface{}, error) {
	body, err := json.Marshal(map[string]string{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", ghGqlURL, bytes.NewBuffer(body))

	token := os.Getenv("GH_ACCESS_TOKEN")

	req.Header.Set("Authorization", "bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// FetchTopUserSummary = fetch all top user using GQL
func FetchTopUserSummary() (map[string]interface{}, error) {
	return FetchGhGql(SummaryQuery, "")
}

// FetchRepo = fetch repo by username
func FetchRepo(username string, after *string) (*UserRepositoryResponse, error) {
	variables, _ := json.Marshal(map[string]interface{}{
		"username": username,
		"after":    after,
	})
	data, err := FetchGhGql(UserQuery, string(variables))
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(data)
	var resp UserRepositoryResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// FetchAllRepos = fetch all repos by username and create their summary
func FetchAllRepos(username string) (*RepoData, error) {
	avatarURL := ""
	starCount := 0
	repoCount := 0
	forkCount := 0
	langMap := make(map[string]int32)
	var cursor *string

	for {
		data, err := FetchRepo(username, cursor)
		if err != nil {
			return nil, err
		}

		avatarURL = data.Data.User.AvatarURL
		repoCount = data.Data.User.Repositories.TotalCount

		for i := 0; i < len(data.Data.User.Repositories.Edges); i++ {
			edge := data.Data.User.Repositories.Edges[i]
			starCount += edge.Node.Stargazers.TotalCount
			forkCount += edge.Node.ForkCount
			if edge.Node.PrimaryLanguage != nil {
				langMap[edge.Node.PrimaryLanguage.Name]++
			} else {
				langMap["Others"]++
			}
		}

		if data.Data.User.Repositories.PageInfo.HasNextPage {
			*cursor = data.Data.User.Repositories.PageInfo.EndCursor
		} else {
			break
		}
	}

	return &RepoData{
		AvatarURL:   avatarURL,
		StarCount:   starCount,
		RepoCount:   repoCount,
		ForkCount:   forkCount,
		LanguageMap: langMap,
	}, nil
}
