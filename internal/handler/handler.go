package handler

import (
	"html/template"
	"net/http"

	"github.com/Jutlandia/ByensHotel/internal/form"
	"github.com/Jutlandia/ByensHotel/internal/tmpl"
	"github.com/gorilla/csrf"
)

type authPageData struct {
	Form form.Form
	CSRF template.HTML
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, "index.html", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl.Render(w, "login.html", authPageData{
			Form: &form.Login{},
			CSRF: csrf.TemplateField(r),
		})
	case http.MethodPost:
		lf := &form.Login{
			Username: r.PostFormValue("username"),
			Password: r.PostFormValue("password"),
		}
		if !lf.IsValid() {
			tmpl.Render(w, "login.html", authPageData{
				Form: lf,
				CSRF: csrf.TemplateField(r),
			})
			return
		}
		// TODO: verify that the credentials are correct
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl.Render(w, "register.html", authPageData{
			Form: &form.Register{},
			CSRF: csrf.TemplateField(r),
		})
	case http.MethodPost:
		rf := &form.Register{
			Username:        r.PostFormValue("username"),
			Email:           r.PostFormValue("email"),
			Password:        r.PostFormValue("password"),
			ConfirmPassword: r.PostFormValue("confirmPassword"),
		}
		if !rf.IsValid() {
			tmpl.Render(w, "register.html", authPageData{
				Form: rf,
				CSRF: csrf.TemplateField(r),
			})
			return
		}
		// TODO: check if user already exists
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
