package client

import (
	"fmt"
	"net/http"

	"github.com/Jutlandia/ByensHotel/internal/storage"
	"github.com/gorilla/sessions"
)

var (
	users storage.UserStorage
	store sessions.Store
)

func SetUp(s storage.UserStorage, sessionKey string, env string) {
	users = s
	store = newStore(sessionKey, env)
}

func Authenticate(w http.ResponseWriter, r *http.Request, username string, password string) error {
	user, err := users.GetWithCredentials(username, password)
	if err != nil {
		return err
	}
	session, _ := store.Get(r, "user")
	session.Values["username"] = user.Username()
	return session.Save(r, w)
}

func IsAuthenticated(r *http.Request) bool {
	session, _ := store.Get(r, "user")
	if _, found := session.Values["username"]; !found {
		return false
	}
	return session.Options.MaxAge > 0
}

func ClearSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := store.Get(r, "user")
	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = 0
	return session.Save(r, w)
}

func CreateUser(username string, email string, password string) error {
	err := alreadyExist(username, email)
	if err != nil {
		return err
	}
	return users.Create(username, email, password)
}

func alreadyExist(username string, email string) error {
	user, err := users.GetByUsername(username)
	if err != nil {
		return err
	}
	if user != nil {
		return fmt.Errorf("a user with that username already exist")
	}
	user, err = users.GetByEmail(email)
	if err != nil {
		return err
	}
	if user != nil {
		return fmt.Errorf("a user with that email already exist")
	}
	return nil
}

func newStore(sessionKey string, env string) *sessions.CookieStore {
	s := sessions.NewCookieStore([]byte(sessionKey))
	s.Options.HttpOnly = true
	s.Options.Secure = env == "production"
	return s
}
