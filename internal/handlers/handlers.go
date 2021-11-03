package handlers

import (
	"html/template"
	"net/http"

	"github.com/Jutlandia/ByensHotel/internal/forms"
	"github.com/Jutlandia/ByensHotel/internal/templates"
	"github.com/gorilla/csrf"
)

type authPageData struct {
	Form forms.Form
	CSRF template.HTML
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	templates.Render(w, "index.html", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templates.Render(w, "login.html", authPageData{
			Form: &forms.LoginForm{},
			CSRF: csrf.TemplateField(r),
		})
	case http.MethodPost:
		lf := &forms.LoginForm{
			Username: r.PostFormValue("username"),
			Password: r.PostFormValue("password"),
		}
		if !lf.IsValid() {
			templates.Render(w, "login.html", authPageData{
				Form: lf,
				CSRF: csrf.TemplateField(r),
			})
			return
		}
		// TODO: verify that the credentials are correct
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templates.Render(w, "register.html", authPageData{
			Form: &forms.RegisterForm{},
			CSRF: csrf.TemplateField(r),
		})
	case http.MethodPost:
		rf := &forms.RegisterForm{
			Username:        r.PostFormValue("username"),
			Email:           r.PostFormValue("email"),
			Password:        r.PostFormValue("password"),
			ConfirmPassword: r.PostFormValue("confirmPassword"),
		}
		if !rf.IsValid() {
			templates.Render(w, "register.html", authPageData{
				Form: rf,
				CSRF: csrf.TemplateField(r),
			})
			return
		}
		// TODO: check if user already exists
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
