package cleanup

import (
	"petcolor/internal/domain"
	"petcolor/internal/service"
	"sort"
	"time"
)

type Report struct {
	Scanned  int               `json:"scanned"`
	Expired  []string          `json:"expired"`
	Deleted  []string          `json:"deleted"`
	Retained []string          `json:"retained"`
	Errors   map[string]string `json:"errors"`
}
type Policy struct {
	DeleteAfter time.Duration
	BatchSize   int
}
type Runner struct {
	service *service.Service
	policy  Policy
}

func NewRunner(application *service.Service, policy Policy) *Runner {
	if policy.BatchSize <= 0 {
		policy.BatchSize = 100
	}
	if policy.DeleteAfter < 0 {
		policy.DeleteAfter = 0
	}
	return &Runner{service: application, policy: policy}
}

func (r *Runner) Run(now time.Time) Report {
	report := Report{Errors: map[string]string{}}
	clips, err := r.service.ListClips(domain.ClipFilter{})
	if err != nil {
		report.Errors["list"] = err.Error()
		return report
	}
	sort.Slice(clips, func(i, j int) bool { return clips[i].ExpiresAt.Before(clips[j].ExpiresAt) })
	for _, clip := range clips {
		if report.Scanned >= r.policy.BatchSize {
			break
		}
		report.Scanned++
		if now.Before(clip.ExpiresAt) {
			report.Retained = append(report.Retained, clip.ID)
			continue
		}
		if clip.State == domain.ClipReady || clip.State == domain.ClipPending {
			if err := r.service.ExpireClip(clip.ID, "cleanup"); err != nil {
				report.Errors[clip.ID] = err.Error()
				continue
			}
			report.Expired = append(report.Expired, clip.ID)
			clip.State = domain.ClipExpired
		}
		if clip.State == domain.ClipExpired && !now.Before(clip.ExpiresAt.Add(r.policy.DeleteAfter)) {
			if err := r.service.DeleteClip(clip.ID, "cleanup"); err != nil {
				report.Errors[clip.ID] = err.Error()
				continue
			}
			report.Deleted = append(report.Deleted, clip.ID)
		}
	}
	return report
}

func (r Report) Successful() bool { return len(r.Errors) == 0 }
func (r Report) Changed() int     { return len(r.Expired) + len(r.Deleted) }
