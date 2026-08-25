package service

import (
	"errors"
	"fmt"
	"path/filepath"
	"petcolor/internal/domain"
	"petcolor/internal/grading"
	"petcolor/internal/media"
	"petcolor/internal/preview"
	"petcolor/internal/storage"
	"sync"
	"time"
)

type Clock interface{ Now() time.Time }
type FixedClock struct{ Value time.Time }

func (c FixedClock) Now() time.Time { return c.Value }

type Service struct {
	store         *storage.Store
	uploads       media.UploadManager
	planner       preview.Planner
	clock         Clock
	hook          func(string)
	eventMu       sync.Mutex
	eventSequence uint64
}

func New(store *storage.Store, uploads media.UploadManager, planner preview.Planner, clock Clock) (*Service, error) {
	if store == nil {
		return nil, errors.New("store is required")
	}
	if clock == nil {
		return nil, errors.New("clock is required")
	}
	return &Service{store: store, uploads: uploads, planner: planner, clock: clock}, nil
}

func (s *Service) SetConcurrencyHook(hook func(string)) { s.hook = hook }

func (s *Service) nextEventID() string {
	s.eventMu.Lock()
	defer s.eventMu.Unlock()
	s.eventSequence++
	return fmt.Sprintf("event-%06d", s.eventSequence)
}

func (s *Service) event(kind, target, actor, summary string, revision uint64) domain.AuditEvent {
	event, _ := domain.NewAuditEvent(s.nextEventID(), kind, target, actor, summary, revision, s.clock.Now())
	return event
}

func (s *Service) RegisterUpload(request domain.UploadRequest, actor string) (domain.ClipSummary, error) {
	clip, err := domain.NewClipAsset(request)
	if err != nil {
		return domain.ClipSummary{}, err
	}
	grade := domain.DefaultGrade(clip.ID, s.clock.Now())
	event := s.event("clip.uploaded", clip.ID, actor, "短片上传并创建默认调色会话", grade.Revision)
	err = s.store.Update(func(unit *storage.UnitOfWork) error {
		if err := unit.PutClip(clip); err != nil {
			return err
		}
		if err := unit.PutGrade(grade); err != nil {
			return err
		}
		return unit.PutEvent(event)
	})
	if err != nil {
		return domain.ClipSummary{}, err
	}
	return domain.SummarizeClip(clip, grade, nil), nil
}

func (s *Service) ApplyPreset(clipID string, name domain.PresetName, actor string) (domain.GradeSession, error) {
	clip, err := s.store.Clip(clipID)
	if err != nil {
		return domain.GradeSession{}, err
	}
	if !clip.AvailableAt(s.clock.Now()) {
		return domain.GradeSession{}, domain.ErrExpired
	}
	current, err := s.store.Grade(clipID)
	if err != nil {
		return domain.GradeSession{}, err
	}
	patch, err := grading.PatchForPreset(name)
	if err != nil {
		return domain.GradeSession{}, err
	}
	updated := current.Apply(patch, s.clock.Now())
	if err := s.store.SaveGrade(updated); err != nil {
		return domain.GradeSession{}, err
	}
	if err := s.store.SaveEvent(s.event("grade.preset", clipID, actor, "应用场景预设 "+string(name), updated.Revision)); err != nil {
		return domain.GradeSession{}, err
	}
	return updated, nil
}

func (s *Service) UpdateGrade(clipID string, patch domain.GradePatch, actor string) (domain.GradeSession, error) {
	clip, err := s.store.Clip(clipID)
	if err != nil {
		return domain.GradeSession{}, err
	}
	if !clip.AvailableAt(s.clock.Now()) {
		return domain.GradeSession{}, domain.ErrExpired
	}
	current, err := s.store.Grade(clipID)
	if err != nil {
		return domain.GradeSession{}, err
	}
	if s.hook != nil {
		s.hook("read")
	}
	updated := current.Apply(patch, s.clock.Now())
	if err := updated.Validate(); err != nil {
		return domain.GradeSession{}, err
	}
	if s.hook != nil {
		s.hook("write")
	}
	if err := s.store.SaveGrade(updated); err != nil {
		return domain.GradeSession{}, err
	}
	if err := s.store.SaveEvent(s.event("grade.updated", clipID, actor, "更新调色微调参数", updated.Revision)); err != nil {
		return domain.GradeSession{}, err
	}
	return updated, nil
}

