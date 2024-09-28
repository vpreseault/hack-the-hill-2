package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func renderTemplate(w http.ResponseWriter, fileName string) error {
	tmpl, err := template.ParseFiles(filepath.Join("..", "frontend", fileName + ".html"))
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}
