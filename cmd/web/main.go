package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

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

	var snippet models.Snippet

	dberr := db.Update(func(tx *bolt.Tx) error {
		// Create the bucket
		bucket, e := tx.CreateBucketIfNotExists([]byte("Snippets"))
		if e != nil {
			return e
		}
		id, _ := bucket.NextSequence()

		snippet = models.Snippet{
			ID:      id,
			Title:   "Example title",
			Content: "Example content",
			Created: time.Now(),
			Expires: time.Now(),
		}
		// Encode the record
		encodedRecord, e := json.Marshal(snippet)
		if e != nil {
			return e
		}

		// Store the record
		if e = bucket.Put(itob(snippet.ID), encodedRecord); e != nil {
			return e
		}
		fmt.Printf("Saved %d\r\n", snippet.ID)
		return nil
	})

	if dberr != nil {
		log.Fatal(dberr)
	}
	read := &models.Snippet{}

	dberr = db.View(func(tx *bolt.Tx) error {
		// Get the bucket
		bucket := tx.Bucket([]byte("Snippets"))
		if bucket == nil {
			return bolt.ErrBucketNotFound
		}

		// Retrieve the record
		encoded := bucket.Get(itob(snippet.ID))
		if len(encoded) < 1 {
			return bolt.ErrInvalid
		}

		// Decode the record
		e := json.Unmarshal(encoded, &read)
		if e != nil {
			return e
		}

		return nil
	})

	if dberr != nil {
		log.Fatal(dberr)
	}

	log.Println(read)

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
