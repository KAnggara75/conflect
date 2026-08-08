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
 * @author KAnggara75 on Sat 08/08/26
 * @project conflect http
 */
package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KAnggara75/conflect/internal/config"
	"github.com/KAnggara75/conflect/internal/service"
)

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
