package filesystem_test

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Jutlandia/ByensHotel/internal/filesystem"
)

var temp string

var (
	folders = [3]string{"css", "js", "img"}
	files   = [3]string{"css/main.css", "js/main.js", "img/profile.png"}
)

func setUp() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	dir, err := os.MkdirTemp(pwd, "temp")
	if err != nil {
		return err
	}
	temp = dir
	for _, folderName := range folders {
		err = os.Mkdir(filepath.Join(temp, folderName), 0777)
		if err != nil {
			return err
		}
	}
	for _, file := range files {
		split := strings.Split(file, "/")
		_, err := os.Create(filepath.Join(temp, split[0], split[1]))
		if err != nil {
			return err
		}
	}
	return nil
}

func TestOpenReturnsFile(t *testing.T) {
	err := setUp()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(temp)

	filesys := filesystem.FileSystem{http.Dir(temp)}

	for _, filename := range files {
		_, err := filesys.Open(filename)
		if err != nil {
			t.Errorf("Expected file: %s\n", filename)
		}
	}
}

func TestOpenReturnsErrorIfFileDoesNotExist(t *testing.T) {
	err := setUp()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(temp)

	filesys := filesystem.FileSystem{http.Dir(temp)}

	_, err = filesys.Open("this-does-not-exist")
	if err == nil {
		t.Error("Expected opening a non-existing file to return an error")
	}
}

func TestOpenReturnsErrorIfFolderIsRequested(t *testing.T) {
	err := setUp()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(temp)

	filesys := filesystem.FileSystem{http.Dir(temp)}

	for _, folder := range folders {
		_, err := filesys.Open(folder)
		if err == nil {
			t.Errorf("Expected an error: %s\n", folder)
		}
	}
}
