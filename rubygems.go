package main

import "fmt"

type RubyGem struct {
	Name             string `json:"name"`
	Version          string `json:"version"`
	Downloads        int    `json:"downloads"`
	HomepageURI      string `json:"homepage_uri"`
	DocumentationURI string `json:"documentation_uri"`
}

func rubygem(gem string) chan *RubyGem {
	rubyGemChan := make(chan *RubyGem, 1)

	go func() {
		if gem == "" {
			rubyGemChan <- nil
			return
		}
		var info RubyGem
		err := get(fmt.Sprintf("https://rubygems.org/api/v1/gems/%s.json", gem), &info)
		if err != nil {
			rubyGemChan <- nil
			panic(fmt.Errorf("error fetching rubygems info for %s: %v", gem, err))
		}
		rubyGemChan <- &info
	}()

	return rubyGemChan
}
