package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"petcolor/internal/domain"
	"petcolor/internal/grading"
	"petcolor/internal/query"
	"petcolor/internal/service"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	application *service.Service
	catalog     *query.Catalog
	mux         *http.ServeMux
}

func NewServer(application *service.Service) (*Server, error) {
	if application == nil {
		return nil, errors.New("application service is required")
	}
	server := &Server{
		application: application,
		catalog:     query.NewCatalog(application),
		mux:         http.NewServeMux(),
	}
	server.routes()
	return server, nil
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /", s.handlePage)
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /api/presets", s.handlePresets)
	s.mux.HandleFunc("GET /api/dashboard", s.handleDashboard)
	s.mux.HandleFunc("GET /api/clips", s.handleClips)
	s.mux.HandleFunc("POST /api/clips", s.handleCreateClip)
	s.mux.HandleFunc("GET /api/clips/{id}", s.handleClip)
	s.mux.HandleFunc("DELETE /api/clips/{id}", s.handleDeleteClip)
	s.mux.HandleFunc("POST /api/clips/{id}/preset", s.handlePreset)
	s.mux.HandleFunc("PATCH /api/clips/{id}/grade", s.handleGrade)
	s.mux.HandleFunc("POST /api/clips/{id}/previews", s.handlePreview)
	s.mux.HandleFunc("GET /api/clips/{id}/activity", s.handleActivity)
}

func (s *Server) Handler() http.Handler {
	return s.recoverPanic(s.requestID(s.mux))
}

func (s *Server) handlePage(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Header().Set("Cache-Control", "no-store")
	writer.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(writer, pageHTML)
}

func (s *Server) handleHealth(writer http.ResponseWriter, request *http.Request) {
	health := s.application.Store().Health()
	status := http.StatusOK
	if !health.Writable {
		status = http.StatusServiceUnavailable
	}
	writeJSON(writer, status, health)
}

func (s *Server) handlePresets(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, map[string]any{"items": grading.List()})
}

func (s *Server) handleDashboard(writer http.ResponseWriter, request *http.Request) {
	dashboard, err := s.application.Dashboard(8)
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, dashboard)
}

func (s *Server) handleClips(writer http.ResponseWriter, request *http.Request) {
	limit, err := parseLimit(request.URL.Query().Get("limit"), 50)
	if err != nil {
		writeError(writer, err)
		return
	}
	search := query.SearchRequest{
		Text:   request.URL.Query().Get("q"),
		Pet:    domain.PetKind(request.URL.Query().Get("pet")),
		State:  domain.ClipState(request.URL.Query().Get("state")),
		Preset: domain.PresetName(request.URL.Query().Get("preset")),
		Limit:  limit,
	}
	result, err := s.catalog.Search(search)
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, result)
}

type createClipPayload struct {
	ID         string         `json:"id"`
	PetKind    domain.PetKind `json:"pet_kind"`
	SourceName string         `json:"source_name"`
	MediaType  string         `json:"media_type"`
	SizeBytes  int64          `json:"size_bytes"`
	DurationMS int64          `json:"duration_ms"`
	Width      int            `json:"width"`
	Height     int            `json:"height"`
	Checksum   string         `json:"checksum"`
	TTLMinutes int            `json:"ttl_minutes"`
}

func (s *Server) handleCreateClip(writer http.ResponseWriter, request *http.Request) {
	var payload createClipPayload
	if err := decodeJSON(request, &payload); err != nil {
		writeError(writer, err)
		return
	}
	if payload.TTLMinutes == 0 {
		payload.TTLMinutes = 60
	}
	upload := domain.UploadRequest{
		ID:         payload.ID,
		PetKind:    payload.PetKind,
		SourceName: payload.SourceName,
		MediaType:  payload.MediaType,
		SizeBytes:  payload.SizeBytes,
		DurationMS: payload.DurationMS,
		Width:      payload.Width,
		Height:     payload.Height,
		Checksum:   payload.Checksum,
		CreatedAt:  time.Now().UTC(),
		TTL:        time.Duration(payload.TTLMinutes) * time.Minute,
	}
	summary, err := s.application.RegisterUpload(upload, actor(request))
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, summary)
}

