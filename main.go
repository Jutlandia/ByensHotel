package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Jutlandia/ByensHotel/internal/filesystem"
	"github.com/Jutlandia/ByensHotel/internal/handlers"
	"github.com/Jutlandia/ByensHotel/internal/templates"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	env     string
	port    string
	csrfKey string
)

func init() {
	env = os.Getenv("HOTEL_ENV")
	if env == "" {
		env = "development"
	}
	if env != "production" {
		godotenv.Load()
	}
	csrfKey = os.Getenv("CSRF_KEY")
	if csrfKey == "" {
		log.Fatal("missing csrf key")
	}
	port = os.Getenv("PORT")
	if port == "" {
		port = "1337"
	}
	templates.Load([]string{
		"index.html",
		"auth/login.html",
		"auth/register.html",
		"layouts/base.html",
	})
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	assets := filepath.Join(pwd, "web", "assets")
	fileserver := http.FileServer(filesystem.FileSystem{http.Dir(assets)})

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HomeHandler).Methods(http.MethodGet)
	r.HandleFunc("/login", handlers.LoginHandler).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/register", handlers.RegisterHandler).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/favicon.ico", favIconHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileserver))

	srv := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%s", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: csrf.Protect(
			[]byte(csrfKey),
			csrf.Secure(env == "production"),
		)(r),
	}

	log.Printf("Hosting environment: %s\n", env)
	log.Printf("Now listening on: %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

// TODO: Delete this handler when favicon.ico has been added to web/assets
func favIconHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
