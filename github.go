package dashboard

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	gh "github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
)

const accessTokenEnvVar = "GITHUB_ACCESS_TOKEN"
const graphqlQuery = `
query jekyllFetchDashboardData($ids: [ID!]!) {
  nodes(ids: $ids) {
    ... on Repository {
      id
      owner {
        login
      }
      name
      pullRequests(states: [OPEN]) {
        totalCount
      }
      issues(states: [OPEN]) {
        totalCount
      }
      releases(first: 5, orderBy: {field: CREATED_AT, direction: DESC}) {
        nodes {
          tag {
            name
            target {
              __typename
              ... on Commit {
                history {
                  totalCount
                }
              }
              ... on Tag {
                target {
                  ... on Commit {
                    history {
                      totalCount
                    }
                  }
                }
              }
            }
          }
          publishedAt
          isPrerelease
        }
      }
      defaultBranchRef {
        target {
          ... on Commit {
            history {
              totalCount
			  nodes {
                statusCheckRollup {
                  contexts(first: 100) {
                    nodes {
                      __typename
                      ... on CheckRun {
                        name
                        status
                        conclusion
                        url
                      }
                      ... on StatusContext {
                        description
                        state
                        targetUrl
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`

type githubGraphQLResults struct {
	once    sync.Once
	fetched bool

	Data struct {
		Nodes []struct {
			GlobalRelayID string `json:"id"`
			Owner         struct {
				Login string `json:"login"`
			} `json:"owner"`
			Name         string `json:"name"`
			PullRequests struct {
				TotalCount int `json:"totalCount"`
			} `json:"pullRequests"`
			Issues struct {
				TotalCount int `json:"totalCount"`
			} `json:"issues"`
			Releases struct {
				Nodes []struct {
					Tag struct {
						Name   string `json:"name"`
						Target struct {
							TypeName string `json:"__typename"`
							Target   struct {
								History struct {
									TotalCount int `json:"totalCount"`
								} `json:"history"`
							} `json:"target"`
							History struct {
								TotalCount int `json:"totalCount"`
							} `json:"history"`
						} `json:"target"`
					} `json:"tag"`
					PublishedAt  time.Time `json:"publishedAt"`
					IsPreRelease bool      `json:"isPrerelease"`
				} `json:"nodes"`
			} `json:"releases"`
			DefaultBranchRef struct {
				Target struct {
					History struct {
						TotalCount int `json:"totalCount"`
						Nodes      []struct {
							StatusCheckRollup struct {
								Contexts struct {
									Nodes []struct {
										TypeName string `json:"__typename"`

										// On CheckRun
										Name       string `json:"name"`
										Status     string `json:"status"`
										Conclusion string `json:"conclusion"`
										URL        string `json:"url"`

										// On StatusContext
										Description string `json:"description"`
										State       string `json:"state"`
										TargetURL   string `json:"targetUrl"`
									} `json:"nodes"`
								} `json:"contexts"`
							} `json:"statusCheckRollup"`
						} `json:"nodes"`
					} `json:"history"`
				} `json:"target"`
			} `json:"defaultBranchRef"`
		} `json:"nodes"`
	} `json:"data"`
}

var githubClient *gh.Client

type GitHub struct {
	Owner                     string            `json:"owner"`
	Name                      string            `json:"name"`
	CommitsThisWeek           int               `json:"commits_this_week"`
	OpenPRs                   int               `json:"open_prs"`
	OpenIssues                int               `json:"open_issues"`
	CommitsSinceLatestRelease int               `json:"commits_since_latest_release"`
	LatestReleaseTag          string            `json:"latest_release_tag"`
	LatestCommitCIData        []githubCIContext `json:"latest_commit_ci_data"`
}

type githubCIContext struct {
	Name     string `json:"name"`
	State    string `json:"state"`
	URL      string `json:"url"`
	TypeName string `json:"__typename"`
}

func gitHubToken() string {
	return os.Getenv(accessTokenEnvVar)
}

