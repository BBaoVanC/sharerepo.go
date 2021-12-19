package main

import (
    "log"
    "text/template"
    "net/http"
    "github.com/GeertJohan/go.rice"
)

type RepoPageData struct {
    URL string
}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/picker", handlePicker)

    staticBox := rice.MustFindBox("static")
    staticFileServer := http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox()))
    http.Handle("/static/", staticFileServer)

    http.HandleFunc("/favicon.ico", handleNoContent)

    http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
    tmplBox, err := rice.FindBox("templates")
    if err != nil {
        log.Fatal(err)
    }

    tmplString, err := tmplBox.String("repo.html")
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Processing a request")
    repos, ok := r.URL.Query()["repo"]
    if !ok || len(repos[0]) < 1 {
        log.Println("Missing repo parameter, redirecting to picker")
        http.Redirect(w, r, "/picker", 303)
        return
    }

    repo := repos[0]
    data := RepoPageData {
        URL: repo,
    }
    log.Println("Repo is", repo)
    tmpl, err := template.New("repo").Parse(tmplString)
    if err != nil {
        log.Fatal(err)
    }
    tmpl.Execute(w, data)
}

func handlePicker(w http.ResponseWriter, r *http.Request) {
    tmplBox, err := rice.FindBox("templates")
    if err != nil {
	log.Fatal(err)
    }

    tmplString, err := tmplBox.String("picker.html")
    if err != nil {
            log.Fatal(err)
    }

    tmpl, err := template.New("picker").Parse(tmplString)
    if err != nil {
            log.Fatal(err)
    }

    tmpl.Execute(w, nil)
}

func handleNoContent(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNoContent)
}
