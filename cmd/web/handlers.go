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


type UserSignupForm struct{
    Name string `form:"name"`
    Email string `form:"email"`
    Password string `form:"password"`
    validator.Validator `form:"-"`
}

type UserLoginForm struct{
    Email string `form:"email"`
    Password string `form:"password"`
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

// get request to see the userSingup form
func (app *App) userSignup(w http.ResponseWriter,r *http.Request){

    data:= app.newTemplateData(r)
    data.Form = UserSignupForm{}

    app.render(w,r,http.StatusOK,"signup.tmpl.html",data)



}

// post req for the usersignup form

func (app *App) userSignupPost(w http.ResponseWriter, r *http.Request){
    var form UserSignupForm

    err:= app.decodePostForm(r,&form)
  // form decode korte somossa hole to client error ei hobe
    if err!=nil{
        app.ClientError(w,http.StatusBadRequest)
        return

    }

    form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
    form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
    form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
    form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
    form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
        return
    }

    // Try to create a new user record in the database. If the email already
    // exists then add an error message to the form and re-display it.
    err = app.users.Insert(form.Name, form.Email, form.Password)
    if err != nil {
        if errors.Is(err, models.ErrDuplicateEmail) {
            form.AddFieldError("email", "Email address is already in use")

            data := app.newTemplateData(r)
            data.Form = form
            app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
        } else {
            app.ServerError(w, r, err)
        }

        return
    }

    // Otherwise add a confirmation flash message to the session confirming that
    // their signup worked.
    app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

    // And redirect the user to the login page.
    http.Redirect(w, r, "/user/login", http.StatusSeeOther)



}

func (app *App) userLogin(w http.ResponseWriter, r *http.Request){
    data:= app.newTemplateData(r)
    data.Form = UserLoginForm{}

    app.render(w,r,http.StatusOK,"login.tmpl.html",data)
}

func (app *App) userLoginPost(w http.ResponseWriter, r *http.Request) {
    // Decode the form data into the userLoginForm struct.
    var form UserLoginForm

    err := app.decodePostForm(r, &form)
    if err != nil {
        app.ClientError(w, http.StatusBadRequest)
        return
    }

    // Do some validation checks on the form. We check that both email and
    // password are provided, and also check the format of the email address as
    // a UX-nicety (in case the user makes a typo).
    form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
    form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
    form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
        return
    }

    // Check whether the credentials are valid. If they're not, add a generic
    // non-field error message and re-display the login page.
    id, err := app.users.Authenticate(form.Email, form.Password)
    if err != nil {
        if errors.Is(err, models.ErrInvalidCredentials) {
            form.AddNonFieldError("Email or password is incorrect")

            data := app.newTemplateData(r)
            data.Form = form
            app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
        } else {
            app.ServerError(w, r, err)
        }
        return
    }

   






    // Use the RenewToken() method on the current session to change the session
    // ID. It's good practice to generate a new session ID when the 
    // authentication state or privilege levels changes for the user (e.g. login
    // and logout operations).
    err = app.sessionManager.RenewToken(r.Context())
    if err != nil {
        app.ServerError(w, r, err)
        return
    }




    // Add the ID of the current user to the session, so that they are now
    // 'logged in'.
    app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

    // Redirect the user to the create snippet page.
    http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}


func (app *App) userLogoutPost(w http.ResponseWriter, r *http.Request){
    // create a new Authentication session:

    app.sessionManager.RenewToken(r.Context())

    // remove the authetication id:
    app.sessionManager.Remove(r.Context(),"authenticatedUserID")

    app.sessionManager.Put(r.Context(),"flash","You've been logged out successfully!")


    http.Redirect(w,r,"/",http.StatusSeeOther)

    


}