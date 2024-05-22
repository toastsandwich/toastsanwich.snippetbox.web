package web

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/toastsandwich/letsgo-api/internal/models"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	SnippetModel  *models.SnippetModel
	TemplateCache map[string]*template.Template
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/snippet/create", app.SnippetCreate)
	mux.HandleFunc("/snippet/view", app.SnippetView)
	return mux
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func Start() {
	// creating a logger for info | error
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//define host for the app
	host := flag.String("host", "localhost", "HTTP Host Address")
	//define port for the app
	port := flag.String("port", "4000", "HTTP Port Number")
	// define dsn for MYSQL
	dsn := flag.String("dsn", "gomon:smpmsmim@/snippetbox?parseTime=true", "MYSQL data source name")

	flag.Parse()

	db, err := OpenDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTempateCache()
	if err != nil {
		errorLog.Fatal(err.Error())
	}

	app := &Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		SnippetModel:  &models.SnippetModel{DB: db},
		TemplateCache: templateCache,
	}

	addr := *host + ":" + *port
	infoLog.Println("starting server at ", addr)
	// errorLog.Fatal(http.ListenAndServe(addr, mux))

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}

	for key, val := range templateCache {
		fmt.Println(key, ":", *val)
	}

	err = srv.ListenAndServe()
	errorLog.Fatal(err.Error())
}
