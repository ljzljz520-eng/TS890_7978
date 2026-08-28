package cleanup

import (
	"path/filepath"
	"petcolor/internal/domain"
	"petcolor/internal/media"
	"petcolor/internal/preview"
	"petcolor/internal/service"
	"petcolor/internal/storage"
	"testing"
	"time"
)

func TestRunnerExpiresAndDeletes(t *testing.T) {
	created := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	store, err := storage.Open(filepath.Join(t.TempDir(), "cleanup.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	uploads, _ := media.NewUploadManager(filepath.Join(t.TempDir(), "media"), 1000)
	planner, _ := preview.NewPlanner(320, 180, 4)
	application, _ := service.New(store, uploads, planner, service.FixedClock{Value: created})
	_, err = application.RegisterUpload(domain.UploadRequest{ID: "old", PetKind: domain.PetCat, SourceName: "cat.mp4", MediaType: "video/mp4", SizeBytes: 100, DurationMS: 1000, Width: 640, Height: 360, Checksum: "abcdefgh", CreatedAt: created, TTL: time.Minute}, "test")
	if err != nil {
		t.Fatal(err)
	}
	runner := NewRunner(application, Policy{DeleteAfter: time.Minute, BatchSize: 10})
	report := runner.Run(created.Add(3 * time.Minute))
	if !report.Successful() || len(report.Deleted) != 1 {
		t.Fatalf("report=%+v", report)
	}
}
