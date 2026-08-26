package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.etcd.io/bbolt"
	"petcolor/internal/domain"
)

var clipsBucket = []byte("clips")
var gradesBucket = []byte("grades")
var previewsBucket = []byte("previews")
var eventsBucket = []byte("events")
var metadataBucket = []byte("metadata")

type Store struct {
	db   *bbolt.DB
	path string
	mu   sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("storage path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}
	db, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: time.Second, NoGrowSync: false})
	if err != nil {
		return nil, fmt.Errorf("open storage: %w", err)
	}
	s := &Store{db: db, path: path}
	if err = s.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range [][]byte{clipsBucket, gradesBucket, previewsBucket, eventsBucket, metadataBucket} {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
func (s *Store) Path() string { return s.path }

func encode(value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode entity: %w", err)
	}
	return data, nil
}
func decode(data []byte, value any) error {
	if len(data) == 0 {
		return domain.ErrNotFound
	}
	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("decode entity: %w", err)
	}
	return nil
}

func (s *Store) SaveClip(clip domain.ClipAsset) error {
	if err := clip.Validate(); err != nil {
		return err
	}
	data, err := encode(clip)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(clipsBucket).Put([]byte(clip.ID), data) })
}

func (s *Store) Clip(id string) (domain.ClipAsset, error) {
	var out domain.ClipAsset
	err := s.db.View(func(tx *bbolt.Tx) error { return decode(tx.Bucket(clipsBucket).Get([]byte(id)), &out) })
	return out, err
}

func (s *Store) DeleteClip(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(clipsBucket).Delete([]byte(id)) })
}

func (s *Store) Clips() ([]domain.ClipAsset, error) {
	out := []domain.ClipAsset{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(clipsBucket).ForEach(func(_, v []byte) error {
			var item domain.ClipAsset
			if err := decode(v, &item); err != nil {
				return err
			}
			out = append(out, item)
			return nil
		})
	})
	return out, err
}

func (s *Store) SaveGrade(grade domain.GradeSession) error {
	if err := grade.Validate(); err != nil {
		return err
	}
	data, err := encode(grade)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(gradesBucket).Put([]byte(grade.ClipID), data) })
}

func (s *Store) Grade(clipID string) (domain.GradeSession, error) {
	var out domain.GradeSession
	err := s.db.View(func(tx *bbolt.Tx) error { return decode(tx.Bucket(gradesBucket).Get([]byte(clipID)), &out) })
	return out, err
}

func (s *Store) SaveGradeIfRevision(grade domain.GradeSession, expected uint64) error {
	if err := grade.Validate(); err != nil {
		return err
	}
	data, err := encode(grade)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(gradesBucket)
		var current domain.GradeSession
		if err := decode(bucket.Get([]byte(grade.ClipID)), &current); err != nil {
			return err
		}
		if current.Revision != expected {
			return domain.ErrConflict
		}
		return bucket.Put([]byte(grade.ClipID), data)
	})
}

// UpdateGradeAtomic runs merge against the freshest grade for clipID inside a
// single read-write transaction. Because the read, the patch, and the write all
// share one transaction, a competing update cannot land between them: the merge
// always sees the latest committed revision and the write is conditional on
// that same revision. This is record isolation for grade updates — concurrent
// patches to distinct fields merge instead of one silently clobbering the other.
func (s *Store) UpdateGradeAtomic(clipID string, merge func(current domain.GradeSession) (domain.GradeSession, error)) (domain.GradeSession, error) {
	if merge == nil {
		return domain.GradeSession{}, errors.New("merge function is required")
	}
	var updated domain.GradeSession
	err := s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(gradesBucket)
		var current domain.GradeSession
		if err := decode(bucket.Get([]byte(clipID)), &current); err != nil {
			return err
		}
		merged, err := merge(current)
		if err != nil {
			return err
		}
		if err := merged.Validate(); err != nil {
			return err
		}
		data, err := encode(merged)
		if err != nil {
			return err
		}
		if err := bucket.Put([]byte(clipID), data); err != nil {
			return err
		}
		updated = merged
		return nil
	})
	return updated, err
}

