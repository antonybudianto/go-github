package config

import (
	"log"
	"os"
)

// GithubAccessToken get GH_ACCESS_TOKEN from os env
// if not found it will panic
func GithubAccessToken() string {
	if val, ok := os.LookupEnv("GH_ACCESS_TOKEN"); ok {
		return val
	}

	log.Fatal("GH_ACCESS_TOKEN not set")
	return ""
}
