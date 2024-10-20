package main

import (
	"path/filepath"
	"text/template"
	"time"

	"snippetbox.tanvirRifat.io/internal/models"
)

// ei struct ti create kora hoise
// jate amra template e
// Snippet,Snippets,CurrentYear use korte pari
type templateData struct {

  CurrentYear int

	// view.tmpl.html e dekhanor jonne
	Snippet models.Snippet

	// home.tmpl.html e dekhanor jonne
	Snippets []models.Snippet
}


// custom template function:

func humanDate(t time.Time) string{
	    return t.Format("02 Jan 2006 at 15:04")

}

var functions = template.FuncMap{
	"humanDate":humanDate,
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

				// instead of this explicit pattern 
				// of the files 

				// I just dynamically add this:

        // files := []string{
        //     "./ui/html/base.tmpl.html",
        //     "./ui/html/partials/nav.tmpl.html",
        //     page,
        // }

				// ts,err:= template.ParseFiles("./ui/html/base.tmpl.html")

				// for the injecting humanDate functions to the template:
				ts,err:= template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")

				if err!=nil{
					return nil,err
				}

				// now add the dynamically any partials to the ts:

				ts,err= ts.ParseGlob("./ui/html/partials/nav.tmpl.html")

				if err!=nil{
					return nil,err
				}

				    

        // Parse the files into a template set.
        ts, err = ts.ParseFiles(page)
        if err != nil {
            return nil, err
        }


        cache[name] = ts
    }

    return cache, nil
}



