package service

import (
	"path/filepath"
	"petcolor/internal/domain"
	"petcolor/internal/media"
	"petcolor/internal/preview"
	"petcolor/internal/storage"
	"testing"
	"time"
)

func TestCheckedUpdateRejectsStaleRevision(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	store, err := storage.Open(filepath.Join(t.TempDir(), "service.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	uploads, _ := media.NewUploadManager(filepath.Join(t.TempDir(), "media"), 1000)
	planner, _ := preview.NewPlanner(320, 180, 4)
	application, _ := New(store, uploads, planner, FixedClock{Value: now})
	_, err = application.RegisterUpload(domain.UploadRequest{ID: "c", PetKind: domain.PetDog, SourceName: "dog.mp4", MediaType: "video/mp4", SizeBytes: 100, DurationMS: 2000, Width: 640, Height: 360, Checksum: "abcdefgh", CreatedAt: now, TTL: time.Hour}, "test")
	if err != nil {
		t.Fatal(err)
	}
	value := 10
	if _, err := application.UpdateGradeChecked("c", 1, domain.GradePatch{Exposure: &value}, "one"); err != nil {
		t.Fatal(err)
	}
	if _, err := application.UpdateGradeChecked("c", 1, domain.GradePatch{Saturation: &value}, "two"); err != domain.ErrConflict {
		t.Fatalf("err=%v", err)
	}
}
