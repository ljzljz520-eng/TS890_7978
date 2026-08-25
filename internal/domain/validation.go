package domain

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("revision conflict")
	ErrExpired  = errors.New("clip expired")
)

func (r UploadRequest) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("clip id is required")
	}
	if r.PetKind != PetCat && r.PetKind != PetDog {
		return errors.New("pet kind must be cat or dog")
	}
	if strings.TrimSpace(r.SourceName) == "" {
		return errors.New("source name is required")
	}
	ext := strings.ToLower(filepath.Ext(r.SourceName))
	if ext != ".mp4" && ext != ".mov" && ext != ".webm" {
		return fmt.Errorf("unsupported video extension %q", ext)
	}
	if r.MediaType != "video/mp4" && r.MediaType != "video/quicktime" && r.MediaType != "video/webm" {
		return errors.New("unsupported media type")
	}
	if r.SizeBytes <= 0 || r.SizeBytes > 512*1024*1024 {
		return errors.New("video size is outside limits")
	}
	if r.DurationMS < 250 || r.DurationMS > 180000 {
		return errors.New("duration is outside limits")
	}
	if r.Width < 160 || r.Height < 120 {
		return errors.New("video dimensions are too small")
	}
	if r.Width > 7680 || r.Height > 4320 {
		return errors.New("video dimensions are too large")
	}
	if len(r.Checksum) < 8 {
		return errors.New("checksum is too short")
	}
	if r.CreatedAt.IsZero() {
		return errors.New("created time is required")
	}
	if r.TTL < time.Minute || r.TTL > 24*time.Hour {
		return errors.New("ttl is outside limits")
	}
	return nil
}

func NewClipAsset(r UploadRequest) (ClipAsset, error) {
	if err := r.Validate(); err != nil {
		return ClipAsset{}, err
	}
	return ClipAsset{ID: r.ID, PetKind: r.PetKind, SourceName: filepath.Base(r.SourceName), MediaType: r.MediaType, SizeBytes: r.SizeBytes, DurationMS: r.DurationMS, Width: r.Width, Height: r.Height, State: ClipReady, CreatedAt: r.CreatedAt.UTC(), ExpiresAt: r.CreatedAt.Add(r.TTL).UTC(), Checksum: r.Checksum}, nil
}

func (c ClipAsset) Validate() error {
	if c.ID == "" {
		return errors.New("clip id is required")
	}
	if c.State != ClipPending && c.State != ClipReady && c.State != ClipExpired && c.State != ClipDeleted {
		return errors.New("invalid clip state")
	}
	if c.ExpiresAt.Before(c.CreatedAt) {
		return errors.New("expiry precedes creation")
	}
	if c.DurationMS <= 0 {
		return errors.New("duration must be positive")
	}
	return nil
}

func (c ClipAsset) AvailableAt(now time.Time) bool {
	return c.State == ClipReady && now.Before(c.ExpiresAt)
}

func (c ClipAsset) AspectRatio() float64 {
	if c.Height == 0 {
		return 0
	}
	return float64(c.Width) / float64(c.Height)
}

func (c ClipAsset) Orientation() string {
	if c.Width > c.Height {
		return "landscape"
	}
	if c.Width < c.Height {
		return "portrait"
	}
	return "square"
}

func DefaultGrade(clipID string, now time.Time) GradeSession {
	return GradeSession{ClipID: clipID, Preset: PresetIndoor, Exposure: 0, Saturation: 0, Revision: 1, UpdatedAt: now.UTC()}
}

func (g GradeSession) Validate() error {
	if g.ClipID == "" {
		return errors.New("clip id is required")
	}
	if g.Preset != PresetIndoor && g.Preset != PresetOutdoor && g.Preset != PresetNight {
		return errors.New("invalid preset")
	}
	values := []int{g.Exposure, g.Saturation, g.Temperature, g.Contrast, g.Highlights, g.Shadows}
	for _, value := range values {
		if value < -100 || value > 100 {
			return errors.New("grade value outside range")
		}
	}
	if g.Revision == 0 {
		return errors.New("revision must be positive")
	}
	return nil
}

func clamp(value, low, high int) int {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}

func (g GradeSession) Apply(p GradePatch, now time.Time) GradeSession {
	if p.Preset != nil {
		g.Preset = *p.Preset
	}
	if p.Exposure != nil {
		g.Exposure = clamp(*p.Exposure, -100, 100)
	}
	if p.Saturation != nil {
		g.Saturation = clamp(*p.Saturation, -100, 100)
	}
	if p.Temperature != nil {
		g.Temperature = clamp(*p.Temperature, -100, 100)
	}
	if p.Contrast != nil {
		g.Contrast = clamp(*p.Contrast, -100, 100)
	}
	if p.Highlights != nil {
		g.Highlights = clamp(*p.Highlights, -100, 100)
	}
	if p.Shadows != nil {
		g.Shadows = clamp(*p.Shadows, -100, 100)
	}
	g.Revision++
	g.UpdatedAt = now.UTC()
	return g
}

func (p PreviewRequest) Validate(clip ClipAsset) error {
	if p.ID == "" || p.ClipID == "" {
		return errors.New("preview identity is required")
	}
	if p.ClipID != clip.ID {
		return errors.New("preview clip mismatch")
	}
	if p.Sequence < 1 {
		return errors.New("sequence must be positive")
	}
	if p.TimestampMS < 0 || p.TimestampMS > clip.DurationMS {
		return errors.New("preview timestamp outside clip")
	}
	if p.Width < 120 || p.Height < 90 || p.Width > 1920 || p.Height > 1080 {
		return errors.New("preview dimensions outside limits")
	}
	if p.RequestedAt.IsZero() {
		return errors.New("requested time is required")
	}
	return nil
}

func NewPreviewFrame(r PreviewRequest, clip ClipAsset, grade GradeSession) (PreviewFrame, error) {
	if err := r.Validate(clip); err != nil {
		return PreviewFrame{}, err
	}
	return PreviewFrame{ID: r.ID, ClipID: r.ClipID, Sequence: r.Sequence, TimestampMS: r.TimestampMS, GradeRevision: grade.Revision, Status: PreviewQueued, Width: r.Width, Height: r.Height, CreatedAt: r.RequestedAt.UTC()}, nil
}

func (p PreviewFrame) Start() PreviewFrame { p.Status = PreviewRendering; p.Error = ""; return p }
func (p PreviewFrame) Complete(path string, now time.Time) PreviewFrame {
	p.Status = PreviewReady
	p.OutputPath = path
	p.CompletedAt = now.UTC()
	p.Error = ""
	return p
}
func (p PreviewFrame) Fail(message string, now time.Time) PreviewFrame {
	p.Status = PreviewFailed
	p.CompletedAt = now.UTC()
	p.Error = message
	return p
}

func NewAuditEvent(id, kind, target, actor, summary string, revision uint64, at time.Time) (AuditEvent, error) {
	if id == "" || kind == "" || target == "" {
		return AuditEvent{}, errors.New("audit identity is required")
	}
	if actor == "" {
		actor = "system"
	}
	if len(summary) > 240 {
		return AuditEvent{}, errors.New("audit summary too long")
	}
	return AuditEvent{ID: id, Kind: kind, TargetID: target, Actor: actor, Summary: summary, Revision: revision, OccurredAt: at.UTC()}, nil
}
