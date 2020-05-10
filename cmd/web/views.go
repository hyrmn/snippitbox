package main

import (
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/hyrmn/snippetbox/pkg/models"
)

type HTMLData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func friendlyDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func (app *App) RenderHTML(w http.ResponseWriter, page string, data *HTMLData) {
	files := []string{
		filepath.Join(app.HTMLDir, "base.html"),
		filepath.Join(app.HTMLDir, page),
	}

	fm := template.FuncMap{
		"friendlyDate": friendlyDate,
	}

	ts, err := template.New("").Funcs(fm).ParseFiles(files...)
	if err != nil {
		app.ServerError(w, err) // Use the new app.ServerError() helper.
		return
	}

	buf := new(bytes.Buffer)

	err = ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServerError(w, err) // Use the new app.ServerError() helper.
		return
	}

	buf.WriteTo(w)
}
