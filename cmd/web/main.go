package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.tanvirRifat.io/internal/models"
)


type App struct{
    logger *slog.Logger
    snippets *models.SnippetModel
    templateCache map[string]*template.Template
    formDecoder *form.Decoder
    sessionManager *scs.SessionManager
}

func main() {


   logger:= slog.New(slog.NewTextHandler(os.Stdout,nil))
   dsn:= flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

   db,err:= openDB(*dsn)

   if err != nil {
         logger.Error(err.Error())
         os.Exit(1)
   }

   defer db.Close()

   templateCache,err:=newTemplateCache()

   if err!=nil{
    logger.Error(err.Error())
    os.Exit(1)
   }

   formDecoder:= form.NewDecoder()


       sessionManager := scs.New()
    sessionManager.Store = mysqlstore.New(db)
    sessionManager.Lifetime = 12 * time.Hour

    // session is used only by https:
    sessionManager.Cookie.Secure = true


    tlsConfig := &tls.Config{
        CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
    }



    // dependency injection:
    app:= &App{
        logger: logger,
        snippets: &models.SnippetModel{DB:db},
        templateCache: templateCache,
        formDecoder: formDecoder,
        sessionManager: sessionManager,

    }

    addr:=flag.String("addr", ":4000", "HTTP network address")

    flag.Parse()

    srv:= &http.Server{
        Addr: *addr,
        Handler: app.routes(),
        // handle any error together with our slog.

        ErrorLog: slog.NewLogLogger(logger.Handler(),slog.LevelError),

        TLSConfig: tlsConfig,
        IdleTimeout:  time.Minute,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,

    }
 
    logger.Info("starting server","addr",srv.Addr)
    
    err = srv.ListenAndServeTLS("./tls/cert.pem","./tls/key.pem")

    logger.Error(err.Error())

    os.Exit(1)
}



func openDB(dsn string) (*sql.DB,error){
    db,err:=sql.Open("mysql", dsn)

    if err != nil {
        return nil,err
    }

    if err = db.Ping(); err != nil {

         return nil,err    
    }

    return db,nil
}



// hide the static path

type neuteredFileSystem struct {
    fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
    f, err := nfs.fs.Open(path)
    if err != nil {
        return nil, err
    }

    s, err := f.Stat()
    if err != nil {
        return nil, err
    }
    
    if s.IsDir() {
        index := filepath.Join(path, "index.html")
        if _, err := nfs.fs.Open(index); err != nil {
            closeErr := f.Close()
            if closeErr != nil {
                return nil, closeErr
            }

            return nil, err
        }
    }

    return f, nil
}    

