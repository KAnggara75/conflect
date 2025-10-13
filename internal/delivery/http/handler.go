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
 * @author KAnggara75 on Mon 22/09/25 07.41
 * @project conflect http
 * https://github.com/KAnggara75/conflect/tree/main/internal/delivery/http
 */

package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/KAnggara75/conflect/internal/config"
	"github.com/KAnggara75/conflect/internal/delivery/http/middleware"
	"github.com/KAnggara75/conflect/internal/service"
)

type Server struct {
	cfg           *config.Config
	queue         *service.Queue
	configService *service.ConfigService
}

func NewServer(cfg *config.Config, q *service.Queue, cs *service.ConfigService) *Server {
	return &Server{
		cfg:           cfg,
		queue:         q,
		configService: cs,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)

	// Route yang butuh middleware
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/webhook", s.handleWebhook)
	protectedMux.HandleFunc("/", s.handleConfig)

	// Wrap middleware
	authCfg := middleware.AuthConfig{Token: s.cfg.Token}

	// Chain untuk endpoint yang dilindungi
	protectedHandler := middleware.Chain(
		protectedMux,
		middleware.Logging,
		middleware.RateLimitMiddleware(s.cfg.Limit, time.Minute),
		middleware.AuthMiddleware(authCfg),
	)

	// Gabungkan kedua mux
	rootMux := http.NewServeMux()
	rootMux.Handle("/health", mux)
	rootMux.Handle("/", protectedHandler)

	srv := &http.Server{
		Addr:    ":" + s.cfg.Port,
		Handler: rootMux,
	}

	log.Printf("🚀 Conflect server running at :%s", s.cfg.Port)
	return srv.ListenAndServe()

}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]interface{}{
		"status": "ok",
		"code":   http.StatusOK,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Printf("[%s-HEALTH] %s %s in %v", r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
}

func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Ref string `json:"ref"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	parts := strings.Split(payload.Ref, "/")
	branch := parts[len(parts)-1]

	s.queue.Enqueue(branch)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "accepted",
		"branch": branch,
	})
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, `{"error":"invalid path, expected /{app}/{env}/{label?}"}`, http.StatusBadRequest)
		return
	}

	appName := parts[0]
	env := parts[1]
	label := ""
	if len(parts) > 2 {
		label = parts[2]
	}

	resp := s.configService.LoadConfig(appName, env, label)

	w.Header().Set("Content-Type", "application/json")

	// kalau tidak ada property sources, return 404
	if len(resp.PropertySources) == 0 {
		w.WriteHeader(http.StatusNotFound)
		resp.Error = "config for " + appName + " with env " + env + " found"
	} else {
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(resp)
}