func (s *Service) UpdateGradeChecked(clipID string, expected uint64, patch domain.GradePatch, actor string) (domain.GradeSession, error) {
	clip, err := s.store.Clip(clipID)
	if err != nil {
		return domain.GradeSession{}, err
	}
	if !clip.AvailableAt(s.clock.Now()) {
		return domain.GradeSession{}, domain.ErrExpired
	}
	current, err := s.store.Grade(clipID)
	if err != nil {
		return domain.GradeSession{}, err
	}
	if current.Revision != expected {
		return domain.GradeSession{}, domain.ErrConflict
	}
	updated := current.Apply(patch, s.clock.Now())
	if err := s.store.SaveGradeIfRevision(updated, expected); err != nil {
		return domain.GradeSession{}, err
	}
	if err := s.store.SaveEvent(s.event("grade.updated", clipID, actor, "带修订校验更新调色参数", updated.Revision)); err != nil {
		return domain.GradeSession{}, err
	}
	return updated, nil
}

func (s *Service) RefreshPreview(clipID string, timestamps []int64, prefix string, actor string) ([]domain.PreviewFrame, error) {
	clip, err := s.store.Clip(clipID)
	if err != nil {
		return nil, err
	}
	if !clip.AvailableAt(s.clock.Now()) {
		return nil, domain.ErrExpired
	}
	grade, err := s.store.Grade(clipID)
	if err != nil {
		return nil, err
	}
	frames, err := s.planner.At(clip, grade, timestamps, prefix, s.clock.Now())
	if err != nil {
		return nil, err
	}
	for index := range frames {
		frame := frames[index].Start()
		output := filepath.Join(s.uploads.Root, "previews", frame.ID+".jpg")
		frame = frame.Complete(output, s.clock.Now())
		if err := s.store.SavePreview(frame); err != nil {
			return nil, err
		}
		frames[index] = frame
	}
	event := s.event("preview.refreshed", clipID, actor, fmt.Sprintf("刷新 %d 个预览帧", len(frames)), grade.Revision)
	if err := s.store.SaveEvent(event); err != nil {
		return nil, err
	}
	return frames, nil
}

func (s *Service) PreviewCommands(clipID string, frames []domain.PreviewFrame) ([]media.Command, error) {
	clip, err := s.store.Clip(clipID)
	if err != nil {
		return nil, err
	}
	grade, err := s.store.Grade(clipID)
	if err != nil {
		return nil, err
	}
	input, err := s.uploads.Destination(clip.ID, clip.SourceName)
	if err != nil {
		return nil, err
	}
	return media.BatchPreviewCommands(input, filepath.Join(s.uploads.Root, "previews"), frames, grade)
}

func (s *Service) Clip(id string) (domain.ClipSummary, error) {
	clip, err := s.store.Clip(id)
	if err != nil {
		return domain.ClipSummary{}, err
	}
	grade, err := s.store.Grade(id)
	if err != nil {
		return domain.ClipSummary{}, err
	}
	frames, err := s.store.Previews(id)
	if err != nil {
		return domain.ClipSummary{}, err
	}
	return domain.SummarizeClip(clip, grade, frames), nil
}

func (s *Service) ExpireClip(id, actor string) error {
	clip, err := s.store.Clip(id)
	if err != nil {
		return err
	}
	if clip.State == domain.ClipDeleted {
		return nil
	}
	clip.State = domain.ClipExpired
	if err := s.store.SaveClip(clip); err != nil {
		return err
	}
	return s.store.SaveEvent(s.event("clip.expired", id, actor, "短片临时文件到期", 0))
}

func (s *Service) DeleteClip(id, actor string) error {
	clip, err := s.store.Clip(id)
	if err != nil {
		return err
	}
	path, pathErr := s.uploads.Destination(clip.ID, clip.SourceName)
	if pathErr == nil {
		if err := s.uploads.Remove(path); err != nil {
			return err
		}
	}
	clip.State = domain.ClipDeleted
	if err := s.store.SaveClip(clip); err != nil {
		return err
	}
	if err := s.store.DeletePreviews(id); err != nil {
		return err
	}
	return s.store.SaveEvent(s.event("clip.deleted", id, actor, "删除短片及预览临时文件", 0))
}

func (s *Service) Store() *storage.Store { return s.store }
