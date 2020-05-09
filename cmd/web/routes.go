package main

import (
	"net/http"
)

func (app *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", Home(app))
	mux.Handle("/snippet", ShowSnippet(app))
	mux.Handle("/snippet/new", NewSnippet(app))

	fileServer := http.FileServer(http.Dir(app.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
