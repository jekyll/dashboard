package main

import "fmt"

type TravisReport struct {
	Nwo    string       `json:"nwo"`
	Branch TravisBranch `json:"branch"`
}

type TravisBranch struct {
	Id    int    `json:"id"`
	State string `json:"state"`
}

func travis(nwo, branch string) chan *TravisReport {
	travisChan := make(chan *TravisReport, 1)

	go func() {
		if branch == "" {
			travisChan <- nil
			return
		}

		var info TravisReport
		info.Nwo = nwo
		err := get(fmt.Sprintf("https://api.travis-ci.org/repos/%s/branches/%s", nwo, branch), &info)
		if err != nil {
			travisChan <- nil
			panic(fmt.Errorf("error fetching travis info for %s/%s: %v", nwo, branch, err))
		}
		travisChan <- &info
	}()

	return travisChan
}
