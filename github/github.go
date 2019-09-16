package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RepoPayload = response payload from github
type RepoPayload struct {
	StarCount       int    `json:"stargazers_count"`
	ForkCount       int    `json:"forks_count"`
	WatcherCount    int    `json:"watchers_count"`
	SubscriberCount int    `json:"subscribers_count"`
	NetworkCount    int    `json:"network_count"`
	FullName        string `json:"full_name"`
	Description     string `json:"description"`
	Language        string `json:"language"`
}

// RepoData = generated summary from raw data
type RepoData struct {
	StarCount       int
	RepoCount       int
	ForkCount       int
	WatcherCount    int
	SubscriberCount int
	LanguageMap     map[string]int32
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
}

// FetchTopUsers = fetch top user by our own custom criteria
func FetchTopUsers(location string, follower string, language string) (*UserPayload, error) {
	url := fmt.Sprintf("https://api.github.com/search/users?q=location:%s+followers:%s+language:%s+type:user", location, follower, language)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data UserPayload
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

// FetchRepo = fetch repo by username
func FetchRepo(username string, page int) ([]RepoPayload, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?page=%d&per_page=100", username, page)

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

	for {
		repos, err := FetchRepo(username, page)
		if err != nil {
			return nil, err
		}
		repoCount += len(repos)

		for i := 0; i < len(repos); i++ {
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
