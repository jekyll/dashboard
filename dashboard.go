package dashboard

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

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

func Listen() {
	mux := goji.NewMux()
	mux.HandleFuncC(pat.Get("/:name.json"), show)
	mux.HandleFunc(pat.Get("/"), index)

	bind := ":"
	if port := os.Getenv("PORT"); port != "" {
		bind += port
	} else {
		bind += "8000"
	}
	log.Fatal(http.ListenAndServe(bind, mux))
}
