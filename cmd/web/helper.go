package web

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (a *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.ErrorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a *Application) NotFound(w http.ResponseWriter) {
	a.ClientError(w, http.StatusNotFound)
}

func (a *Application) Render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := a.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		a.ServerError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		a.ServerError(w, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (a *Application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		Flash:       a.SessionManager.PopString(r.Context(), "flash"),
	}
}
