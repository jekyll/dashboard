package dashboard

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	gh "github.com/google/go-github/v50/github"
)

var allowRequestEvery = time.Second / 30
var throttle <-chan time.Time = time.Tick(allowRequestEvery)

func logHTTP(method, url string, f func()) {
	log.Println("==> ", method, url)
	start := time.Now()
	f()
	log.Println("==> ", method, url, "finished in", time.Since(start))
}

func get(url string, data interface{}) error {
	<-throttle
	var resp *http.Response
	var err error
	logHTTP("GET", url, func() {
		resp, err = http.Get(url)
	})
	if err != nil {
		return err
	}
	return json.NewDecoder(resp.Body).Decode(data)
}

func doGraphql(client *gh.Client, query string, variables map[string]interface{}, output interface{}) error {
	req, err := githubClient.NewRequest(
		"POST",
		"/graphql",
		map[string]interface{}{"query": query, "variables": variables},
	)
	if err != nil {
		return err
	}

	logHTTP(req.Method, req.URL.String(), func() {
		_, err = githubClient.Do(context.Background(), req, output)
	})
	return err
}

func getRetry(retries int, url string, data interface{}) error {
	var err error
	tries := 0
	for tries <= retries {
		tries += 1
		err = get(url, data)
		if err == nil {
			break
		}
	}
	return err
}
