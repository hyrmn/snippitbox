package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/hyrmn/snippetbox/pkg/forms"
	"github.com/hyrmn/snippetbox/pkg/models"
	"github.com/julienschmidt/httprouter"
)

//Home handles all requests for the home page
func (app *App) Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	app.RenderHTML(w, r, "new.page.html", nil)
}

func (app *App) CreateSnippet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	expires, _ := strconv.Atoi(r.PostForm.Get("expires"))

	form := &forms.NewSnippet{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	if !form.Valid() {
		fmt.Fprint(w, form.Failures)
		return
	}

	newSnippet, err := models.SaveSnippet(app.Store, form.Title, form.Content, form.GetExpirationTime())
	if err != nil {
		app.ServerError(w, err)
		return
	}
	// If successful, send a 303 See Other response redirecting the user to the
	// page with their new snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", newSnippet.ID), http.StatusSeeOther)
}
