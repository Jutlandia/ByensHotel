package client

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/Jutlandia/ByensHotel/internal/config"
	"github.com/go-ldap/ldap"
	"github.com/gorilla/sessions"
)

var (
	cfg   config.LDAP
	store sessions.Store
)

func SetUp(sessionKey string, env string, c config.LDAP) {
	cfg = c
	store = newStore(sessionKey, env)
}

func Authenticate(w http.ResponseWriter, r *http.Request, username string, password string) error {
	userDN, err := getUserDN(username, password)
	if err != nil {
		return err
	}
	session, _ := store.Get(r, "user")
	session.Values["dn"] = userDN
	return session.Save(r, w)
}

func IsAuthenticated(r *http.Request) bool {
	session, _ := store.Get(r, "user")
	if _, found := session.Values["dn"]; !found {
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

func getUserDN(username string, password string) (string, error) {
	l, err := ldap.DialURL(fmt.Sprintf("ldap://%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return "", err
	}
	defer l.Close()
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return "", err
	}
	err = l.Bind(cfg.BindUsername, cfg.BindPassword)
	if err != nil {
		return "", err
	}
	sr := ldap.NewSearchRequest(
		cfg.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=*)(uid=%s))", ldap.EscapeFilter(username)),
		[]string{"dn"},
		nil,
	)
	res, err := l.Search(sr)
	if err != nil {
		return "", err
	}
	if len(res.Entries) != 1 {
		return "", fmt.Errorf("invalid username or password")
	}
	userdn := res.Entries[0].DN
	err = l.Bind(userdn, password)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}
	err = l.Bind(cfg.BindUsername, cfg.BindPassword)
	if err != nil {
		return "", err
	}
	return userdn, nil
}

func newStore(sessionKey string, env string) *sessions.CookieStore {
	s := sessions.NewCookieStore([]byte(sessionKey))
	s.Options.HttpOnly = true
	s.Options.Secure = env == "production"
	return s
}
