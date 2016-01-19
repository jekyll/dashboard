package main

import (
	"encoding/json"
	"log"
	"net/http"

	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
)

func show(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	name := pat.Param(ctx, "name")
	json.NewEncoder(w).Encode(getProject(name))
}

func index(w http.ResponseWriter, r *http.Request) {
	indexTmpl.Execute(w, templateInfo{Projects: getAllProjects()})
}

func main() {
	mux := goji.NewMux()
	mux.HandleFuncC(pat.Get("/:name.json"), show)
	mux.HandleFunc(pat.Get("/"), index)
	log.Fatal(http.ListenAndServe(":8000", mux))
}
