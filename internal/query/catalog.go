package query

import (
	"petcolor/internal/domain"
	"petcolor/internal/service"
	"sort"
	"strings"
	"time"
)

type Catalog struct{ service *service.Service }

func NewCatalog(application *service.Service) *Catalog { return &Catalog{service: application} }

type SearchRequest struct {
	Text         string
	Pet          domain.PetKind
	State        domain.ClipState
	Preset       domain.PresetName
	CreatedAfter time.Time
	Limit        int
}
type SearchResult struct {
	Items  []domain.ClipSummary `json:"items"`
	Total  int                  `json:"total"`
	Active int                  `json:"active"`
	Query  string               `json:"query"`
}

func (c *Catalog) Search(request SearchRequest) (SearchResult, error) {
	clips, err := c.service.ListClips(domain.ClipFilter{PetKind: request.Pet, State: request.State, CreatedAfter: request.CreatedAfter, Search: request.Text})
	if err != nil {
		return SearchResult{}, err
	}
	items := []domain.ClipSummary{}
	active := 0
	for _, clip := range clips {
		summary, err := c.service.Clip(clip.ID)
		if err != nil {
			return SearchResult{}, err
		}
		if request.Preset != "" && summary.Grade.Preset != request.Preset {
			continue
		}
		if summary.Clip.State == domain.ClipReady {
			active++
		}
		items = append(items, summary)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Clip.CreatedAt.After(items[j].Clip.CreatedAt) })
	total := len(items)
	if request.Limit > 0 && len(items) > request.Limit {
		items = items[:request.Limit]
	}
	return SearchResult{Items: items, Total: total, Active: active, Query: strings.TrimSpace(request.Text)}, nil
}

type Activity struct {
	Events []domain.AuditEvent `json:"events"`
	Counts map[string]int      `json:"counts"`
}

func (c *Catalog) Activity(target string, limit int) (Activity, error) {
	events, err := c.service.ListEvents(domain.EventFilter{TargetID: target, Limit: limit})
	if err != nil {
		return Activity{}, err
	}
	return Activity{Events: events, Counts: domain.GroupEventsByKind(events)}, nil
}

type ManagementView struct {
	Dashboard    service.Dashboard     `json:"dashboard"`
	ExpiringSoon []domain.ClipAsset    `json:"expiring_soon"`
	Failures     []domain.PreviewFrame `json:"preview_failures"`
}

func (c *Catalog) Management(now time.Time, horizon time.Duration) (ManagementView, error) {
	dashboard, err := c.service.Dashboard(8)
	if err != nil {
		return ManagementView{}, err
	}
	clips, err := c.service.ListClips(domain.ClipFilter{State: domain.ClipReady, ExpiresBefore: now.Add(horizon)})
	if err != nil {
		return ManagementView{}, err
	}
	all, err := c.service.Store().Previews("")
	if err != nil {
		return ManagementView{}, err
	}
	failures := []domain.PreviewFrame{}
	for _, frame := range all {
		if frame.Status == domain.PreviewFailed {
			failures = append(failures, frame)
		}
	}
	return ManagementView{Dashboard: dashboard, ExpiringSoon: clips, Failures: failures}, nil
}
