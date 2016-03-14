package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func serveIndex(w http.ResponseWriter, _ *http.Request) error {
	return render(w, "index", nil)
}

func renderError(w http.ResponseWriter, owner, repo string) error {
	return render(w, "error", nil)
}

func renderSuccess(w http.ResponseWriter, owner, repo string) error {
	return render(w, "success", nil)
}

func render(w http.ResponseWriter, partial string, data interface{}) error {
	tmplName := func(name string) string { return fmt.Sprintf("templates/%s.html", name) }
	tmpl, err := template.ParseFiles(tmplName("layout"), tmplName(partial))
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, "layout", data)
}
