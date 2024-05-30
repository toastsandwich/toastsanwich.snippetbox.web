package web

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/toastsandwich/letsgo-api/internal/models"
	"github.com/toastsandwich/letsgo-api/internal/validator"
)

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

type UserCreateForm struct {
	Name, Email, Password string
	validator.Validator
}

type LoginForm struct {
	Email    string
	Password string
	validator.Validator
}

func (a *Application) Home(w http.ResponseWriter, r *http.Request) {
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

	params := httprouter.ParamsFromContext(r.Context())
	idRaw := params.ByName("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil && id < 1 {
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
	data := a.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	a.Render(w, http.StatusOK, "create.tmpl", data)
}

func (a *Application) SnippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		a.ClientError(w, http.StatusBadRequest)
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		a.ClientError(w, http.StatusBadRequest)
	}

	form := snippetCreateForm{
		Title:     r.PostForm.Get("title"),
		Content:   r.PostForm.Get("content"),
		Expires:   expires,
		Validator: validator.NewValidator(),
	}

	form.CheckField(validator.NotABlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.NotABlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.MaxCharLimit(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.PermittedInts(expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.Render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	a.SessionManager.Put(r.Context(), "flash", "Snippet Successfully Created")

	id, err := a.SnippetModel.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		a.ServerError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (a *Application) UserSignUp(w http.ResponseWriter, r *http.Request) {
	form := UserCreateForm{}
	data := a.newTemplateData(r)
	data.Form = form
	a.Render(w, http.StatusOK, "signup.tmpl", data)
}

func (a *Application) UserSignUpPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		a.ServerError(w, err)
		return
	}
	form := UserCreateForm{
		Name:      r.PostForm.Get("name"),
		Email:     r.PostForm.Get("email"),
		Password:  r.PostForm.Get("password"),
		Validator: validator.NewValidator(),
	}
	form.CheckField(validator.NotABlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotABlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotABlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChar(form.Password, 8), "password", "This field must be atleast 8 characters long")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.Render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err := a.UserModel.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "email already in use")
			data := a.newTemplateData(r)
			data.Form = form
			a.Render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			a.ServerError(w, err)
		}
		return
	}
	a.SessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *Application) UserLogin(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = LoginForm{}
	a.Render(w, http.StatusOK, "login.tmpl", data)
}

func (a *Application) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		a.ServerError(w, err)
	}

	form := LoginForm{
		Email:     r.PostForm.Get("email"),
		Password:  r.PostForm.Get("password"),
		Validator: validator.NewValidator(),
	}
	form.CheckField(validator.NotABlank(form.Email), "email", "This field cannot be blank.")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email.")
	form.CheckField(validator.NotABlank(form.Password), "password", "This field cannot be blank.")

	id, err := a.UserModel.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := a.newTemplateData(r)
			data.Form = form
			a.Render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			a.ServerError(w, err)
		}
		return
	}
	err = a.SessionManager.RenewToken(r.Context())
	if err != nil {
		a.ServerError(w, err)
		return
	}
	a.SessionManager.Put(r.Context(), "authenticatedUserID", id)
	a.SessionManager.Put(r.Context(), "flash", "Login Success")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *Application) UserLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := a.SessionManager.RenewToken(r.Context())
	if err != nil {
		a.ServerError(w, err)
	}
	a.SessionManager.Remove(r.Context(), "authencicatedUserID")
	a.SessionManager.Put(r.Context(), "flash", "Logged out success")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
