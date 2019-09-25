package github

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

const ghGqlURL = "https://api.github.com/graphql"

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

type UserRepositoryEdge struct {
	Node struct {
		Name            string `json:"name"`
		ForkCount       int    `json:"forkCount"`
		PrimaryLanguage struct {
			Name string `json:"name"`
		} `json:"primaryLanguage"`
		Stargazers struct {
			TotalCount int `json:"totalCount"`
		} `json:"stargazers"`
	} `json:"node"`
}

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

	token := os.Getenv("GH_ACCESS_TOKEN")

	req.Header.Set("Authorization", "bearer "+token)
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
	query := `
	query topSummary {
		topPHPDev: search(query: "location:Indonesia language:PHP followers:>=200", type: USER, first: 10) {
		  edges {
			node {
			  ... on User {
				name
				avatarUrl
				login
				bio
				company
				location
				following {
				  totalCount
				}
				followers {
				  totalCount
				}
			  }
			}
		  }
		}
		topJsDev: search(query: "location:Indonesia language:JavaScript followers:>=200", type: USER, first: 10) {
		  edges {
			node {
			  ... on User {
				name
				avatarUrl
				login
				bio
				company
				location
				following {
				  totalCount
				}
				followers {
				  totalCount
				}
			  }
			}
		  }
		}
		
		topJavaDev: search(query: "location:Indonesia language:Java followers:>=200", type: USER, first: 10) {
		  edges {
			node {
			  ... on User {
				name
				avatarUrl
				login
				bio
				company
				location
				following {
				  totalCount
				}
				followers {
				  totalCount
				}
			  }
			}
		  }
		}
		
		topPythonDev: search(query: "location:Indonesia language:Python followers:>=150", type: USER, first: 10) {
		  edges {
			node {
			  ... on User {
				name
				avatarUrl
				login
				bio
				company
				location
				following {
				  totalCount
				}
				followers {
				  totalCount
				}
			  }
			}
		  }
		}
		
		topGoDev: search(query: "location:Indonesia language:Go followers:>=100", type: USER, first: 10) {
		  edges {
			node {
			  ... on User {
				name
				avatarUrl
				login
				bio
				company
				location
				following {
				  totalCount
				}
				followers {
				  totalCount
				}
			  }
			}
		  }
		}
	  
		topJakartaDev: search(query: "location:Jakarta followers:>=300", type: USER, first: 10) {
		  edges {
			node {
			  ... on User {
				name
				avatarUrl
				login
				bio
				company
				location
				following {
				  totalCount
				}
				followers {
				  totalCount
				}
			  }
			}
		  }
		}
		
		topBandungDev: search(query: "location:Bandung followers:>=200", type: USER, first: 10) {
		  edges {
			node {
			  ... on User {
				name
				avatarUrl
				login
				bio
				company
				location
				following {
				  totalCount
				}
				followers {
				  totalCount
				}
			  }
			}
		  }
		}
		
		topYogyakartaDev: search(query: "location:Yogyakarta followers:>=100", type: USER, first: 10) {
		  edges {
			node {
			  ... on User {
				name
				avatarUrl
				login
				bio
				company
				location
				following {
				  totalCount
				}
				followers {
				  totalCount
				}
			  }
			}
		  }
		}
		
		topMalangDev: search(query: "location:Malang followers:>=100", type: USER, first: 10) {
		  edges {
			node {
			  ... on User {
				name
				avatarUrl
				login
				bio
				company
				location
				following {
				  totalCount
				}
				followers {
				  totalCount
				}
			  }
			}
		  }
		}

		topSurabayaDev: search(query: "location:Surabaya followers:>=100", type: USER, first: 10) {
		  edges {
			node {
			... on User {
				name
				avatarUrl
				login
				bio
				company
				location
				following {
				totalCount
				}
				followers {
				totalCount
				}
			  }
			}
		  }
		}

	  }
	`

	return FetchGhGql(query, "")
}

// FetchRepo = fetch repo by username
func FetchRepo(username string, after *string) (*UserRepositoryResponse, error) {
	query := `
	query getUserRepo($username: String!, $after: String) {
		user(login:$username){
		  avatarUrl
		  repositories(after:$after, first:100, ownerAffiliations:OWNER, isFork:false, privacy:PUBLIC){
			totalCount
			pageInfo{
			  endCursor
			  hasNextPage
			}
			edges{
			  node{
				name
				forkCount
				stargazers {
				  totalCount
				}
			  }
			}
		  }
		}
	}
	`
	variables, _ := json.Marshal(map[string]interface{}{
		"username": username,
		"after":    after,
	})
	data, err := FetchGhGql(query, string(variables))
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

// FetchAllRepos = fetch all repos by username
func FetchAllRepos(username string) (*RepoData, error) {
	starCount := 0
	repoCount := 0
	forkCount := 0
	watcherCount := 0
	subscriberCount := 0
	// langMap := make(map[string]int32)
	avatarUrl := ""
	var cursor *string

	for {
		data, err := FetchRepo(username, cursor)
		if err != nil {
			break
		}

		avatarUrl = data.Data.User.AvatarURL
		repoCount = data.Data.User.Repositories.TotalCount

		for i := 0; i < len(data.Data.User.Repositories.Edges); i++ {
			edge := data.Data.User.Repositories.Edges[i]
			starCount += edge.Node.Stargazers.TotalCount
			forkCount += edge.Node.ForkCount
		}

		if data.Data.User.Repositories.PageInfo.HasNextPage {
			*cursor = data.Data.User.Repositories.PageInfo.EndCursor
		} else {
			break
		}
	}

	return &RepoData{
		AvatarURL:       avatarUrl,
		StarCount:       starCount,
		RepoCount:       repoCount,
		ForkCount:       forkCount,
		WatcherCount:    watcherCount,
		SubscriberCount: subscriberCount,
		LanguageMap:     nil,
	}, nil
}
