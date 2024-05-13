package web

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/toastsandwich/letsgo-api/internal/models"
)

func (a *Application) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.NotFound(w)
		return
	}
	a.InfoLog.Println("Home")
	snippets, err := a.SnippetModel.Latest()
	if err != nil {
		a.ServerError(w, err)
		return
	}

	data := a.newTemplateData(r)
	data.Snippets = snippets

	a.Render(w, http.StatusOK, "home.tmpl", data)
}

func (a *Application) SnippetView(w http.ResponseWriter, r *http.Request) {
	idRaw := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil || id < 1 {
		a.NotFound(w)
		return
	}
	snippet, err := a.SnippetModel.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			a.NotFound(w)
		} else {
			fmt.Println("getting id error")
			a.ServerError(w, err)
		}
		return
	}

	data := a.newTemplateData(r)
	data.Snippet = snippet

	a.Render(w, http.StatusOK, "view.tmpl", data)
}

func (a *Application) SnippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allowed", http.MethodPost)
		a.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := a.SnippetModel.Insert(title, content, expires)
	if err != nil {
		a.ServerError(w, err)
		return
	}

	a.InfoLog.Println("snippet created")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