func (s *Server) handleClip(writer http.ResponseWriter, request *http.Request) {
	summary, err := s.application.Clip(request.PathValue("id"))
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, summary)
}

func (s *Server) handleDeleteClip(writer http.ResponseWriter, request *http.Request) {
	if err := s.application.DeleteClip(request.PathValue("id"), actor(request)); err != nil {
		writeError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

type presetPayload struct {
	Name domain.PresetName `json:"name"`
}

func (s *Server) handlePreset(writer http.ResponseWriter, request *http.Request) {
	var payload presetPayload
	if err := decodeJSON(request, &payload); err != nil {
		writeError(writer, err)
		return
	}
	grade, err := s.application.ApplyPreset(request.PathValue("id"), payload.Name, actor(request))
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, grade)
}

type gradePayload struct {
	Revision    uint64 `json:"revision"`
	Exposure    *int   `json:"exposure"`
	Saturation  *int   `json:"saturation"`
	Temperature *int   `json:"temperature"`
	Contrast    *int   `json:"contrast"`
	Highlights  *int   `json:"highlights"`
	Shadows     *int   `json:"shadows"`
}

func (s *Server) handleGrade(writer http.ResponseWriter, request *http.Request) {
	var payload gradePayload
	if err := decodeJSON(request, &payload); err != nil {
		writeError(writer, err)
		return
	}
	patch := domain.GradePatch{
		Exposure:    payload.Exposure,
		Saturation:  payload.Saturation,
		Temperature: payload.Temperature,
		Contrast:    payload.Contrast,
		Highlights:  payload.Highlights,
		Shadows:     payload.Shadows,
	}
	grade, err := s.application.UpdateGradeChecked(request.PathValue("id"), payload.Revision, patch, actor(request))
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, grade)
}

type previewPayload struct {
	Timestamps []int64 `json:"timestamps_ms"`
	Prefix     string  `json:"prefix"`
}

func (s *Server) handlePreview(writer http.ResponseWriter, request *http.Request) {
	var payload previewPayload
	if err := decodeJSON(request, &payload); err != nil {
		writeError(writer, err)
		return
	}
	if payload.Prefix == "" {
		payload.Prefix = "preview"
	}
	frames, err := s.application.RefreshPreview(request.PathValue("id"), payload.Timestamps, payload.Prefix, actor(request))
	if err != nil {
		writeError(writer, err)
		return
	}
	commands, err := s.application.PreviewCommands(request.PathValue("id"), frames)
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, map[string]any{"frames": frames, "commands": commands})
}

func (s *Server) handleActivity(writer http.ResponseWriter, request *http.Request) {
	limit, err := parseLimit(request.URL.Query().Get("limit"), 100)
	if err != nil {
		writeError(writer, err)
		return
	}
	activity, err := s.catalog.Activity(request.PathValue("id"), limit)
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, activity)
}

func parseLimit(raw string, fallback int) (int, error) {
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 || value > 500 {
		return 0, errors.New("limit must be between 1 and 500")
	}
	return value, nil
}

func decodeJSON(request *http.Request, destination any) error {
	if !strings.HasPrefix(request.Header.Get("Content-Type"), "application/json") {
		return errors.New("content type must be application/json")
	}
	decoder := json.NewDecoder(io.LimitReader(request.Body, 1024*1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("decode request: %w", err)
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return errors.New("request body must contain one json value")
	}
	return nil
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func writeError(writer http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, domain.ErrNotFound) {
		status = http.StatusNotFound
	}
	if errors.Is(err, domain.ErrConflict) {
		status = http.StatusConflict
	}
	if errors.Is(err, domain.ErrExpired) {
		status = http.StatusGone
	}
	writeJSON(writer, status, map[string]string{"error": err.Error()})
}

func actor(request *http.Request) string {
	value := strings.TrimSpace(request.Header.Get("X-Operator"))
	if value == "" {
		return "web-user"
	}
	if len(value) > 64 {
		return value[:64]
	}
	return value
}
