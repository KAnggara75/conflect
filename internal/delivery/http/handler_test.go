/*
 * Copyright (c) 2025 KAnggara75
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * See <https://www.gnu.org/licenses/gpl-3.0.html>.
 *
 * @author KAnggara75
 * @project conflect http
 */

package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/KAnggara75/conflect/internal/config"
	"github.com/KAnggara75/conflect/internal/repository"
	"github.com/KAnggara75/conflect/internal/service"
)

func TestNewServer(t *testing.T) {
	cfg := &config.Config{Port: "8080"}
	q := service.NewQueue(10)
	srv := NewServer(cfg, q, nil)

	if srv == nil || srv.cfg != cfg || srv.queue != q {
		t.Errorf("NewServer initialized incorrectly")
	}
}

func TestHandleWebhook_InvalidJSON(t *testing.T) {
	q := service.NewQueue(10)
	srv := &Server{
		queue: q,
	}

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString("invalid-json"))
	w := httptest.NewRecorder()

	srv.handleWebhook(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleWebhook_EnqueueSuccess(t *testing.T) {
	q := service.NewQueue(10)
	srv := &Server{
		queue: q,
	}

	payload := map[string]string{
		"ref":   "refs/heads/main",
		"after": "1234567890abcdef1234567890abcdef12345678",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	srv.handleWebhook(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if resp["status"] != "accepted" {
		t.Errorf("expected status 'accepted', got '%s'", resp["status"])
	}
	if resp["branch"] != "main" {
		t.Errorf("expected branch 'main', got '%s'", resp["branch"])
	}
}

func TestHandleWebhook_QueueFull(t *testing.T) {
	q := service.NewQueue(0) // Full queue (capacity 0)
	srv := &Server{
		queue: q,
	}

	payload := map[string]string{
		"ref":   "refs/heads/main",
		"after": "1234567890abcdef1234567890abcdef12345678",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	srv.handleWebhook(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}
}

func TestHandleWebhook_UpToDate(t *testing.T) {
	tmpDir := t.TempDir()
	mainDir := filepath.Join(tmpDir, "main")
	_ = os.MkdirAll(mainDir, 0755)

	cfg := &config.Config{RepoPath: tmpDir, DefaultBranch: "main"}
	repo := repository.NewGitRepo(tmpDir, "")
	cs := service.NewConfigServiceFromRepo(repo, cfg)
	q := service.NewQueue(10)

	srv := &Server{
		cfg:           cfg,
		queue:         q,
		configService: cs,
	}

	payload := map[string]string{
		"ref":   "refs/heads/main",
		"after": "",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	srv.handleWebhook(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, w.Code)
	}
}

func TestHandleConfig(t *testing.T) {
	tmpDir := t.TempDir()
	envDir := filepath.Join(tmpDir, "main", "production")
	_ = os.MkdirAll(envDir, 0755)
	_ = os.WriteFile(filepath.Join(envDir, "myapp-production.yaml"), []byte("key: value"), 0644)

	cfg := &config.Config{
		RepoPath:      tmpDir,
		DefaultBranch: "main",
	}
	repo := repository.NewGitRepo(tmpDir, "")
	cs := service.NewConfigServiceFromRepo(repo, cfg)
	srv := &Server{configService: cs}

	t.Run("Invalid path format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/myapp", nil)
		rec := httptest.NewRecorder()
		srv.handleConfig(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("Not found config", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/nonexistent/production/main", nil)
		rec := httptest.NewRecorder()
		srv.handleConfig(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})

	t.Run("Success load config", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/myapp/production/main", nil)
		rec := httptest.NewRecorder()
		srv.handleConfig(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestHealth(t *testing.T) {
	srv := &Server{
		cfg: &config.Config{},
	}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	srv.health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestServer_StartAndShutdown(t *testing.T) {
	cfg := &config.Config{Port: "0"} // free port
	q := service.NewQueue(10)
	srv := NewServer(cfg, q, nil)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	time.Sleep(100 * time.Millisecond)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("unexpected error on server shutdown: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("server start did not return after SIGINT")
	}
}
