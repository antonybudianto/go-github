package main

import (
	"encoding/json"
	"fmt"
	"gogithub/github"
	"gogithub/model"
	"net/http"
	"strings"
)

// ProfilePayload for profile response payload
type ProfilePayload struct {
	Username        string         `json:"username"`
	StarCount       int            `json:"star_count"`
	RepoCount       int            `json:"repo_count"`
	ForkCount       int            `json:"fork_count"`
	WatcherCount    int            `json:"watcher_count"`
	SubscriberCount int            `json:"subscriber_count"`
	LanguageCount   int            `json:"language_count"`
	LanguageMap     map[string]int `json:"language_map"`
}

func handleGithubProfile(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	urlPath = strings.TrimPrefix(urlPath, "/")
	urlSeg := strings.SplitN(urlPath, "/", 3)
	username := urlSeg[1]

	data, err := github.FetchAllRepos(username)

	if err != nil {
		fmt.Println("ERR", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		payload := model.ResponsePayload{
			Error: "Error fetch repo",
		}
		b, _ := json.Marshal(payload)
		w.Write(b)
		return
	}

	profilePayload := ProfilePayload{
		Username:        username,
		StarCount:       data.StarCount,
		RepoCount:       data.RepoCount,
		ForkCount:       data.ForkCount,
		WatcherCount:    data.WatcherCount,
		SubscriberCount: data.SubscriberCount,
		LanguageCount:   len(data.LanguageMap),
		LanguageMap:     data.LanguageMap,
	}

	payload := model.ResponsePayload{
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
