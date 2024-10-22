package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.tanvirRifat.io/internal/models"
	"snippetbox.tanvirRifat.io/internal/validator"
)

// for snippetcreate struct

type SnippetCreate struct{
    Title string `form:"title"`
    Content string `form:"content"`
    Expires int    `form:"expires"`
    // FieldErrors map[string]string

    // ignore this when parsing using the go-playground/form decoder package
    validator.Validator `form:"-"`
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
    

    // using the go-playground/form decode package

    var form SnippetCreate

    err:= app.decodePostForm(r,&form)

    if err != nil {
        app.ClientError(w, http.StatusBadRequest)
        return
    }




 form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
    form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
    form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
    form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

    



        if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
        return
    }


    id,err:=app.snippets.Insert(form.Title,form.Content,form.Expires)


    if err!=nil{
        app.ServerError(w,r,err)
        return
    }


    // creating a flash/toash message and pass it to the view handler
    // because after successfully created  we redirect to the snippetview handler
    app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

    http.Redirect(w,r,fmt.Sprintf("/snippet/view/%d",id),http.StatusSeeOther)




}


func (app *App) userSignup(w http.ResponseWriter,r *http.Request){

    fmt.Fprintln(w,"display the user signup form...")

}

func (app *App) userSignupPost(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w,"Displaying post signup...")
}

func (app *App) userLogin(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w,"Displaying user login...")
}

func (app *App) userLoginPost(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w,"Displaying user login post ...")
}

func (app *App) userLogoutPost(w http.ResponseWriter, r *http.Request){
    fmt.Fprintln(w,"Displaying the user logout post...")
}