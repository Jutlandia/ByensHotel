package handler

import (
	"html/template"
	"net/http"

	"github.com/Jutlandia/ByensHotel/internal/client"
	"github.com/Jutlandia/ByensHotel/internal/form"
	"github.com/Jutlandia/ByensHotel/internal/tmpl"
	"github.com/gorilla/csrf"
)

type authPageData struct {
	Form form.Form
	CSRF template.HTML
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "index.html", nil)
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	if client.IsAuthenticated(r) {
		err := client.ClearSession(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	redirectIfAuthenticated(w, r)
	switch r.Method {
	case http.MethodGet:
		tmpl.Render(w, r, "login.html", authPageData{
			Form: &form.Login{},
			CSRF: csrf.TemplateField(r),
		})
	case http.MethodPost:
		lf := &form.Login{
			Username: r.PostFormValue("username"),
			Password: r.PostFormValue("password"),
		}
		if !lf.IsValid() {
			tmpl.Render(w, r, "login.html", authPageData{
				Form: lf,
				CSRF: csrf.TemplateField(r),
			})
			return
		}
		err := client.Authenticate(w, r, lf.Username, lf.Password)
		if err != nil {
			// TODO: make this better with custom error names, e.g. ErrBadCredentials.
			if err.Error() == "invalid username or password" {
				lf.AddError("Overall", "Invalid username or password")
			} else {
				lf.AddError("Overall", "Something went wrong")
			}
			tmpl.Render(w, r, "login.html", authPageData{
				Form: lf,
				CSRF: csrf.TemplateField(r),
			})
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	redirectIfAuthenticated(w, r)
	switch r.Method {
	case http.MethodGet:
		tmpl.Render(w, r, "register.html", authPageData{
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
			tmpl.Render(w, r, "register.html", authPageData{
				Form: rf,
				CSRF: csrf.TemplateField(r),
			})
			return
		}
		// TODO: check if user already exists
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// redirectIfAuthenticated redirects users to "/" if they try to access an auth route
// and already already authenticated.
// TODO: come up with a better way of doing this.
func redirectIfAuthenticated(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet || r.Method == http.MethodPost {
		if client.IsAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}
