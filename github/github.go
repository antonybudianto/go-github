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

type GithubData struct {
	StarCount       int
	RepoCount       int
	ForkCount       int
	WatcherCount    int
	SubscriberCount int
	LanguageMap     map[string]int
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
func FetchAllRepos(username string) (*GithubData, error) {
	page := 1
	starCount := 0
	repoCount := 0
	forkCount := 0
	watcherCount := 0
	subscriberCount := 0
	langMap := make(map[string]int)

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
			return &GithubData{
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
