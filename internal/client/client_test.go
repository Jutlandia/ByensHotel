package client_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jutlandia/ByensHotel/internal/client"
	"github.com/Jutlandia/ByensHotel/internal/storage"
)

var (
	env        = "test"
	sessionKey = "my-32-byte-session-key"
)

type testUser struct {
}

func (tu testUser) Username() string { return "fry" }

func (tu testUser) Email() string { return "fry@planetexpress.com" }

type testStorage struct {
}

func (ts testStorage) GetWithCredentials(username string, password string) (storage.User, error) {
	if username == "fry" && password == "fry" {
		return testUser{}, nil
	}
	return nil, fmt.Errorf("invalid username or password")
}

func (ts testStorage) GetByUsername(username string) (storage.User, error) {
	return nil, nil
}

func (ts testStorage) GetByEmail(email string) (storage.User, error) {
	return nil, nil
}

func (ts testStorage) Create(username string, email string, password string) error {
	return nil
}

func TestAuthenticate(t *testing.T) {
	ts := testStorage{}
	client.SetUp(ts, sessionKey, env)
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
	ts := testStorage{}
	client.SetUp(ts, sessionKey, env)
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
