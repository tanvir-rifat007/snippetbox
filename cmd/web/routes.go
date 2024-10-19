package main

import "net/http"


func (app *App) routes() *http.ServeMux{
	// static files:

		mux := http.NewServeMux()

		fileServer:=http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

		mux.Handle("/static/", http.StripPrefix("/static", fileServer))

		mux.HandleFunc("GET /{$}", app.home)
		mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
		mux.HandleFunc("GET /snippet/create", app.snippetCreate)
		mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

		return mux
}