package query

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

func TestCatalogSearchAndActivity(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	store, err := storage.Open(filepath.Join(t.TempDir(), "query.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	uploads, _ := media.NewUploadManager(filepath.Join(t.TempDir(), "media"), 1000)
	planner, _ := preview.NewPlanner(320, 180, 4)
	application, _ := service.New(store, uploads, planner, service.FixedClock{Value: now})
	_, err = application.RegisterUpload(domain.UploadRequest{ID: "cat-one", PetKind: domain.PetCat, SourceName: "window-cat.mp4", MediaType: "video/mp4", SizeBytes: 100, DurationMS: 1000, Width: 640, Height: 360, Checksum: "abcdefgh", CreatedAt: now, TTL: time.Hour}, "test")
	if err != nil {
		t.Fatal(err)
	}
	catalog := NewCatalog(application)
	result, err := catalog.Search(SearchRequest{Text: "window", Pet: domain.PetCat})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 {
		t.Fatalf("result=%+v", result)
	}
	activity, err := catalog.Activity("cat-one", 10)
	if err != nil {
		t.Fatal(err)
	}
	if activity.Counts["clip.uploaded"] != 1 {
		t.Fatalf("activity=%+v", activity)
	}
}
