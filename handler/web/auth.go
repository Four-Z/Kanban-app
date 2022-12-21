package web

import (
	"a21hc3NpZ25tZW50/client"
	"embed"
	"fmt"
	"log"
	"net/http"
	"path"
	"text/template"
)

type AuthWeb interface {
	Login(w http.ResponseWriter, r *http.Request)
	LoginProcess(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	RegisterProcess(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type authWeb struct {
	userClient client.UserClient
	embed      embed.FS
}

func NewAuthWeb(userClient client.UserClient, embed embed.FS) *authWeb {
	return &authWeb{userClient, embed}
}

func (a *authWeb) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	var filepath = path.Join("views", "auth", "login.html")
	var header = path.Join("views", "general", "header.html")

	var tmpl = template.Must(template.ParseFS(a.embed, filepath, header))

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *authWeb) LoginProcess(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	userId, status, err := a.userClient.Login(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status == 200 {
		http.SetCookie(w, &http.Cookie{
			Name:   "user_id",
			Value:  fmt.Sprintf("%d", userId),
			Path:   "/",
			MaxAge: 31536000,
			Domain: "",
		})

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func (a *authWeb) Register(w http.ResponseWriter, r *http.Request) {

	id := fmt.Sprintf("%s", r.Context().Value("id"))
	log.Print(id)

	if id == "" {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	var filepath = path.Join("views", "auth", "register.html")
	var header = path.Join("views", "general", "header.html")

	var tmpl = template.Must(template.ParseFS(a.embed, filepath, header))

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (a *authWeb) RegisterProcess(w http.ResponseWriter, r *http.Request) {
	fullname := r.FormValue("fullname")
	email := r.FormValue("email")
	password := r.FormValue("password")

	userId, status, err := a.userClient.Register(fullname, email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status == 200 {
		http.SetCookie(w, &http.Cookie{
			Name:   "user_id",
			Value:  fmt.Sprintf("%d", userId),
			Path:   "/",
			MaxAge: 31536000,
			Domain: "",
		})

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}
}

func (a *authWeb) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
		Domain: "",
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
