package domain

import (
	"sort"
	"strings"
	"time"
)

type ClipFilter struct {
	PetKind       PetKind
	State         ClipState
	CreatedAfter  time.Time
	ExpiresBefore time.Time
	Search        string
	Limit         int
}
type EventFilter struct {
	Kind     string
	TargetID string
	Since    time.Time
	Until    time.Time
	Limit    int
}

func FilterClips(items []ClipAsset, f ClipFilter) []ClipAsset {
	out := make([]ClipAsset, 0, len(items))
	needle := strings.ToLower(strings.TrimSpace(f.Search))
	for _, item := range items {
		if f.PetKind != "" && item.PetKind != f.PetKind {
			continue
		}
		if f.State != "" && item.State != f.State {
			continue
		}
		if !f.CreatedAfter.IsZero() && !item.CreatedAt.After(f.CreatedAfter) {
			continue
		}
		if !f.ExpiresBefore.IsZero() && !item.ExpiresAt.Before(f.ExpiresBefore) {
			continue
		}
		if needle != "" && !strings.Contains(strings.ToLower(item.SourceName), needle) {
			continue
		}
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	if f.Limit > 0 && len(out) > f.Limit {
		out = out[:f.Limit]
	}
	return out
}

func FilterEvents(items []AuditEvent, f EventFilter) []AuditEvent {
	out := make([]AuditEvent, 0, len(items))
	for _, item := range items {
		if f.Kind != "" && item.Kind != f.Kind {
			continue
		}
		if f.TargetID != "" && item.TargetID != f.TargetID {
			continue
		}
		if !f.Since.IsZero() && item.OccurredAt.Before(f.Since) {
			continue
		}
		if !f.Until.IsZero() && item.OccurredAt.After(f.Until) {
			continue
		}
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].OccurredAt.After(out[j].OccurredAt) })
	if f.Limit > 0 && len(out) > f.Limit {
		out = out[:f.Limit]
	}
	return out
}

func LatestPreview(items []PreviewFrame) *PreviewFrame {
	if len(items) == 0 {
		return nil
	}
	latest := items[0]
	for _, item := range items[1:] {
		if item.Sequence > latest.Sequence || (item.Sequence == latest.Sequence && item.CreatedAt.After(latest.CreatedAt)) {
			latest = item
		}
	}
	return &latest
}

func SummarizeClip(clip ClipAsset, grade GradeSession, previews []PreviewFrame) ClipSummary {
	return ClipSummary{Clip: clip, Grade: grade, PreviewCount: len(previews), LatestPreview: LatestPreview(previews)}
}

func GroupEventsByKind(items []AuditEvent) map[string]int {
	out := map[string]int{}
	for _, item := range items {
		out[item.Kind]++
	}
	return out
}

func ActiveStorageBytes(items []ClipAsset, now time.Time) int64 {
	var total int64
	for _, item := range items {
		if item.AvailableAt(now) {
			total += item.SizeBytes
		}
	}
	return total
}
