package handler_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Jutlandia/ByensHotel/internal/client"
	"github.com/Jutlandia/ByensHotel/internal/handler"
	"github.com/Jutlandia/ByensHotel/internal/storage"
	"github.com/Jutlandia/ByensHotel/internal/tmpl"
)

var (
	env             = "test"
	sessionKey      = "my-32-byte-session-key"
	templatesLoaded = false
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

func loadTemplates(wd string) error {
	err := os.Chdir(filepath.Join(wd, "..", ".."))
	if err != nil {
		return err
	}
	tmpl.Load([]string{
		"index.html",
		"404.html",
		"auth/login.html",
		"auth/register.html",
	})
	templatesLoaded = true
	err = os.Chdir(wd)
	if err != nil {
		return err
	}
	return nil
}

func setUp() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if !templatesLoaded {
		err = loadTemplates(pwd)
		if err != nil {
			log.Fatal(err)
		}
	}
	ts := testStorage{}
	client.SetUp(ts, sessionKey, env)
}

func createResponse(method string, path string, body io.Reader, h http.HandlerFunc) *http.Response {
	setUp()
	req := httptest.NewRequest(method, path, body)
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	w := httptest.NewRecorder()
	h(w, req)
	return w.Result()
}

func TestLoginGet(t *testing.T) {
	resp := createResponse(http.MethodGet, "/login", nil, handler.Login)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code: %d\nGot status code: %d\n",
			http.StatusOK, resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	expectedContents := []string{
		"<form method=\"POST\" class=\"px-3\">",
		"<input type=\"text\" id=\"username\" name=\"username\" class=\"input\" value=\"\" required/>",
		"<input type=\"password\" id=\"password\" name=\"password\" class=\"input\" value=\"\" required/>",
	}
	for _, content := range expectedContents {
		if !bytes.Contains(body, []byte(content)) {
			t.Errorf("Expected: %s\n", content)
		}
	}
}

func TestLoginRedirectIfSuccess(t *testing.T) {
	formData := strings.NewReader("username=fry&password=fry")
	resp := createResponse(http.MethodPost, "/login", formData, handler.Login)
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Expected status code: %d\nGot: %d\n",
			http.StatusSeeOther, resp.StatusCode)
	}
	if resp.Header.Get("Location") != "/" {
		t.Errorf("Expected Location: /\nGot: %s\n",
			resp.Header.Get("Location"))
	}
}

func TestLoginErrorMsg(t *testing.T) {
	formData := strings.NewReader("username=&password=")
	resp := createResponse(http.MethodPost, "/login", formData, handler.Login)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code: %d\nGot: %d\n",
			http.StatusOK, resp.StatusCode)
	}
	expectedContents := []string{
		"<span class=\"help is-danger\">Please enter your username</span>",
		"<span class=\"help is-danger\">Please enter your password</span>",
	}
	body, _ := io.ReadAll(resp.Body)
	for _, content := range expectedContents {
		if !bytes.Contains(body, []byte(content)) {
			t.Errorf("Expected: %s\n", content)
		}
	}
}

func TestRegisterGet(t *testing.T) {
	resp := createResponse(http.MethodPost, "/register", nil, handler.Register)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code: %d\nGot status code: %d\n",
			http.StatusOK, resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	expectedContents := []string{
		"<form method=\"POST\" class=\"px-3\">",
		"<input type=\"text\" id=\"username\" name=\"username\" class=\"input\" value=\"\" required/>",
		"<input type=\"email\" id=\"email\" name=\"email\" class=\"input\" value=\"\" required/>",
		"<input type=\"password\" id=\"password\" name=\"password\" class=\"input\" value=\"\" required/>",
		"<input type=\"password\" id=\"confirmPassword\" name=\"confirmPassword\" class=\"input\" value=\"\" required/>",
	}
	for _, content := range expectedContents {
		if !bytes.Contains(body, []byte(content)) {
			t.Errorf("Expected: %s\n", content)
		}
	}
}

func TestRegisterRedirectIfSuccess(t *testing.T) {
	formData := "username=alice&email=test@mail.com&password=123123&confirmPassword=123123"
	resp := createResponse(http.MethodPost, "/register", strings.NewReader(formData), handler.Register)
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Expected status code: %d\nGot: %d\n",
			http.StatusSeeOther, resp.StatusCode)
	}
	if resp.Header.Get("Location") != "/login" {
		t.Errorf("Expected Location: /login\nGot: %s\n", resp.Header.Get("Location"))
	}
}

func TestRegisterErrorMsg(t *testing.T) {
	formData := strings.NewReader("username=&email=&password=&confirmPassword=")
	resp := createResponse(http.MethodPost, "/register", formData, handler.Register)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code: %d\nGot: %d\n",
			http.StatusOK, resp.StatusCode)
	}
	expectedContents := []string{
		"<span class=\"help is-danger\">Please enter a username</span>",
		"<span class=\"help is-danger\">Please enter an email</span>",
		"<span class=\"help is-danger\">Please enter a password</span>",
		"<span class=\"help is-danger\">Please confirm your password</span>",
	}
	body, _ := io.ReadAll(resp.Body)
	for _, content := range expectedContents {
		if !bytes.Contains(body, []byte(content)) {
			t.Errorf("Expected: %s\n", content)
		}
	}
}

func TestNotFound(t *testing.T) {
	res := createResponse(http.MethodGet, "/not-found", nil, handler.NotFound)
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code: %d\nGot: %d\n",
			http.StatusNotFound, res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	if !bytes.Contains(body, []byte("404 - Not Found")) {
		t.Error("Expected 404 page to contain \"404 - Not Found\"")
	}
}
