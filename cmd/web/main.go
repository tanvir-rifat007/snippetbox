package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.tanvirRifat.io/internal/models"
)


type App struct{
    logger *slog.Logger
    snippets *models.SnippetModel
    templateCache map[string]*template.Template
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

    // dependency injection:
    app:= &App{
        logger: logger,
        snippets: &models.SnippetModel{DB:db},
        templateCache: templateCache,

    }

    addr:=flag.String("addr", ":4000", "HTTP network address")

    flag.Parse()
 
    logger.Info("starting server","addr",*addr)
    
    err = http.ListenAndServe(*addr, app.routes())

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

