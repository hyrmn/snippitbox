package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *App) Routes() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", app.Home)
	router.GET("/snippet/", app.NewSnippet)
	router.GET("/snippet/:id", app.ShowSnippet)
	//fileServer := http.FileServer(http.Dir(app.StaticDir))

	router.ServeFiles("/static/*filepath", http.Dir(app.StaticDir))

	return router
}
