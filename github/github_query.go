package github

// UserQuery = query used when fetching profile
var UserQuery = `
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
			primaryLanguage {
			  name
			}
			stargazers {
			  totalCount
			}
		  }
		}
	  }
	}
}
`

// SummaryQuery = query used when fetch all summary
var SummaryQuery = `
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
