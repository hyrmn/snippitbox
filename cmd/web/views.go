package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func (app *App) RenderHTML(w http.ResponseWriter, page string) {
	files := []string{
		filepath.Join(app.HTMLDir, "base.html"),
		filepath.Join(app.HTMLDir, page),
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.ServerError(w, err) // Use the new app.ServerError() helper.
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.ServerError(w, err) // Use the new app.ServerError() helper.
	}
}
