package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

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
	StarCount       int
	RepoCount       int
	ForkCount       int
	WatcherCount    int
	SubscriberCount int
	LanguageMap     map[string]int32
	AvatarURL       string
}

// UserPayload = response payload from github user search
type UserPayload struct {
	TotalCount        int        `json:"total_count"`
	IncompleteResults bool       `json:"incomplete_results"`
	Items             []UserItem `json:"items"`
}

// UserItem = user on userpayload
type UserItem struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Stars     int    `json:"stars"`
}

func genClientQuery() string {
	clientID := os.Getenv("GH_CLIENT_ID")
	clientSecret := os.Getenv("GH_CLIENT_SECRET")
	return fmt.Sprintf("client_id=%s&client_secret=%s", clientID, clientSecret)
}

// FetchTopUsers = fetch top user by our own custom criteria
func FetchTopUsers(location string, follower string, language string) (*UserPayload, error) {
	url := fmt.Sprintf("https://api.github.com/search/users?q=location:%s+followers:%s+language:%s+type:user&"+genClientQuery(), location, follower, language)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data UserPayload
	err = json.NewDecoder(resp.Body).Decode(&data)
	// for i := 0; i < len(data.Items); i++ {
	// 	var stars int
	// 	repoData, err := FetchAllRepos(data.Items[i].Login)
	// 	if err != nil {
	// 		stars = 0
	// 	} else {
	// 		stars = repoData.StarCount
	// 	}
	// 	data.Items[i].Stars = stars
	// }

	if err != nil {
		return nil, err
	}

	return &data, nil
}

// FetchRepo = fetch repo by username
func FetchRepo(username string, page int) ([]RepoPayload, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?page=%d&per_page=100&"+genClientQuery(), username, page)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data []RepoPayload
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// FetchAllRepos = fetch all repos by username
func FetchAllRepos(username string) (*RepoData, error) {
	page := 1
	starCount := 0
	repoCount := 0
	forkCount := 0
	watcherCount := 0
	subscriberCount := 0
	langMap := make(map[string]int32)
	avatarUrl := ""

	for {
		repos, err := FetchRepo(username, page)
		if err != nil {
			return nil, err
		}
		repoCount += len(repos)

		for i := 0; i < len(repos); i++ {
			if i == 0 {
				avatarUrl = repos[i].Owner.AvatarURL
			}
			starCount += repos[i].StarCount
			forkCount += repos[i].ForkCount
			watcherCount += repos[i].WatcherCount
			subscriberCount += repos[i].SubscriberCount
			if repos[i].Language != "" {
				langMap[repos[i].Language]++
			} else {
				langMap["Others"]++
			}
		}

		if len(repos) == 0 {
			return &RepoData{
				AvatarURL:       avatarUrl,
				StarCount:       starCount,
				RepoCount:       repoCount,
				ForkCount:       forkCount,
				WatcherCount:    watcherCount,
				SubscriberCount: subscriberCount,
				LanguageMap:     langMap,
			}, nil
		}

		page++
	}
}
