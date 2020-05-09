package models

import (
	bolt "go.etcd.io/bbolt"
)

type Store struct {
	DB *bolt.DB
}
