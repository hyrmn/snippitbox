package forms

import (
	"strings"
	"time"
	"unicode/utf8"
)

type NewSnippet struct {
	Title    string
	Content  string
	Expires  int
	Failures map[string]string
}

func (f *NewSnippet) GetExpirationTime() time.Time {
	return time.Now().Add(time.Duration(f.Expires) * time.Second)
}

func (f *NewSnippet) Valid() bool {
	f.Failures = make(map[string]string)
	// Check that the Title field is not blank and is not more than 100 characters
	// long. If it fails either of those checks, add a message to the f.Failures
	// map using the field name as the key.
	if strings.TrimSpace(f.Title) == "" {
		f.Failures["Title"] = "Title is required"
	} else if utf8.RuneCountInString(f.Title) > 100 {
		f.Failures["Title"] = "Title cannot be longer than 100 characters"
	}
	// Validate the Content and Expires fields aren't blank in a similar way.
	if strings.TrimSpace(f.Content) == "" {
		f.Failures["Content"] = "Content is required"
	}
	// Check that the Expires field isn't blank and is one of a fixed list. Using
	// a lookup on a map keyed with the permitted options and values of true is a
	// neat trick which saves you looping over the permitted values.
	permitted := map[int]bool{3600: true, 86400: true, 31536000: true}
	if f.Expires == 0 {
		f.Failures["Expires"] = "Expiry time is required"
	} else if !permitted[f.Expires] {
		f.Failures["Expires"] = "Expiry time must be 3600, 86400 or 31536000 seconds"
	}
	// If there are no failure messages, return true.
	return len(f.Failures) == 0
}
