package handler_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Jutlandia/ByensHotel/internal/handler"
	"github.com/Jutlandia/ByensHotel/internal/tmpl"
)

var templatesLoaded = false

func loadTemplates(wd string) error {
	err := os.Chdir(filepath.Join(wd, "..", ".."))
	if err != nil {
		return err
	}
	tmpl.Load([]string{
		"index.html",
		"auth/login.html",
		"auth/register.html",
		"layouts/base.html",
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
}

func getResult(method string, path string, body io.Reader, h http.HandlerFunc) *http.Response {
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
	resp := getResult(http.MethodGet, "/login", nil, handler.Login)
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
	formData := strings.NewReader("username=alice&password=123123")
	resp := getResult(http.MethodPost, "/login", formData, handler.Login)
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Expected status code: %d\nGot: %d\n",
			http.StatusSeeOther, resp.StatusCode)
	}
	if resp.Header.Get("Location") != "/" {
		t.Errorf("Expected Location: /\nGot: %s\n", resp.Header.Get("Location"))
	}
}

func TestLoginErrorMsg(t *testing.T) {
	formData := strings.NewReader("username=&password=")
	resp := getResult(http.MethodPost, "/login", formData, handler.Login)
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
	resp := getResult(http.MethodPost, "/register", nil, handler.Register)
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
	resp := getResult(http.MethodPost, "/register", strings.NewReader(formData), handler.Register)
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
	resp := getResult(http.MethodPost, "/register", formData, handler.Register)
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