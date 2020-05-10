package main

import (
	"net/http"
	"strconv"

	"github.com/hyrmn/snippetbox/pkg/models"
	"github.com/julienschmidt/httprouter"
)

//Home handles all requests for the home page
func (app *App) Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.URL.Path != "/" {
		app.NotFound(w) // Use the app.NotFound() helper.
		return
	}

	snippets, err := models.GetLatestSnippets(app.Store, 10)

	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.RenderHTML(w, r, "home.page.html", &HTMLData{
		Snippets: snippets,
	})
}

//ShowSnippet returns the snippet with the requested ID (using the url param "id")
func (app *App) ShowSnippet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestedID := ps.ByName("id")
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

	app.RenderHTML(w, r, "show.page.html", &HTMLData{
		Snippet: snippet,
	})
}

func (app *App) NewSnippet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Display the new snippet form..."))
}
