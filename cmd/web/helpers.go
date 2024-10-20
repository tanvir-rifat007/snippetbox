package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
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


		// creating a buffer 
		buf:= new(bytes.Buffer)

		// writing the template to the buffer not the whole response like(w)

    err := ts.ExecuteTemplate(buf, "base", data)
    if err != nil {
        app.ServerError(w, r, err)
    }
		 w.WriteHeader(status)

		 buf.WriteTo(w)

}


func (app *App) newTemplateData(r *http.Request)templateData {
	return templateData{
		 CurrentYear: time.Now().Year(),
	}

}