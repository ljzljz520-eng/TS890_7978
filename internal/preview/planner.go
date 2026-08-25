package preview

import (
	"errors"
	"fmt"
	"petcolor/internal/domain"
	"sort"
	"time"
)

type Planner struct {
	Width     int
	Height    int
	MaxFrames int
}

func NewPlanner(width, height, maxFrames int) (Planner, error) {
	if width < 120 || height < 90 {
		return Planner{}, errors.New("preview dimensions too small")
	}
	if maxFrames < 1 || maxFrames > 24 {
		return Planner{}, errors.New("frame limit outside range")
	}
	return Planner{Width: width, Height: height, MaxFrames: maxFrames}, nil
}

func (p Planner) At(clip domain.ClipAsset, grade domain.GradeSession, timestamps []int64, prefix string, now time.Time) ([]domain.PreviewFrame, error) {
	if len(timestamps) == 0 {
		return nil, errors.New("timestamps are required")
	}
	if len(timestamps) > p.MaxFrames {
		return nil, errors.New("too many preview frames")
	}
	sorted := append([]int64(nil), timestamps...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	frames := make([]domain.PreviewFrame, 0, len(sorted))
	last := int64(-1)
	for index, stamp := range sorted {
		if stamp == last {
			continue
		}
		request := domain.PreviewRequest{ID: fmt.Sprintf("%s-%03d", prefix, index+1), ClipID: clip.ID, Sequence: index + 1, TimestampMS: stamp, Width: p.Width, Height: p.Height, RequestedAt: now}
		frame, err := domain.NewPreviewFrame(request, clip, grade)
		if err != nil {
			return nil, err
		}
		frames = append(frames, frame)
		last = stamp
	}
	return frames, nil
}

func (p Planner) EvenlySpaced(clip domain.ClipAsset, grade domain.GradeSession, count int, prefix string, now time.Time) ([]domain.PreviewFrame, error) {
	if count < 1 || count > p.MaxFrames {
		return nil, errors.New("frame count outside range")
	}
	timestamps := make([]int64, 0, count)
	if count == 1 {
		timestamps = append(timestamps, clip.DurationMS/2)
	} else {
		margin := clip.DurationMS / 20
		span := clip.DurationMS - 2*margin
		for index := 0; index < count; index++ {
			timestamps = append(timestamps, margin+int64(index)*span/int64(count-1))
		}
	}
	return p.At(clip, grade, timestamps, prefix, now)
}

func SelectRefresh(previous []domain.PreviewFrame, revision uint64) []domain.PreviewFrame {
	out := []domain.PreviewFrame{}
	for _, frame := range previous {
		if frame.GradeRevision != revision || frame.Status == domain.PreviewFailed {
			out = append(out, frame)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Sequence < out[j].Sequence })
	return out
}

func CompletionRate(frames []domain.PreviewFrame) float64 {
	if len(frames) == 0 {
		return 0
	}
	ready := 0
	for _, frame := range frames {
		if frame.Status == domain.PreviewReady {
			ready++
		}
	}
	return float64(ready) / float64(len(frames))
}
