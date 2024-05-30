package web

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/toastsandwich/letsgo-api/internal/models"
)

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	UserModel      *models.UserModel
	SnippetModel   *models.SnippetModel
	TemplateCache  map[string]*template.Template
	SessionManager *scs.SessionManager
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

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &Application{
		ErrorLog:       errorLog,
		InfoLog:        infoLog,
		UserModel:      &models.UserModel{DB: db},
		SnippetModel:   &models.SnippetModel{DB: db},
		TemplateCache:  templateCache,
		SessionManager: sessionManager,
	}

	addr := *host + ":" + *port
	infoLog.Println("starting server at", addr)

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
	}

	srv := &http.Server{
		Addr:         addr,
		ErrorLog:     errorLog,
		Handler:      app.Routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err.Error())
}