func previewKey(frame domain.PreviewFrame) []byte {
	return []byte(fmt.Sprintf("%s/%010d/%s", frame.ClipID, frame.Sequence, frame.ID))
}

func (s *Store) SavePreview(frame domain.PreviewFrame) error {
	data, err := encode(frame)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(previewsBucket).Put(previewKey(frame), data) })
}

func (s *Store) Preview(id string) (domain.PreviewFrame, error) {
	var out domain.PreviewFrame
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(previewsBucket).ForEach(func(_, v []byte) error {
			var item domain.PreviewFrame
			if err := decode(v, &item); err != nil {
				return err
			}
			if item.ID == id {
				out = item
				return errors.New("found")
			}
			return nil
		})
	})
	if err != nil && err.Error() == "found" {
		return out, nil
	}
	if err != nil {
		return out, err
	}
	return out, domain.ErrNotFound
}

func (s *Store) Previews(clipID string) ([]domain.PreviewFrame, error) {
	out := []domain.PreviewFrame{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(previewsBucket).ForEach(func(_, v []byte) error {
			var item domain.PreviewFrame
			if err := decode(v, &item); err != nil {
				return err
			}
			if clipID == "" || item.ClipID == clipID {
				out = append(out, item)
			}
			return nil
		})
	})
	return out, err
}

func (s *Store) DeletePreviews(clipID string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(previewsBucket)
		keys := [][]byte{}
		if err := bucket.ForEach(func(k, v []byte) error {
			var item domain.PreviewFrame
			if err := decode(v, &item); err != nil {
				return err
			}
			if item.ClipID == clipID {
				keys = append(keys, append([]byte(nil), k...))
			}
			return nil
		}); err != nil {
			return err
		}
		for _, key := range keys {
			if err := bucket.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) SaveEvent(event domain.AuditEvent) error {
	data, err := encode(event)
	if err != nil {
		return err
	}
	key := []byte(fmt.Sprintf("%020d/%s", event.OccurredAt.UnixNano(), event.ID))
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(eventsBucket).Put(key, data) })
}

func (s *Store) Events() ([]domain.AuditEvent, error) {
	out := []domain.AuditEvent{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(eventsBucket).ForEach(func(_, v []byte) error {
			var item domain.AuditEvent
			if err := decode(v, &item); err != nil {
				return err
			}
			out = append(out, item)
			return nil
		})
	})
	return out, err
}

func (s *Store) DeleteEventsFor(target string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(eventsBucket)
		keys := [][]byte{}
		if err := bucket.ForEach(func(k, v []byte) error {
			var item domain.AuditEvent
			if err := decode(v, &item); err != nil {
				return err
			}
			if item.TargetID == target {
				keys = append(keys, append([]byte(nil), k...))
			}
			return nil
		}); err != nil {
			return err
		}
		for _, key := range keys {
			if err := bucket.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) SetMetadata(key, value string) error {
	if key == "" {
		return errors.New("metadata key is required")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(metadataBucket).Put([]byte(key), []byte(value)) })
}
func (s *Store) Metadata(key string) (string, error) {
	var value string
	err := s.db.View(func(tx *bbolt.Tx) error {
		raw := tx.Bucket(metadataBucket).Get([]byte(key))
		if raw == nil {
			return domain.ErrNotFound
		}
		value = string(raw)
		return nil
	})
	return value, err
}

func (s *Store) Snapshot() (map[string]int, error) {
	result := map[string]int{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		for name, bucketName := range map[string][]byte{"clips": clipsBucket, "grades": gradesBucket, "previews": previewsBucket, "events": eventsBucket} {
			count := 0
			if err := tx.Bucket(bucketName).ForEach(func(_, v []byte) error {
				if v != nil {
					count++
				}
				return nil
			}); err != nil {
				return err
			}
			result[name] = count
		}
		return nil
	})
	return result, err
}
