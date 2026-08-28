package web

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"petcolor/internal/media"
	"petcolor/internal/preview"
	"petcolor/internal/service"
	"petcolor/internal/storage"
	"strings"
	"testing"
	"time"
)

func TestServerPageAndHealth(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "web.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	uploads, _ := media.NewUploadManager(filepath.Join(t.TempDir(), "media"), 1000)
	planner, _ := preview.NewPlanner(320, 180, 4)
	application, _ := service.New(store, uploads, planner, service.FixedClock{Value: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)})
	server, err := NewServer(application)
	if err != nil {
		t.Fatal(err)
	}
	page := httptest.NewRecorder()
	server.Handler().ServeHTTP(page, httptest.NewRequest(http.MethodGet, "/", nil))
	if page.Code != http.StatusOK || !strings.Contains(page.Body.String(), "宠物短片调色台") {
		t.Fatalf("page=%d", page.Code)
	}
	health := httptest.NewRecorder()
	server.Handler().ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/health", nil))
	if health.Code != http.StatusOK {
		t.Fatalf("health=%d body=%s", health.Code, health.Body.String())
	}
}
