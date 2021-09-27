package dashboard

import (
	"log"
	"sync"
	"time"
)

var defaultProjectMap = sync.Map{}

func init() {
	for _, p := range defaultProjects {
		defaultProjectMap.Store(p.Name, p)
	}
	go resetProjectsPeriodically()
	go prefillAllProjectsFromGitHub()
}

func resetProjectsPeriodically() {
	for range time.Tick(time.Hour / 2) {
		log.Println("resetting projects' cache")
		resetProjects()
		prefillAllProjectsFromGitHub()
	}
}

func resetProjects() {
	for _, p := range defaultProjects {
		p.reset()
	}
}

type Project struct {
	GlobalRelayID string `json:"id"`
	Name          string `json:"name"`
	Nwo           string `json:"nwo"`
	Branch        string `json:"branch"`
	GemName       string `json:"gem_name"`
	Stars         int    `json:"star_count"`

	Gem     *RubyGem `json:"gem"`
	GitHub  *GitHub  `json:"github"`
	fetched bool
}

func (p *Project) fetch() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		p.fetchGitHubData()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		p.fetchRubyGemData()
		wg.Done()
	}()

	wg.Wait()

	p.fetched = true
}

func (p *Project) fetchRubyGemData() {
	if p.Gem != nil {
		return
	}

	p.Gem = <-rubygem(p.GemName)
}

func (p *Project) fetchGitHubData() {
	if p.GitHub != nil {
		return
	}

	p.GitHub = <-github(p.GlobalRelayID)
}

func (p *Project) reset() {
	p.fetched = false
	p.Gem = nil
	p.GitHub = nil
}

func getProject(name string) *Project {
	if p, ok := defaultProjectMap.Load(name); ok {
		proj := p.(*Project)
		if !proj.fetched {
			proj.fetch()
		}
		return proj
	}

	return nil
}

func getProjects() []*Project {
	return defaultProjects
}
