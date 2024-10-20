package main

import (
	"path/filepath"
	"text/template"

	"snippetbox.tanvirRifat.io/internal/models"
)


type templateData struct {
	Snippet models.Snippet
	Snippets []models.Snippet
}


func newTemplateCache() (map[string]*template.Template, error) {
    cache := map[string]*template.Template{}


    pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
    if err != nil {
        return nil, err
    }

    for _, page := range pages {
        // get the only filename from the path(like:home.tmpl.html)
        name := filepath.Base(page)


        files := []string{
            "./ui/html/base.tmpl.html",
            "./ui/html/partials/nav.tmpl.html",
            page,
        }

        // Parse the files into a template set.
        ts, err := template.ParseFiles(files...)
        if err != nil {
            return nil, err
        }


        cache[name] = ts
    }

    return cache, nil
}
