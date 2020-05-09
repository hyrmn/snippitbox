package models

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

type Store struct {
	DB *bolt.DB
}

func (s *Store) update(bucket string, key []byte, dataStruct interface{}) error {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		// Create the bucket
		bucket, e := tx.CreateBucketIfNotExists([]byte(bucket))
		if e != nil {
			return e
		}

		// Encode the record
		encoded, e := json.Marshal(dataStruct)
		if e != nil {
			return e
		}

		// Store the record
		e = bucket.Put(key, encoded)
		if e != nil {
			return e
		}
		return nil
	})
	return err
}

func (s *Store) get(bucket string, key []byte, dataStruct interface{}) error {
	err := s.DB.View(func(tx *bolt.Tx) error {
		// Get the bucket
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return bolt.ErrBucketNotFound
		}

		// Retrieve the record
		encoded := bucket.Get(key)
		if len(encoded) < 1 {
			return bolt.ErrInvalid
		}

		// Decode the record
		e := json.Unmarshal(encoded, &dataStruct)
		if e != nil {
			return e
		}

		return nil
	})

	return err
}
