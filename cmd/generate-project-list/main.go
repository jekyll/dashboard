package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/v73/github"
	"github.com/jekyll/dashboard"
)

var tmpl = template.Must(template.New(`configuration.go`).Parse(`//go:generate go run ./cmd/generate-project-list
// THIS FILE IS AUTO-GENERATED WITH 'go generate .'
// LAST UPDATED {{.Now.Format "Jan 02, 2006 15:04:05 UTC" }}
package dashboard

var defaultProjects = []*Project{{"{"}}{{range .Repositories}}
	{
		Name:          {{.Name | printf "%q"}},
		Nwo:           {{.Nwo | printf "%q"}},
		Branch:        {{.Branch | printf "%q"}},
		GemName:       {{.GemName | printf "%q"}},
		GlobalRelayID: {{.GlobalRelayID | printf "%q" }},
	},{{end}}
}
`))

type templateData struct {
	Repositories []*dashboard.Project
	Now          time.Time
}

var additionalProjectNames = map[string]bool{
	"classifier-reborn": true,
	"directory":         true,
	"github-metadata":   true,
	"jekyll":            true,
	"jemoji":            true,
	"mercenary":         true,
	"minima":            true,
}

func relevantProject(name string) bool {
	return strings.HasPrefix(name, "jekyll-") && !strings.HasPrefix(name, "jekyll-test") || additionalProjectNames[name]
}

func projectGemName(name string) string {
	switch name {
	case "github-metadata":
		return "jekyll-github-metadata"
	case "directory":
		return ""
	default:
		return name
	}
}

func main() {
	flag.Parse()

	client := github.NewClient(http.DefaultClient)
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	repositories, _, err := client.Repositories.ListByOrg(context.Background(), "jekyll", opt)
	if err != nil {
		log.Fatalf("unable to list repositories: %v", err)
	}

	repoInfos := make([]*dashboard.Project, 0, len(repositories))

	for _, repository := range repositories {
		name := repository.GetName()
		if !relevantProject(name) || repository.GetArchived() || repository.GetPrivate() {
			continue
		}
		info := &dashboard.Project{
			Name:          repository.GetName(),
			Nwo:           repository.GetFullName(),
			Branch:        repository.GetDefaultBranch(),
			GemName:       projectGemName(repository.GetName()),
			GlobalRelayID: repository.GetNodeID(),
			Stars:         repository.GetStargazersCount(),
		}
		repoInfos = append(repoInfos, info)
		log.Printf("repo: %#v", info)
	}

	sort.Slice(repoInfos, func(i, j int) bool {
		return repoInfos[i].Stars > repoInfos[j].Stars
	})

	s := &strings.Builder{}
	err = tmpl.Execute(s, templateData{Repositories: repoInfos, Now: time.Now().UTC()})
	if err != nil {
		log.Fatalf("unable to execute template: %v", err)
	}

	path := filepath.Join("configuration.go")
	err = ioutil.WriteFile(path, []byte(s.String()), 0644)
	if err != nil {
		log.Fatalf("unable to write %s: %v", path, err)
	}
}
