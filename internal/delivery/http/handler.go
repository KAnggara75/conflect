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
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/KAnggara75/conflect/internal/config"
	"github.com/KAnggara75/conflect/internal/delivery/http/middleware"
	"github.com/KAnggara75/conflect/internal/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests processed.",
		},
		[]string{"method", "path", "status"},
	)
	verifyFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "webhook_verification_failures_total",
			Help: "Number of failed webhook signature verifications.",
		},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(verifyFailures)
}

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

	// Wrap middleware
	authCfg := middleware.AuthConfig{Token: s.cfg.Token}

	// Route yang butuh middleware
	webhookMux := http.NewServeMux()
	webhookMux.HandleFunc("/webhook", s.handleWebhook)

	// Chain untuk endpoint yang dilindungi
	webhookHandler := middleware.Chain(
		webhookMux,
		middleware.Logging,
		middleware.RateLimitMiddleware(s.cfg.Limit, time.Minute),
		middleware.VerifySignature(authCfg),
	)

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/", s.handleConfig)

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
	rootMux.Handle("/metrics", promhttp.Handler())
	rootMux.Handle("/webhook", webhookHandler)
	rootMux.Handle("/", protectedHandler)

	srv := &http.Server{
		Addr:         ":" + s.cfg.Port,
		Handler:      rootMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for errors from server
	serverErrors := make(chan error, 1)

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Conflect server running at :%s", s.cfg.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	// Channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal or server error
	select {
	case err := <-serverErrors:
		return err
	case sig := <-quit:
		log.Printf("⚠️  Received signal %v, initiating graceful shutdown...", sig)
	}

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Graceful shutdown failed: %v, forcing exit", err)
		return srv.Close()
	}

	log.Println("✅ Server gracefully stopped")
	return nil
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

	w.Header().Set("Content-Type", "application/json")

	if !s.queue.Enqueue(branch) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "queue_full",
			"error":  "server busy, please retry later",
		})
		return
	}

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
		resp.Error = "config for " + appName + " with env " + env + " not found"
	} else {
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(resp)
}
