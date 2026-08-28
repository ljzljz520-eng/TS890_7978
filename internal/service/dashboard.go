package service

import (
	"petcolor/internal/domain"
	"sort"
)

type Dashboard struct {
	TotalClips    int                       `json:"total_clips"`
	ActiveClips   int                       `json:"active_clips"`
	ExpiredClips  int                       `json:"expired_clips"`
	PreviewFrames int                       `json:"preview_frames"`
	StorageBytes  int64                     `json:"storage_bytes"`
	ByPet         map[domain.PetKind]int    `json:"by_pet"`
	ByPreset      map[domain.PresetName]int `json:"by_preset"`
	Recent        []domain.ClipSummary      `json:"recent"`
}

func (s *Service) Dashboard(limit int) (Dashboard, error) {
	clips, err := s.store.Clips()
	if err != nil {
		return Dashboard{}, err
	}
	frames, err := s.store.Previews("")
	if err != nil {
		return Dashboard{}, err
	}
	dashboard := Dashboard{TotalClips: len(clips), PreviewFrames: len(frames), ByPet: map[domain.PetKind]int{}, ByPreset: map[domain.PresetName]int{}}
	summaries := []domain.ClipSummary{}
	for _, clip := range clips {
		dashboard.ByPet[clip.PetKind]++
		if clip.AvailableAt(s.clock.Now()) {
			dashboard.ActiveClips++
			dashboard.StorageBytes += clip.SizeBytes
		} else if clip.State == domain.ClipExpired {
			dashboard.ExpiredClips++
		}
		grade, gradeErr := s.store.Grade(clip.ID)
		if gradeErr != nil {
			continue
		}
		dashboard.ByPreset[grade.Preset]++
		clipFrames, frameErr := s.store.Previews(clip.ID)
		if frameErr != nil {
			return Dashboard{}, frameErr
		}
		summaries = append(summaries, domain.SummarizeClip(clip, grade, clipFrames))
	}
	sort.Slice(summaries, func(i, j int) bool { return summaries[i].Clip.CreatedAt.After(summaries[j].Clip.CreatedAt) })
	if limit > 0 && len(summaries) > limit {
		summaries = summaries[:limit]
	}
	dashboard.Recent = summaries
	return dashboard, nil
}

func (s *Service) ListClips(filter domain.ClipFilter) ([]domain.ClipAsset, error) {
	items, err := s.store.Clips()
	if err != nil {
		return nil, err
	}
	return domain.FilterClips(items, filter), nil
}
func (s *Service) ListEvents(filter domain.EventFilter) ([]domain.AuditEvent, error) {
	items, err := s.store.Events()
	if err != nil {
		return nil, err
	}
	return domain.FilterEvents(items, filter), nil
}
