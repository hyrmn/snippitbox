package main

import (
	"encoding/binary"
	"flag"
	"log"
	"net/http"

	"github.com/hyrmn/snippetbox/pkg/models"
	bolt "go.etcd.io/bbolt"
)

func main() {
	port := flag.String("port", "4000", "HTTP port to listen on")
	htmlDir := flag.String("html-dir", "./ui/html", "Path to HTML templates")
	staticDir := flag.String("static-dir", "./ui/static", "Path to static assets")
	dataDir := flag.String("data-dir", ".\\", "Path to BoltDB")
	flag.Parse()

	db, err := bolt.Open(*dataDir+"my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &App{
		Store:     &models.Store{DB: db},
		Port:      *port,
		HTMLDir:   *htmlDir,
		StaticDir: *staticDir,
	}

	log.Printf("Starting server on %s", app.Port)

	err = http.ListenAndServe(":"+app.Port, app.Routes())
	log.Fatal(err)
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
