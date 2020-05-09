package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hyrmn/snippetbox/pkg/models"
)

//Home handles all requests for the home page
func Home(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			app.NotFound(w) // Use the app.NotFound() helper.
			return
		}
		app.RenderHTML(w, "home.page.html") // Use the app.RenderHTML() helper.
	})
}

//ShowSnippet returns the snippet with the requested ID (using the url param "id")
func ShowSnippet(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedID := r.URL.Query().Get("id")
		id, err := strconv.ParseUint(requestedID, 10, 64)
		if err != nil || id < 1 {
			app.NotFound(w)
			return
		}

		snippet, err := models.GetSnippet(app.Store, id)
		if err == models.ErrNoResult {
			app.NotFound(w)
			return
		}
		if err != nil {
			app.ServerError(w, err)
		}

		js, err := json.Marshal(snippet)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	})
}

func NewSnippet(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Display the new snippet form..."))
	})
}
