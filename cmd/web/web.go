package main

import (
	"encoding/json"
	"fmt"
	"gogithub/github"
	"gogithub/model"
	"net/http"
	"strings"
	"time"
)

var cacheSummary []byte
var lastCache time.Time

const cacheHours = 24

// ProfilePayload for profile response payload
type ProfilePayload struct {
	Username      string                     `json:"username"`
	StarCount     int                        `json:"star_count"`
	RepoCount     int                        `json:"repo_count"`
	ForkCount     int                        `json:"fork_count"`
	LanguageCount int                        `json:"language_count"`
	LanguageMap   map[string]int32           `json:"language_map"`
	AvatarURL     string                     `json:"avatar_url"`
	TopRepo       *github.UserRepositoryEdge `json:"top_repo"`
}

func handleGithubProfile(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	urlPath = strings.TrimPrefix(urlPath, "/")
	urlSeg := strings.SplitN(urlPath, "/", 3)
	username := urlSeg[2]

	data, err := github.FetchAllRepos(username)

	if err != nil {
		fmt.Println("ERR handleGitHubProfile:", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		payload := model.ResponsePayload{
			Error: "Error fetch profile",
		}
		b, _ := json.Marshal(payload)
		w.Write(b)
		return
	}

	profilePayload := ProfilePayload{
		Username:      username,
		AvatarURL:     data.AvatarURL,
		StarCount:     data.StarCount,
		RepoCount:     data.RepoCount,
		ForkCount:     data.ForkCount,
		LanguageCount: len(data.LanguageMap),
		LanguageMap:   data.LanguageMap,
		TopRepo:       data.TopRepo,
	}

	payload := model.ResponsePayload{
		Data: profilePayload,
	}

	b, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func handleGithubSummary(w http.ResponseWriter, r *http.Request) {
	hoursElapsed := time.Since(lastCache).Hours()
	if hoursElapsed < cacheHours && len(cacheSummary) != 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Gogithub-Cache", "true")
		w.Write(cacheSummary)
		return
	}
	data, err := github.FetchTopUserSummary()
	if err != nil {
		fmt.Println("ERR", err)
		w.WriteHeader(http.StatusBadRequest)
		payload := model.ResponsePayload{
			Error: "Error fetch summary",
		}
		b, _ := json.Marshal(payload)
		w.Write(b)
		return
	}
	b, _ := json.Marshal(data)
	cacheSummary = b
	lastCache = time.Now()
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	data, _ := github.FetchAllRepos("antonybudianto")
	b, _ := json.Marshal(data)
	w.Write(b)
}

func main() {
	http.HandleFunc("/gh/summary", handleGithubSummary)
	http.HandleFunc("/gh/profile/", handleGithubProfile)

	// For testing purpose
	http.HandleFunc("/gh/test", handleTest)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
