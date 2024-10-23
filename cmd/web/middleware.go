// we want this middleware like below:
// middleware ----> servemux(our http.Handler)-----> Application handler

package main

import (
	"fmt"
	"net/http"
)


func commonHeader(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
        w.Header().Set("Content-Security-Policy",
            "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

        w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "deny")
        w.Header().Set("X-XSS-Protection", "0")

        w.Header().Set("Server", "Go")

				next.ServeHTTP(w,r)
	})
}

func (app *App) logRequest(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		        var (
            ip     = r.RemoteAddr
            proto  = r.Proto
            method = r.Method
            uri    = r.URL.RequestURI()
        )

        app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

        next.ServeHTTP(w, r)

	})
  
}

func (app *App) recoverPanic(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		defer func(){
			if err:= recover(); err!=nil{
				w.Header().Set("Connection", "close")
				app.ServerError(w,r,fmt.Errorf("%s",err))
			}
		}()
		next.ServeHTTP(w,r)

	})
}


// middleware for protected routes:

func (app *App) requireAuthentication(next http.Handler) http.Handler{
	
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		if !app.isAuthenticated(r){
			http.Redirect(w,r,"/user/login",http.StatusSeeOther)
			return
		}

		// which pages requires authentication are not stored any browser cache:
    w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w,r)
	})

}




