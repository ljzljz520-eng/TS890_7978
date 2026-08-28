package media

import (
	"bytes"
	"os"
	"path/filepath"
	"petcolor/internal/domain"
	"testing"
)

func TestUploadManagerSaveAndRemove(t *testing.T) {
	manager, err := NewUploadManager(t.TempDir(), 100)
	if err != nil {
		t.Fatal(err)
	}
	path, checksum, size, err := manager.Save("clip", "cat.mp4", bytes.NewBufferString("video-data"))
	if err != nil {
		t.Fatal(err)
	}
	if size != 10 || len(checksum) != 64 {
		t.Fatalf("size=%d checksum=%s", size, checksum)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
	if err := manager.Remove(path); err != nil {
		t.Fatal(err)
	}
}

func TestFFmpegCommands(t *testing.T) {
	grade := domain.GradeSession{ClipID: "c", Preset: domain.PresetOutdoor, Exposure: 10, Saturation: 20, Revision: 2}
	frame := domain.PreviewFrame{ID: "p", TimestampMS: 1250, Width: 640, Height: 360}
	command, err := PreviewCommand("input.mp4", filepath.Join("out", "p.jpg"), frame, grade)
	if err != nil {
		t.Fatal(err)
	}
	if command.Program != "ffmpeg" || len(command.Args) < 10 {
		t.Fatalf("command=%+v", command)
	}
}
