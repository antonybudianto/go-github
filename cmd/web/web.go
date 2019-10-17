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
var cacheTopStar []byte
var lastCacheTopStar time.Time

const (
	cacheHours        = 24
	cacheHoursTopStar = 24 * 14
)

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

func checkCache(lastCache time.Time, cacheHours int, cacheBytes []byte) bool {
	hoursElapsed := time.Since(lastCache).Hours()
	return hoursElapsed < cacheHoursTopStar && len(cacheTopStar) != 0
}

func handleGithubSummary(w http.ResponseWriter, r *http.Request) {
	if checkCache(lastCache, cacheHours, cacheSummary) {
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

func handleTopStars(w http.ResponseWriter, r *http.Request) {
	if checkCache(lastCacheTopStar, cacheHoursTopStar, cacheTopStar) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Gogithub-Cache", "true")
		w.Write(cacheTopStar)
		return
	}
	data, err := github.FetchAllStars()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		payload := model.ResponsePayload{
			Error: "Error fetch star",
		}
		b, _ := json.Marshal(payload)
		w.WriteHeader(400)
		w.Write(b)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(data)
	cacheTopStar = b
	lastCacheTopStar = time.Now()
	w.Write(b)
	return
}

func main() {
	http.HandleFunc("/gh/summary", handleGithubSummary)
	http.HandleFunc("/gh/profile/", handleGithubProfile)
	http.HandleFunc("/gh/topstars", handleTopStars)

	// For testing purpose
	// http.HandleFunc("/gh/test", handleTest)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
