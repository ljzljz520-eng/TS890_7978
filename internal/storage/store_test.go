package storage

import (
	"path/filepath"
	"petcolor/internal/domain"
	"testing"
	"time"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "reopen.db")
	now := time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	clip := domain.ClipAsset{ID: "persisted", PetKind: domain.PetCat, SourceName: "cat.mp4", MediaType: "video/mp4", SizeBytes: 100, DurationMS: 1000, Width: 640, Height: 360, State: domain.ClipReady, CreatedAt: now, ExpiresAt: now.Add(time.Hour), Checksum: "abcdefgh"}
	grade := domain.DefaultGrade(clip.ID, now)
	frame := domain.PreviewFrame{ID: "frame", ClipID: clip.ID, Sequence: 1, TimestampMS: 500, GradeRevision: 1, Status: domain.PreviewReady, Width: 320, Height: 180, CreatedAt: now}
	event, _ := domain.NewAuditEvent("event", "stored", clip.ID, "test", "saved", 1, now)
	if err := first.SaveClip(clip); err != nil {
		t.Fatal(err)
	}
	if err := first.SaveGrade(grade); err != nil {
		t.Fatal(err)
	}
	if err := first.SavePreview(frame); err != nil {
		t.Fatal(err)
	}
	if err := first.SaveEvent(event); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	if got, err := second.Clip(clip.ID); err != nil || got.SourceName != "cat.mp4" {
		t.Fatalf("clip=%+v err=%v", got, err)
	}
	if got, err := second.Grade(clip.ID); err != nil || got.Revision != 1 {
		t.Fatalf("grade=%+v err=%v", got, err)
	}
	previews, err := second.Previews(clip.ID)
	if err != nil || len(previews) != 1 {
		t.Fatalf("previews=%v err=%v", previews, err)
	}
	events, err := second.Events()
	if err != nil || len(events) != 1 {
		t.Fatalf("events=%v err=%v", events, err)
	}
}