func newGitHubClient() *gh.Client {
	if token := gitHubToken(); token != "" {
		return gh.NewClient(oauth2.NewClient(
			oauth2.NoContext,
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			),
		))
	} else {
		log.Printf("%s required for GitHub", accessTokenEnvVar)
		return gh.NewClient(nil)
	}
}

func github(globalRelayID string) chan *GitHub {
	githubChan := make(chan *GitHub, 1)

	go func() {
		if globalRelayID == "" || githubClient == nil {
			githubChan <- nil
			close(githubChan)
			return
		}

		githubChan <- loadGitHubFromGraphQL(globalRelayID)
		close(githubChan)
	}()

	return githubChan
}

func loadGitHubFromGraphQL(globalRelayID string) *GitHub {
	githubGraphQLData := &githubGraphQLResults{}
	githubData := &GitHub{}

	err := doGraphql(githubClient, graphqlQuery, map[string]interface{}{"ids": []string{globalRelayID}}, githubGraphQLData)
	if err != nil {
		log.Printf("error fetching graphql: %+v", err)
	}

	for _, githubProject := range githubGraphQLData.Data.Nodes {
		if githubProject.GlobalRelayID == globalRelayID {
			githubData.Owner = githubProject.Owner.Login
			githubData.Name = githubProject.Name
			githubData.OpenPRs = githubProject.PullRequests.TotalCount
			githubData.OpenIssues = githubProject.Issues.TotalCount
			for _, release := range githubProject.Releases.Nodes {
				if !release.IsPreRelease {
					githubData.LatestReleaseTag = release.Tag.Name
					if release.Tag.Target.TypeName == "Commit" {
						githubData.CommitsSinceLatestRelease = githubProject.DefaultBranchRef.Target.History.TotalCount - release.Tag.Target.History.TotalCount
					} else {
						githubData.CommitsSinceLatestRelease = githubProject.DefaultBranchRef.Target.History.TotalCount - release.Tag.Target.Target.History.TotalCount
					}
					break
				}
			}

			index := 0
			if len(githubProject.DefaultBranchRef.Target.History.Nodes[1].StatusCheckRollup.Contexts.Nodes) > len(githubProject.DefaultBranchRef.Target.History.Nodes[0].StatusCheckRollup.Contexts.Nodes) {
				index = 1
			}
			for _, ciContext := range githubProject.DefaultBranchRef.Target.History.Nodes[index].StatusCheckRollup.Contexts.Nodes {
				var name, state, url string
				if ciContext.TypeName == "CheckRun" {
					name = ciContext.Name
					url = ciContext.URL
					if ciContext.Status == "pending" {
						state = "pending"
					} else {
						state = ciContext.Conclusion
					}
				} else {
					name = ciContext.Description
					state = ciContext.State
					url = ciContext.TargetURL
				}
				githubData.LatestCommitCIData = append(githubData.LatestCommitCIData, githubCIContext{
					TypeName: ciContext.TypeName,
					Name:     name,
					State:    state,
					URL:      url,
				})
			}
			break
		}
	}

	return githubData
}

func prefillAllProjectsFromGitHub() {
	var wg sync.WaitGroup
	for _, project := range getProjects() {
		wg.Add(1)
		project := project
		go func() {
			project.fetchGitHubData()
			wg.Done()
		}()
	}
	wg.Wait()
}

func commitsSinceLatestRelease(owner, repo, latestReleaseTagName string) int {
	var comparison *gh.CommitsComparison
	var err error
	logHTTP("GET", fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/compare/%s...master",
		owner, repo, latestReleaseTagName,
	), func() {
		comparison, _, err = githubClient.Repositories.CompareCommits(
			context.Background(),
			owner, repo,
			latestReleaseTagName, "master",
		)
	})
	if err != nil {
		log.Printf("error fetching commit comparison for %s...master for %s/%s: %v", latestReleaseTagName, owner, repo, err)
		return -1
	}
	return *comparison.TotalCommits
}
