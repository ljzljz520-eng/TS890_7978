package storage

import (
	"fmt"
	"go.etcd.io/bbolt"
	"petcolor/internal/domain"
)

type UnitOfWork struct{ tx *bbolt.Tx }

func (s *Store) Update(fn func(*UnitOfWork) error) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return fn(&UnitOfWork{tx: tx}) })
}

func (u *UnitOfWork) PutClip(clip domain.ClipAsset) error {
	if err := clip.Validate(); err != nil {
		return err
	}
	data, err := encode(clip)
	if err != nil {
		return err
	}
	return u.tx.Bucket(clipsBucket).Put([]byte(clip.ID), data)
}
func (u *UnitOfWork) PutGrade(grade domain.GradeSession) error {
	if err := grade.Validate(); err != nil {
		return err
	}
	data, err := encode(grade)
	if err != nil {
		return err
	}
	return u.tx.Bucket(gradesBucket).Put([]byte(grade.ClipID), data)
}
func (u *UnitOfWork) PutPreview(frame domain.PreviewFrame) error {
	data, err := encode(frame)
	if err != nil {
		return err
	}
	return u.tx.Bucket(previewsBucket).Put(previewKey(frame), data)
}
func (u *UnitOfWork) PutEvent(event domain.AuditEvent) error {
	data, err := encode(event)
	if err != nil {
		return err
	}
	key := []byte(fmt.Sprintf("%020d/%s", event.OccurredAt.UnixNano(), event.ID))
	return u.tx.Bucket(eventsBucket).Put(key, data)
}
func (u *UnitOfWork) Clip(id string) (domain.ClipAsset, error) {
	var out domain.ClipAsset
	err := decode(u.tx.Bucket(clipsBucket).Get([]byte(id)), &out)
	return out, err
}
func (u *UnitOfWork) Grade(id string) (domain.GradeSession, error) {
	var out domain.GradeSession
	err := decode(u.tx.Bucket(gradesBucket).Get([]byte(id)), &out)
	return out, err
}
func (u *UnitOfWork) RemoveClip(id string) error { return u.tx.Bucket(clipsBucket).Delete([]byte(id)) }
func (u *UnitOfWork) RemoveGrade(id string) error {
	return u.tx.Bucket(gradesBucket).Delete([]byte(id))
}
