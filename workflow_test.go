package petcolor_test

import (
	"path/filepath"
	"petcolor/internal/domain"
	"petcolor/internal/media"
	"petcolor/internal/preview"
	"petcolor/internal/service"
	"petcolor/internal/storage"
	"sync"
	"testing"
	"time"
)

func workflowService(t *testing.T) (*service.Service, *storage.Store, time.Time) {
	t.Helper()
	now := time.Date(2026, 8, 26, 9, 30, 0, 0, time.UTC)
	store, err := storage.Open(filepath.Join(t.TempDir(), "workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	uploads, err := media.NewUploadManager(filepath.Join(t.TempDir(), "media"), 1024*1024)
	if err != nil {
		t.Fatal(err)
	}
	planner, err := preview.NewPlanner(640, 360, 8)
	if err != nil {
		t.Fatal(err)
	}
	application, err := service.New(store, uploads, planner, service.FixedClock{Value: now})
	if err != nil {
		t.Fatal(err)
	}
	return application, store, now
}

func registerWorkflowClip(t *testing.T, application *service.Service, now time.Time, id string) {
	t.Helper()
	_, err := application.RegisterUpload(domain.UploadRequest{ID: id, PetKind: domain.PetCat, SourceName: "cat-window.mp4", MediaType: "video/mp4", SizeBytes: 8096, DurationMS: 12000, Width: 1920, Height: 1080, Checksum: "abcdef012345", CreatedAt: now, TTL: 2 * time.Hour}, "tester")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUploadAndPrepareWorkflow(t *testing.T) {
	application, store, now := workflowService(t)
	defer store.Close()
	registerWorkflowClip(t, application, now, "clip-upload")
	summary, err := application.Clip("clip-upload")
	if err != nil {
		t.Fatal(err)
	}
	if summary.Clip.State != domain.ClipReady {
		t.Fatalf("state=%s", summary.Clip.State)
	}
	if summary.Grade.Revision != 1 {
		t.Fatalf("revision=%d", summary.Grade.Revision)
	}
	events, err := application.ListEvents(domain.EventFilter{TargetID: "clip-upload"})
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Kind != "clip.uploaded" {
		t.Fatalf("events=%v", events)
	}
}

func TestPreviewRefreshWorkflow(t *testing.T) {
	application, store, now := workflowService(t)
	defer store.Close()
	registerWorkflowClip(t, application, now, "clip-preview")
	grade, err := application.ApplyPreset("clip-preview", domain.PresetOutdoor, "tester")
	if err != nil {
		t.Fatal(err)
	}
	frames, err := application.RefreshPreview("clip-preview", []int64{1000, 6000, 11000}, "fast", "tester")
	if err != nil {
		t.Fatal(err)
	}
	if len(frames) != 3 {
		t.Fatalf("frames=%d", len(frames))
	}
	for _, frame := range frames {
		if frame.Status != domain.PreviewReady || frame.GradeRevision != grade.Revision {
			t.Fatalf("frame=%+v", frame)
		}
	}
	commands, err := application.PreviewCommands("clip-preview", frames)
	if err != nil {
		t.Fatal(err)
	}
	if len(commands) != 3 || commands[0].Program != "ffmpeg" {
		t.Fatalf("commands=%v", commands)
	}
}

func TestBusinessChain45(t *testing.T) {
	application, store, now := workflowService(t)
	defer store.Close()
	registerWorkflowClip(t, application, now, "clip-concurrent")
	barrier := sync.WaitGroup{}
	barrier.Add(2)
	release := make(chan struct{})
	var reads sync.Once
	count := 0
	var mu sync.Mutex
	application.SetConcurrencyHook(func(stage string) {
		if stage != "read" {
			return
		}
		mu.Lock()
		count++
		if count == 2 {
			reads.Do(func() { close(release) })
		}
		mu.Unlock()
		<-release
	})
	exposure := 24
	saturation := 31
	errorsFound := make(chan error, 2)
	go func() {
		defer barrier.Done()
		_, err := application.UpdateGrade("clip-concurrent", domain.GradePatch{Exposure: &exposure}, "editor-a")
		errorsFound <- err
	}()
	go func() {
		defer barrier.Done()
		_, err := application.UpdateGrade("clip-concurrent", domain.GradePatch{Saturation: &saturation}, "editor-b")
		errorsFound <- err
	}()
	barrier.Wait()
	close(errorsFound)
	for err := range errorsFound {
		if err != nil {
			t.Fatal(err)
		}
	}
	grade, err := store.Grade("clip-concurrent")
	if err != nil {
		t.Fatal(err)
	}
	if grade.Exposure != exposure || grade.Saturation != saturation {
		t.Fatalf("concurrent edits were not preserved: exposure=%d saturation=%d", grade.Exposure, grade.Saturation)
	}
}
