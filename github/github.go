package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gogithub/config"
	"gogithub/model"
	"net/http"

	"github.com/joho/godotenv"
)

const ghGqlURL = "https://api.github.com/graphql"

// RepoData = generated summary from raw data
type RepoData struct {
	StarCount   int
	RepoCount   int
	ForkCount   int
	AvatarURL   string
	LanguageMap map[string]int32
	TopRepo     *UserRepositoryEdge
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

func init() {
	// try load from .env file
	_ = godotenv.Load()
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

	req.Header.Set("Authorization", "bearer "+config.GithubAccessToken())
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
	var bestRepo *UserRepositoryEdge

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

			if bestRepo == nil {
				bestRepo = &edge
			} else if edge.Node.Stargazers.TotalCount > bestRepo.Node.Stargazers.TotalCount {
				bestRepo = &edge
			}

			if edge.Node.PrimaryLanguage != nil {
				langMap[edge.Node.PrimaryLanguage.Name]++
			} else {
				langMap["Others"]++
			}
		}

		if data.Data.User.Repositories.PageInfo.HasNextPage {
			cursor = &data.Data.User.Repositories.PageInfo.EndCursor
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
		TopRepo:     bestRepo,
	}, nil
}

// SummaryDev - summary dev for star fetch purpose
type SummaryDev struct {
	Node struct {
		Login string `json:"login"`
	} `json:"node"`
}

// SummaryData - extract cache for fetching star
type SummaryData struct {
	Data struct {
		TopJakartaDev struct {
			Edges []SummaryDev `json:"edges"`
		} `json:"topJakartaDev"`
	} `json:"data"`
}

// DevStar - for single dev star data
type DevStar struct {
	Login string `json:"login"`
	Stars int    `json:"stars"`
}

// DevChannel - custom dev channel for async
type DevChannel struct {
	Data     *RepoData
	Username string
}

func asyncFetchRepos(ch chan DevChannel, username string) {
	fmt.Println("fetching...", username)
	devData, err := FetchAllRepos(username)
	if err != nil {
		ch <- DevChannel{
			Username: username,
		}
		return
	}
	fmt.Println("done...", username)
	ch <- DevChannel{
		Data:     devData,
		Username: username,
	}
}

// FetchAllStars - fetch all dev star from cache
func FetchAllStars(cache []byte) (*model.ResponsePayload, error) {
	var data SummaryData
	err := json.Unmarshal(cache, &data)
	var devList []SummaryDev
	var devStarList []DevStar
	if err != nil {
		return nil, err
	}
	devList = append(devList, data.Data.TopJakartaDev.Edges...)
	ch := make(chan DevChannel)
	for i := 0; i < len(devList); i++ {
		dev := devList[i]
		go asyncFetchRepos(ch, dev.Node.Login)
	}
	for i := 0; i < len(devList); i++ {
		devStar := DevStar{}
		devData := <-ch
		devStar.Stars = devData.Data.StarCount
		devStar.Login = devData.Username
		devStarList = append(devStarList, devStar)
	}
	return &model.ResponsePayload{
		Data: devStarList,
	}, nil
}
