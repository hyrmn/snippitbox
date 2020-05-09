package models

import (
	"encoding/binary"
	"time"
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

func GetSnippet(s *Store, id uint64) (Snippet, error) {
	snippet := Snippet{}

	err := s.get(bucketName, keyFromID(id), &snippet)

	if err != nil {
		err = ErrNoResult
	}

	return snippet, err
}

func keyFromID(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
