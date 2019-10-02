package github

import "fmt"

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

func generateSummaryQuery(name, location, language, followers string) string {
	return fmt.Sprintf(`
	%s: search(query: "location:%s language:%s followers:%s", type: USER, first: 10) {
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
	`, name, location, language, followers)
}

// SummaryQuery = query used when fetch all summary
var SummaryQuery = fmt.Sprintf(`
query topSummary {
	%s
	%s
	%s
	%s
	%s
	%s
	%s
	%s
	%s
	%s
	%s
	%s
	%s
	%s
	%s
	%s
  }
`,
	generateSummaryQuery("topPHPDev", "Indonesia", "PHP", ">=200"),
	generateSummaryQuery("topJsDev", "Indonesia", "JavaScript", ">=200"),
	generateSummaryQuery("topJavaDev", "Indonesia", "Java", ">=200"),
	generateSummaryQuery("topPythonDev", "Indonesia", "Python", ">=150"),
	generateSummaryQuery("topHTMLDev", "Indonesia", "HTML", ">=150"),
	generateSummaryQuery("topGoDev", "Indonesia", "Go", ">=100"),
	generateSummaryQuery("topRubyDev", "Indonesia", "Ruby", ">=100"),
	generateSummaryQuery("topShellDev", "Indonesia", "Shell", ">=100"),
	generateSummaryQuery("topSwiftDev", "Indonesia", "Swift", ">=50"),

	generateSummaryQuery("topJakartaDev", "Jakarta", "*", ">=300"),
	generateSummaryQuery("topBandungDev", "Bandung", "*", ">=200"),
	generateSummaryQuery("topYogyakartaDev", "Yogyakarta", "*", ">=100"),
	generateSummaryQuery("topMalangDev", "Malang", "*", ">=100"),
	generateSummaryQuery("topBaliDev", "Bali", "*", ">=100"),
	generateSummaryQuery("topSurabayaDev", "Surabaya", "*", ">=100"),
	generateSummaryQuery("topSemarangDev", "Semarang", "*", ">=100"),
)
