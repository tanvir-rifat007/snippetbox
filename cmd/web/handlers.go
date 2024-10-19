package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"snippetbox.tanvirRifat.io/internal/models"
)

func (app *App) home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")
    
    snippets, err := app.snippets.Latest()
    if err != nil {
        app.ServerError(w, r, err)
        return
    }

    for _, snippet := range snippets {
        fmt.Fprintf(w, "%+v\n", snippet)
    }

    // files := []string{
    //     "./ui/html/base.tmpl.html",
    //     "./ui/html/partials/nav.tmpl.html",
    //     "./ui/html/pages/view.tmpl.html",
    // }

    // ts, err := template.ParseFiles(files...)
    // if err != nil {
    //     app.ServerError(w, r, err)
    //     return
    // }

    // err = ts.ExecuteTemplate(w, "base", nil)
    // if err != nil {
    //     app.ServerError(w, r, err)
    // }
}



func (app *App)snippetView(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }


    snippet, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            http.NotFound(w, r)
        } else {
            app.ServerError(w, r, err)
        }
        return
    }

    // Write the snippet data as a plain-text HTTP response body.
    // fmt.Fprintf(w, "%+v", snippet)


    files:= []string{
        "./ui/html/base.tmpl.html",
        "./ui/html/partials/nav.tmpl.html",
        "./ui/html/pages/view.tmpl.html",
    }

    ts, err := template.ParseFiles(files...)
    
    if err != nil {
        app.ServerError(w, r, err)
        return
    }

    data:= templateData{
        Snippet: snippet,
    }

    err = ts.ExecuteTemplate(w, "base", data)

    if err != nil {
        app.ServerError(w, r, err)
    }

}

func (app *App)snippetCreate(w http.ResponseWriter, r *http.Request) {


    title := "O snail"
    content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
    expires := 7

   
    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.ServerError(w, r, err)
        return
    }

    // Redirect the user to the relevant page for the snippet.
    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}

func (app *App)snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Save a new snippet..."))
}
