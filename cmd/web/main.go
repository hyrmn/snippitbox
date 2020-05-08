package main

import (
	"flag"
	"log"
	"net/http"
)

type Config struct {
	Addr      string
	HTMLDir   string
	StaticDir string
}

func main() {
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", "4000", "HTTP network address")
	flag.StringVar(&cfg.HTMLDir, "html-dir", "./ui/html", "Path to HTML templates")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	app := &App{
		HTMLDir: cfg.HTMLDir,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/snippet", app.ShowSnippet)
	mux.HandleFunc("/snippet/new", app.NewSnippet)

	fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Starting server on %s", cfg.Addr)

	err := http.ListenAndServe(":"+cfg.Addr, mux)
	log.Fatal(err)
}
