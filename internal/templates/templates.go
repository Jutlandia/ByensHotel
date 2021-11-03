package templates

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type templateData struct {
	LayoutData  interface{}
	ContentData interface{}
}

var templates map[string]*template.Template

func Load(tmpl []string) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	layoutPath := filepath.Join(pwd, "web", "templates", "layouts", "base.html")
	layoutData, err := os.ReadFile(layoutPath)
	if err != nil {
		log.Fatal(err)
	}
	layout := template.New("base")
	layout, err = layout.Parse(string(layoutData))
	if err != nil {
		log.Fatal(err)
	}
	layoutFiles := []string{layoutPath}
	for _, file := range tmpl {
		fileName := filepath.Base(file)
		files := append(layoutFiles, filepath.Join(pwd, "web", "templates", file))
		templates[fileName], err = layout.Clone()
		if err != nil {
			log.Fatal(err)
		}
		templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
	}
	log.Println("templates loaded successfully")
}

func Render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, found := templates[name]
	if !found {
		http.Error(w,
			fmt.Sprintf("template %s does not exist", name),
			http.StatusInternalServerError)
		return
	}
	tmplData := templateData{
		LayoutData: struct {
			Year int
		}{Year: time.Now().Year()},
		ContentData: data,
	}
	err := tmpl.Execute(w, tmplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
