package main

import (
	"fmt"
	"net/http"
)


func (app *App) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	
	app.logger.Error(err.Error(),"method",r.Method,"path",r.URL.RequestURI())
	http.Error(w,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	
}

func (app *App) ClientError(w http.ResponseWriter, status int) {
    http.Error(w, http.StatusText(status), status)
}


func (app *App) render(w http.ResponseWriter,r *http.Request, status int, page string, data templateData){
	 ts,ok:=app.templateCache[page]

 if !ok {
        err := fmt.Errorf("the template %s does not exist", page)
        app.ServerError(w, r, err)
        return
    }


    w.WriteHeader(status)


    err := ts.ExecuteTemplate(w, "base", data)
    if err != nil {
        app.ServerError(w, r, err)
    }
}