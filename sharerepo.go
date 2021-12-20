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

    staticBox := rice.MustFindBox("static")
    staticFileServer := http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox()))
    http.Handle("/static/", staticFileServer)

    http.HandleFunc("/favicon.ico", handleNoContent)

    http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/picker":
        renderTemplate("picker.html", RepoPageData{}, w, r)
    case "/v2":
    case "/v2/":
        renderTemplate("v2.html", RepoPageData{}, w, r)
    default:
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

        renderTemplate("repo.html", data, w, r)
    }
}

func renderTemplate(tmplFile string, data RepoPageData, w http.ResponseWriter, r *http.Request) {
    tmplBox, err := rice.FindBox("templates")
    if err != nil {
	log.Fatal(err)
    }

    tmplString, err := tmplBox.String(tmplFile)
    if err != nil {
            log.Fatal(err)
    }

    tmpl, err := template.New(tmplFile).Parse(tmplString)
    if err != nil {
            log.Fatal(err)
    }

    tmpl.Execute(w, data)
}

func handleNoContent(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNoContent)
}
