package main

import (
	"net/http"

	"github.com/justinas/alice"
)


func (app *App) routes() http.Handler{
	// static files:

		mux := http.NewServeMux()

		fileServer:=http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

		mux.Handle("/static/", http.StripPrefix("/static", fileServer))

				dynamic:= alice.New(app.sessionManager.LoadAndSave)


		mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
		mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
		

		  mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
    mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
    mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
    mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))


		protected:= dynamic.Append(app.requireAuthentication)

		mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
		mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
    mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))



		// using the default middleware chaining:


		// return app.recoverPanic(app.logRequest(commonHeader(mux)))


		// using the third party alice package:

		standard:= alice.New(app.recoverPanic,app.logRequest,commonHeader)

		return standard.Then(mux)
}