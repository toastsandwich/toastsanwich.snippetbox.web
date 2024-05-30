package web

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *Application) Routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NotFound(w)
	})
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileserver))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileserver))
	// mux.HandleFunc("/", app.Home)
	// mux.HandleFunc("/snippet/create", app.SnippetCreate)
	// mux.HandleFunc("/snippet/view", app.SnippetView)

	dynamic := alice.New(app.SessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.Home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.SnippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.SnippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.SnippetCreatePost))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.UserSignUp))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.UserSignUpPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.UserLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.UserLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.UserLogoutPost))
	standard := alice.New(app.RecoverPanic, app.LogRequest, SecureHeaderMiddleWare)
	return standard.Then(router)
}
