package handler

import (
	"go-19/model"
	"go-19/service"
	"html/template"
	"net/http"
	"strconv"
)

type AuthHandler struct {
	Tmpl    *template.Template
	Service service.UserService
}

func NewAuthHandler(tmpl *template.Template, service service.UserService) AuthHandler {
	return AuthHandler{Tmpl: tmpl, Service: service}
}

func (authHandler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	err := authHandler.Tmpl.ExecuteTemplate(w, "login", nil)
	if err != nil {
		http.Error(w, "Gagal menampilkan halaman login", http.StatusInternalServerError)
	}
}

func (authHandler *AuthHandler) DoLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Gagal membaca form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := authHandler.Service.FindByEmail(email)
	if err != nil || user == nil || user.Password != password {
		http.Error(w, "Email atau password salah", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    strconv.Itoa(user.ID),
		Path:     "/",
		HttpOnly: true,
	})

	if user.Role == "student" {
		http.Redirect(w, r, "/student/home", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/lecturer/home", http.StatusSeeOther)
	}
}

func (authHandler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	err := authHandler.Tmpl.ExecuteTemplate(w, "register", nil)
	if err != nil {
		http.Error(w, "Gagal menampilkan halaman register", http.StatusInternalServerError)
	}
}

func (authHandler *AuthHandler) DoRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Gagal membaca form", http.StatusBadRequest)
		return
	}

	user := &model.User{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Role:     r.FormValue("role"),
	}

	err := authHandler.Service.CreateUser(user)
	if err != nil {
		http.Error(w, "Gagal mendaftar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
