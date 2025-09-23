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
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/KAnggara75/conflect/internal/config"
	"github.com/KAnggara75/conflect/internal/delivery/http/dto"
	"github.com/KAnggara75/conflect/internal/service"
)

type Server struct {
	cfg           *config.Config
	queue         *service.Queue
	configService *service.ConfigService
}

func NewServer(cfg *config.Config, q *service.Queue, cs *service.ConfigService) *Server {
	return &Server{cfg: cfg, queue: q, configService: cs}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", s.handleWebhook)
	mux.HandleFunc("/", s.handleConfig)
	return http.ListenAndServe(":"+s.cfg.Port, mux)
}

func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	// Here you’d verify GitHub signature, omitted for brevity
	s.queue.Enqueue()
	w.WriteHeader(200)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		fmt.Println("Failed to write response:", err)
		return
	}
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, `{"error":"invalid path, expected /{app}/{env}/{label?}"}`, http.StatusBadRequest)
		return
	}

	var propertySources []dto.PropertySource

	// commit hash (version). If error, leave empty string.
	version, _ := s.configService.GetCommitHash()

	appName := parts[0]
	env := parts[1]
	label := ""
	if len(parts) > 2 {
		label = parts[2]
	}

	prefixes := []string{appName, "application"}
	extensionList := []string{".yaml", ".yml", ".json", ".properties"}

	for _, prefix := range prefixes {
		// try with label if exists
		config, err := s.configService.LoadConfig(app, env, label)
		if err != nil && label != "" {
			// fallback: try without label
			config, err = s.configService.LoadConfig(app, env, "")
		}

		filename := fmt.Sprintf("%s-%s", prefix, env)
		if label != "" {
			filename = fmt.Sprintf("%s-%s", filename, *label)
		}

		for _, ext := range extensionList {
			data, err := s.configService.GetFile(env, filename+ext)
			if err != nil {
				continue
			}

			parsed, err := parseConfigFile(data, ext)
			if err != nil {
				// skip file if parse fails
				continue
			}

			// name must be "dev/pakaiwa-dev.yaml" style (relative path)
			rel := filepath.ToSlash(filepath.Join(env, filename+ext))

			propertySources = append(propertySources, dto.PropertySource{
				Name:   rel,
				Source: parsed,
			})
		}
	}

	resp := dto.ConfigResponse{
		Name:            appName,
		Profiles:        []string{env},
		Label:           label,
		Version:         version,
		PropertySources: propertySources,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
