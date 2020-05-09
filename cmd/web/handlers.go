package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func Home(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			app.NotFound(w) // Use the app.NotFound() helper.
			return
		}
		app.RenderHTML(w, "home.page.html") // Use the app.RenderHTML() helper.
	})
}

func ShowSnippet(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedID := r.URL.Query().Get("id")
		id, err := strconv.Atoi(requestedID)
		if err != nil || id < 1 {
			app.NotFound(w)
			return
		}

		fmt.Fprintf(w, "Display a specific snippet (ID %d)...", id)
	})
}

func NewSnippet(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Display the new snippet form..."))
	})
}
