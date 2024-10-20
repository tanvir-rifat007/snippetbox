package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.tanvirRifat.io/internal/models"
)

func (app *App) home(w http.ResponseWriter, r *http.Request) {
    
    snippets, err := app.snippets.Latest()
    if err != nil {
        app.ServerError(w, r, err)
        return
    }


    data := app.newTemplateData(r)
    data.Snippets = snippets




    app.render(w, r, http.StatusOK, "home.tmpl.html", data)

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

    data:= app.newTemplateData(r)
    data.Snippet = snippet

    app.render(w,r,http.StatusOK,"view.tmpl.html",data)



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
