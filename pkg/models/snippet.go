package models

import (
	"encoding/binary"
	"encoding/json"
	"time"

	bolt "go.etcd.io/bbolt"
)

type Snippet struct {
	ID      uint64
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

const bucketName = "Snippets"

// multiple Snippet objects.
type Snippets []*Snippet

func GetSnippet(s *Store, id uint64) (*Snippet, error) {
	snippet := Snippet{}

	err := loadByID(s, id, &snippet)

	return &snippet, err
}

func GetLatestSnippets(s *Store, numToFetch uint) (Snippets, error) {
	snippets := Snippets{}

	err := s.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucketName))

		c := b.Cursor()
		var fetched uint = 0

		for k, v := c.Last(); k != nil && fetched < numToFetch; k, v = c.Prev() {
			snippet := &Snippet{}
			e := json.Unmarshal(v, &snippet)
			if e != nil {
				return e
			}

			snippets = append(snippets, snippet)
			fetched++
		}

		return nil
	})

	return snippets, err
}

func SaveSnippet(s *Store, title string, content string, expires time.Time) (*Snippet, error) {
	snippet := &Snippet{
		Title:   title,
		Content: content,
		Expires: expires,
	}

	err := upsert(s, snippet)

	if err != nil {
		return nil, err
	}

	return snippet, nil
}

func loadByID(s *Store, id uint64, snippet *Snippet) error {
	err := s.DB.View(func(tx *bolt.Tx) error {
		// Get the bucket
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrNoResult
		}

		// Retrieve the record
		encoded := bucket.Get(keyFromID(id))
		if encoded == nil {
			return ErrNoResult
		}

		// Decode the record
		e := json.Unmarshal(encoded, &snippet)
		if e != nil {
			return e
		}

		return nil
	})

	return err
}

func upsert(s *Store, snippet *Snippet) error {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		// Create the bucket
		bucket, e := tx.CreateBucketIfNotExists([]byte(bucketName))
		if e != nil {
			return e
		}

		if snippet.ID == 0 {
			id, _ := bucket.NextSequence()
			snippet.ID = id
		}

		// Encode the record
		encoded, e := json.Marshal(snippet)
		if e != nil {
			return e
		}

		// Store the record
		e = bucket.Put(keyFromID(snippet.ID), encoded)
		if e != nil {
			return e
		}
		return nil
	})
	return err
}

func keyFromID(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
