package preview

import (
	"petcolor/internal/domain"
	"testing"
	"time"
)

func TestPlannerCreatesOrderedFrames(t *testing.T) {
	planner, err := NewPlanner(640, 360, 6)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	clip := domain.ClipAsset{ID: "c", DurationMS: 10000}
	grade := domain.GradeSession{ClipID: "c", Revision: 4}
	frames, err := planner.At(clip, grade, []int64{9000, 1000, 5000}, "p", now)
	if err != nil {
		t.Fatal(err)
	}
	if len(frames) != 3 || frames[0].TimestampMS != 1000 || frames[2].GradeRevision != 4 {
		t.Fatalf("frames=%+v", frames)
	}
}

func TestCompletionRate(t *testing.T) {
	frames := []domain.PreviewFrame{{Status: domain.PreviewReady}, {Status: domain.PreviewFailed}, {Status: domain.PreviewReady}}
	if rate := CompletionRate(frames); rate < 0.66 || rate > 0.67 {
		t.Fatalf("rate=%f", rate)
	}
}
