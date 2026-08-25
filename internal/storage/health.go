package storage

import (
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
)

type Health struct {
	Path     string         `json:"path"`
	Writable bool           `json:"writable"`
	Buckets  map[string]int `json:"buckets"`
	Error    string         `json:"error,omitempty"`
}

func (s *Store) Health() Health {
	h := Health{Path: s.path, Buckets: map[string]int{}}
	err := s.db.View(func(tx *bbolt.Tx) error {
		for _, name := range []string{"clips", "grades", "previews", "events", "metadata"} {
			bucket := tx.Bucket([]byte(name))
			if bucket == nil {
				return fmt.Errorf("bucket %s missing", name)
			}
			h.Buckets[name] = bucket.Stats().KeyN
		}
		return nil
	})
	if err != nil {
		h.Error = err.Error()
		return h
	}
	err = s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(metadataBucket)
		if bucket == nil {
			return errors.New("metadata bucket missing")
		}
		return bucket.Put([]byte("health"), []byte("ok"))
	})
	if err != nil {
		h.Error = err.Error()
		return h
	}
	h.Writable = true
	return h
}

func (s *Store) Compact(destination string) error {
	if destination == "" {
		return errors.New("destination is required")
	}
	return s.db.View(func(tx *bbolt.Tx) error { return tx.CopyFile(destination, 0600) })
}
