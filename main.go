package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ResponsePayload is the basic response payload
type ResponsePayload struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

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

// ProfilePayload for profile response payload
type ProfilePayload struct {
	Username        string         `json:"username"`
	StarCount       int            `json:"star_count"`
	RepoCount       int            `json:"repo_count"`
	ForkCount       int            `json:"fork_count"`
	WatcherCount    int            `json:"watcher_count"`
	SubscriberCount int            `json:"subscriber_count"`
	LanguageMap     map[string]int `json:"language_map"`
}

func fetchRepo(w http.ResponseWriter, r *http.Request, username string, page int) ([]RepoPayload, error) {
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

func handleGithubProfile(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	urlPath = strings.TrimPrefix(urlPath, "/")
	urlSeg := strings.SplitN(urlPath, "/", 3)
	username := urlSeg[1]
	page := 1
	starCount := 0
	repoCount := 0
	forkCount := 0
	watcherCount := 0
	subscriberCount := 0
	langMap := make(map[string]int)

	for {
		repos, err := fetchRepo(w, r, username, page)
		if err != nil {
			fmt.Println("ERR", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			payload := ResponsePayload{
				Error: "Error fetch repo",
			}
			b, _ := json.Marshal(payload)
			w.Write(b)
			return
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
			break
		} else {
			page++
		}
	}

	profilePayload := ProfilePayload{
		Username:        username,
		StarCount:       starCount,
		RepoCount:       repoCount,
		ForkCount:       forkCount,
		WatcherCount:    watcherCount,
		SubscriberCount: subscriberCount,
		LanguageMap:     langMap,
	}

	payload := ResponsePayload{
		Data: profilePayload,
	}

	b, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func main() {
	http.HandleFunc("/profile/", handleGithubProfile)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
