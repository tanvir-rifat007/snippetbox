package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"snippetbox.tanvirRifat.io/internal/models"
)

// for snippetcreate struct

type SnippetCreate struct{
    Title string
    Content string
    Expires int
    FieldErrors map[string]string
}

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


   data:= app.newTemplateData(r)

  data.Form= SnippetCreate{
    Expires: 365,
  }

   app.render(w,r,http.StatusCreated,"create.tmpl.html",data)


}

func (app *App)snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    err:=r.ParseForm()

    if err!=nil{
        app.ClientError(w,http.StatusBadRequest)
        return
    }

    title:=r.PostForm.Get("title")
    content:=r.PostForm.Get("content")
    expires:=r.PostForm.Get("expires")



    intExpires,err:= strconv.Atoi(expires)

    form:= SnippetCreate{
        Title: title,
        Content: content,
        Expires: intExpires,
        FieldErrors: map[string]string{},
    }



    if strings.TrimSpace(title)==""{
        form.FieldErrors["title"]="Title is required"

    }else if utf8.RuneCountInString(title) > 100{
        form.FieldErrors["title"]="Title is too long"
    }

    if strings.TrimSpace(content) == "" {
        form.FieldErrors["content"] = "Content is required"
    }

    if strings.TrimSpace(expires) == "" {
        form.FieldErrors["expires"] = "Expiry time is required"
    }else if intExpires<1 || intExpires>365{
        form.FieldErrors["expires"]="Expiry time must be between 1 and 365 days"
    }

    if err!=nil{
        app.ClientError(w,http.StatusBadRequest)
        return
    }

    if len(form.FieldErrors)>0{
        data:= app.newTemplateData(r)
        data.Form = form

        

        app.render(w,r,http.StatusUnprocessableEntity,"create.tmpl.html",data)
        return
        

    }

    id,err:=app.snippets.Insert(form.Title,form.Content,form.Expires)


    if err!=nil{
        app.ServerError(w,r,err)
        return
    }

    http.Redirect(w,r,fmt.Sprintf("/snippet/view/%d",id),http.StatusSeeOther)




}
