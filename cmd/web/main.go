package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "4000", "HTTP port to listen on")
	htmlDir := flag.String("html-dir", "./ui/html", "Path to HTML templates")
	staticDir := flag.String("static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	app := &App{
		Port:      *port,
		HTMLDir:   *htmlDir,
		StaticDir: *staticDir,
	}

	log.Printf("Starting server on %s", app.Port)

	err := http.ListenAndServe(":"+app.Port, app.Routes())
	log.Fatal(err)
}
