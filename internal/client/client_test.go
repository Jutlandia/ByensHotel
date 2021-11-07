package client_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jutlandia/ByensHotel/internal/client"
	"github.com/Jutlandia/ByensHotel/internal/config"
)

var (
	env        = "test"
	sessionKey = "my-32-byte-session-key"
	testLDAP   = config.LDAP{
		Host:         "localhost",
		Port:         10389,
		BindUsername: "cn=Hubert J. Farnsworth,ou=people,dc=planetexpress,dc=com",
		BindPassword: "professor",
		BaseDN:       "dc=planetexpress,dc=com",
	}
)

func TestAuthenticate(t *testing.T) {
	client.SetUp(sessionKey, env, testLDAP)
	r := httptest.NewRequest(http.MethodPost, "/login", nil)
	w := httptest.NewRecorder()
	err := client.Authenticate(w, r, "fry", "fry")
	if err != nil {
		log.Fatal(err)
	}
	if client.IsAuthenticated(r) == false {
		t.Error("Expected user to be authenticated")
	}
}

func TestClearSession(t *testing.T) {
	client.SetUp(sessionKey, env, testLDAP)
	r := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()
	err := client.Authenticate(w, r, "fry", "fry")
	if err != nil {
		log.Fatal(err)
	}
	err = client.ClearSession(w, r)
	if err != nil {
		log.Fatal(err)
	}
	if client.IsAuthenticated(r) == true {
		t.Error("Expected session to be cleared")
	}
}
