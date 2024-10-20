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

		mux.HandleFunc("GET /{$}", app.home)
		mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
		mux.HandleFunc("GET /snippet/create", app.snippetCreate)
		mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

		// using the default middleware chaining:


		// return app.recoverPanic(app.logRequest(commonHeader(mux)))


		// using the third party alice package:

		standard:= alice.New(app.recoverPanic,app.logRequest,commonHeader)

		return standard.Then(mux)
}