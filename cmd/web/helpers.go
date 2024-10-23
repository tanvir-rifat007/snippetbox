package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
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
         // for floashig or toast message
		 Flash: app.sessionManager.PopString(r.Context(), "flash"),
         IsAuthenticated: app.isAuthenticated(r),

         CSRFToken: nosurf.Token(r),
         

	}

}


func (app *App) decodePostForm(r *http.Request, dst any) error {
    // Call ParseForm() on the request, in the same way that we did in our
    // snippetCreatePost handler.
    err := r.ParseForm()
    if err != nil {
        return err
    }


    err = app.formDecoder.Decode(dst, r.PostForm)
    if err != nil {

			// if we pass the wrong dst then there is the invalidDecoderError
        
        var invalidDecoderError *form.InvalidDecoderError
        
        if errors.As(err, &invalidDecoderError) {
            panic(err)
        }

        return err
    }

    return nil
}

// for authorization
// mane e holo login obosthay asi kina

func (app *App) isAuthenticated(r *http.Request) bool {
    return app.sessionManager.Exists(r.Context(), "authenticatedUserID") 
}

