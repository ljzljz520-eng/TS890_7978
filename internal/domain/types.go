package domain

import "time"

type PetKind string

const (
	PetCat PetKind = "cat"
	PetDog PetKind = "dog"
)

type ClipState string

const (
	ClipPending ClipState = "pending"
	ClipReady   ClipState = "ready"
	ClipExpired ClipState = "expired"
	ClipDeleted ClipState = "deleted"
)

type ClipAsset struct {
	ID         string    `json:"id"`
	PetKind    PetKind   `json:"pet_kind"`
	SourceName string    `json:"source_name"`
	MediaType  string    `json:"media_type"`
	SizeBytes  int64     `json:"size_bytes"`
	DurationMS int64     `json:"duration_ms"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	State      ClipState `json:"state"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	Checksum   string    `json:"checksum"`
}

type PresetName string

const (
	PresetIndoor  PresetName = "indoor"
	PresetOutdoor PresetName = "outdoor"
	PresetNight   PresetName = "night"
)

type GradeSession struct {
	ClipID      string     `json:"clip_id"`
	Preset      PresetName `json:"preset"`
	Exposure    int        `json:"exposure"`
	Saturation  int        `json:"saturation"`
	Temperature int        `json:"temperature"`
	Contrast    int        `json:"contrast"`
	Highlights  int        `json:"highlights"`
	Shadows     int        `json:"shadows"`
	Revision    uint64     `json:"revision"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PreviewStatus string

const (
	PreviewQueued    PreviewStatus = "queued"
	PreviewRendering PreviewStatus = "rendering"
	PreviewReady     PreviewStatus = "ready"
	PreviewFailed    PreviewStatus = "failed"
)

type PreviewFrame struct {
	ID            string        `json:"id"`
	ClipID        string        `json:"clip_id"`
	Sequence      int           `json:"sequence"`
	TimestampMS   int64         `json:"timestamp_ms"`
	GradeRevision uint64        `json:"grade_revision"`
	Status        PreviewStatus `json:"status"`
	OutputPath    string        `json:"output_path"`
	Width         int           `json:"width"`
	Height        int           `json:"height"`
	CreatedAt     time.Time     `json:"created_at"`
	CompletedAt   time.Time     `json:"completed_at"`
	Error         string        `json:"error,omitempty"`
}

type AuditEvent struct {
	ID         string    `json:"id"`
	Kind       string    `json:"kind"`
	TargetID   string    `json:"target_id"`
	Actor      string    `json:"actor"`
	Summary    string    `json:"summary"`
	Revision   uint64    `json:"revision"`
	OccurredAt time.Time `json:"occurred_at"`
}

type UploadRequest struct {
	ID         string
	PetKind    PetKind
	SourceName string
	MediaType  string
	SizeBytes  int64
	DurationMS int64
	Width      int
	Height     int
	Checksum   string
	CreatedAt  time.Time
	TTL        time.Duration
}

type GradePatch struct {
	Preset      *PresetName
	Exposure    *int
	Saturation  *int
	Temperature *int
	Contrast    *int
	Highlights  *int
	Shadows     *int
}

type PreviewRequest struct {
	ID          string
	ClipID      string
	Sequence    int
	TimestampMS int64
	Width       int
	Height      int
	RequestedAt time.Time
}

type ClipSummary struct {
	Clip          ClipAsset     `json:"clip"`
	Grade         GradeSession  `json:"grade"`
	PreviewCount  int           `json:"preview_count"`
	LatestPreview *PreviewFrame `json:"latest_preview,omitempty"`
}
