package domain

import (
	"testing"
	"time"
)

func TestClipAndGradeLifecycle(t *testing.T) {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	clip, err := NewClipAsset(UploadRequest{ID: "c1", PetKind: PetDog, SourceName: "dog.mov", MediaType: "video/quicktime", SizeBytes: 1000, DurationMS: 5000, Width: 1080, Height: 1920, Checksum: "12345678", CreatedAt: now, TTL: time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	if clip.Orientation() != "portrait" || !clip.AvailableAt(now) {
		t.Fatalf("clip=%+v", clip)
	}
	grade := DefaultGrade(clip.ID, now)
	value := 44
	grade = grade.Apply(GradePatch{Exposure: &value}, now.Add(time.Minute))
	if grade.Exposure != 44 || grade.Revision != 2 {
		t.Fatalf("grade=%+v", grade)
	}
}

func TestFiltersAndSummary(t *testing.T) {
	now := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	items := []ClipAsset{{ID: "a", PetKind: PetCat, SourceName: "one.mp4", State: ClipReady, CreatedAt: now}, {ID: "b", PetKind: PetDog, SourceName: "two.mp4", State: ClipReady, CreatedAt: now.Add(time.Minute)}}
	filtered := FilterClips(items, ClipFilter{PetKind: PetDog})
	if len(filtered) != 1 || filtered[0].ID != "b" {
		t.Fatalf("filtered=%v", filtered)
	}
	frames := []PreviewFrame{{ID: "p1", Sequence: 1}, {ID: "p2", Sequence: 2}}
	summary := SummarizeClip(items[0], GradeSession{}, frames)
	if summary.PreviewCount != 2 || summary.LatestPreview.ID != "p2" {
		t.Fatalf("summary=%+v", summary)
	}
}
