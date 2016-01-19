package dashboard

import (
	"fmt"
	"sync"
)

var (
	defaultProjectMap map[string]*Project
	defaultProjects   = []*Project{
		makeProject("jekyll", "jekyll/jekyll", "master", "jekyll"),
		makeProject("jemoji", "jekyll/jemoji", "master", "jemoji"),
		makeProject("mercenary", "jekyll/mercenary", "master", "mercenary"),
		makeProject("jekyll-import", "jekyll/jekyll-import", "master", "jekyll-import"),
		makeProject("jekyll-feed", "jekyll/jekyll-feed", "master", "jekyll-feed"),
		makeProject("jekyll-sitemap", "jekyll/jekyll-sitemap", "master", "jekyll-sitemap"),
		makeProject("jekyll-mentions", "jekyll/jekyll-mentions", "master", "jekyll-mentions"),
		makeProject("jekyll-watch", "jekyll/jekyll-watch", "master", "jekyll-watch"),
		makeProject("jekyll-compose", "jekyll/jekyll-compose", "master", "jekyll-compose"),
		makeProject("jekyll-paginate", "jekyll/jekyll-paginate", "master", "jekyll-paginate"),
		makeProject("jekyll-gist", "jekyll/jekyll-gist", "master", "jekyll-gist"),
		makeProject("jekyll-coffeescript", "jekyll/jekyll-coffeescript", "master", "jekyll-coffeescript"),
		makeProject("jekyll-opal", "jekyll/jekyll-opal", "master", "jekyll-opal"),
		makeProject("classifier-reborn", "jekyll/classifier-reborn", "master", "classifier-reborn"),
		makeProject("jekyll-sass-converter", "jekyll/jekyll-sass-converter", "master", "jekyll-sass-converter"),
		makeProject("jekyll-textile-converter", "jekyll/jekyll-textile-converter", "master", "jekyll-textile-converter"),
		makeProject("plugins.jekyllrb", "jekyll/plugins", "gh-pages", ""),
		makeProject("jekyll-redirect-from", "jekyll/jekyll-redirect-from", "master", "jekyll-redirect-from"),
		makeProject("github-metadata", "jekyll/github-metadata", "master", "jekyll-github-metadata"),
		makeProject("jekyll docker", "jekyll/docker", "", ""),
	}
)

type Project struct {
	Name    string `json:"name"`
	Nwo     string `json:"nwo"`
	Branch  string `json:"branch"`
	GemName string `json:"gem_name"`

	Gem     *RubyGem      `json:"gem"`
	Travis  *TravisReport `json:"travis"`
	fetched bool
}

func (p *Project) fetch() {
	if !p.fetched {
		rubyGemChan := rubygem(p.GemName)
		travisChan := travis(p.Nwo, p.Branch)
		p.Gem = <-rubyGemChan
		p.Travis = <-travisChan
		p.fetched = true
	}
}

func buildProjectMap() {
	defaultProjectMap = map[string]*Project{}
	for _, p := range defaultProjects {
		defaultProjectMap[p.Name] = p
	}
}

func makeProject(name, nwo, branch, rubygem string) *Project {
	return &Project{
		Name:    name,
		Nwo:     nwo,
		Branch:  branch,
		GemName: rubygem,
	}
}

func getProject(name string) Project {
	if defaultProjectMap == nil {
		buildProjectMap()
	}

	if p, ok := defaultProjectMap[name]; ok {
		if !p.fetched {
			p.fetch()
		}
		return *p
	}
	panic(fmt.Sprintf("no project named '%s'", name))
}

func getAllProjects() []*Project {
	var wg sync.WaitGroup
	for _, p := range defaultProjects {
		wg.Add(1)
		go func(project *Project) {
			project.fetch()
			wg.Done()
		}(p)
	}
	wg.Wait()
	return defaultProjects
}
