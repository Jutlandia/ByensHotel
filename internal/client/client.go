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

// Create user:
// tip: https://github.com/go-ldap/ldap/issues/95
// https://github.com/go-ldap/ldap/blob/f61ea45f3b2b60217901f6110d5fef3193812062/v3/add.go#L29

func SetUp(sessionKey string, env string, c config.LDAP) {
	cfg = c
	store = newStore(sessionKey, env)
}

func Authenticate(w http.ResponseWriter, r *http.Request, username string, password string) error {
	l, err := connect()
	if err != nil {
		return err
	}
	defer l.Close()
	userDN, err := getUserDN(l, username, password)
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

func CreateUser(username string, email string, password string) error {
	l, err := connect()
	if err != nil {
		return err
	}
	defer l.Close()
	err = alreadyExist(l, username, email)
	if err != nil {
		return err
	}
	return createUser(l, username, email, password)
}

func getUserDN(l *ldap.Conn, username string, password string) (string, error) {
	err := l.Bind(cfg.BindUsername, cfg.BindPassword)
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
	return userdn, nil
}

func connect() (*ldap.Conn, error) {
	l, err := ldap.DialURL(fmt.Sprintf("ldap://%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return nil, err
	}
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}
	return l, nil
}

func alreadyExist(l *ldap.Conn, username string, email string) error {
	err := l.Bind(cfg.BindUsername, cfg.BindPassword)
	if err != nil {
		return err
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
		return err
	}
	if len(res.Entries) > 0 {
		return fmt.Errorf("a user with that username already exist")
	}
	sr = ldap.NewSearchRequest(
		cfg.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=*)(mail=%s))", ldap.EscapeFilter(email)),
		[]string{"dn"},
		nil,
	)
	res, err = l.Search(sr)
	if err != nil {
		return err
	}
	if len(res.Entries) > 0 {
		return fmt.Errorf("a user with that email already exist")
	}
	return nil
}

func createUser(l *ldap.Conn, username string, email string, password string) error {
	err := l.Bind(cfg.BindUsername, cfg.BindPassword)
	if err != nil {
		return err
	}
	req := ldap.AddRequest{
		DN: fmt.Sprintf("cn=%s,%s", username, cfg.BaseDN),
		Attributes: []ldap.Attribute{
			{
				Type: "objectClass",
				Vals: []string{"inetOrgPerson"},
			},
			{
				Type: "cn",
				Vals: []string{username},
			},
			{
				Type: "sn",
				Vals: []string{username},
			},
			{
				Type: "uid",
				Vals: []string{username},
			},
			{
				Type: "ou",
				Vals: []string{"people"},
			},
			{
				Type: "mail",
				Vals: []string{email},
			},
			{
				Type: "userPassword",
				Vals: []string{password},
			},
		},
	}
	return l.Add(&req)
}

func newStore(sessionKey string, env string) *sessions.CookieStore {
	s := sessions.NewCookieStore([]byte(sessionKey))
	s.Options.HttpOnly = true
	s.Options.Secure = env == "production"
	return s
}
