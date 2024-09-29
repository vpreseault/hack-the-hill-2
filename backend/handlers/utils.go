package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func renderTemplate(w http.ResponseWriter, fileName string) error {
	tmpl, err := template.ParseFiles(filepath.Join(".", "static", fileName + ".html"))
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func saveSessionIDInCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:  "sessionID",
		Value: sessionID,
	}
	http.SetCookie(w, cookie)
}
