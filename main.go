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
}

// ProfilePayload for profile response payload
type ProfilePayload struct {
	Username  string        `json:"username"`
	Star      int           `json:"starcount"`
	RepoCount int           `json:"repocount"`
	Repos     []interface{} `json:"repos"`
}

func handleGithubProfile(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	urlPath = strings.TrimPrefix(urlPath, "/")
	urlSeg := strings.SplitN(urlPath, "/", 3)
	username := urlSeg[1]
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?page=1&per_page=100", username)
	fmt.Println(username)

	resp, err := http.Get(url)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		payload := ResponsePayload{
			Error: "Bad request: Error request",
		}
		b, _ := json.Marshal(payload)
		w.Write(b)
		return
	}

	defer resp.Body.Close()

	var data []RepoPayload
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("ERR2", err)
		payload := ResponsePayload{
			Error: "Bad request: Error decoding",
		}
		b, _ := json.Marshal(payload)
		w.Write(b)
		return
	}

	starCount := 0

	for i := 0; i < len(data); i++ {
		starCount += data[i].StarCount
	}

	profilePayload := ProfilePayload{
		Username:  username,
		Star:      starCount,
		RepoCount: len(data),
		Repos:     []interface{}{data},
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
