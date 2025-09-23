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
 * @author KAnggara75 on Mon 22/09/25 07.40
 * @project conflect service
 * https://github.com/KAnggara75/conflect/tree/main/internal/service
 */

package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/KAnggara75/conflect/internal/config"
	"github.com/KAnggara75/conflect/internal/delivery/http/dto"
	"github.com/KAnggara75/conflect/internal/helper"
	"github.com/KAnggara75/conflect/internal/repository"
)

type ConfigService struct {
	repo *repository.GitRepo
	cfg  *config.Config
}

func NewConfigService(cfg *config.Config, q *Queue) *ConfigService {
	repo := repository.NewGitRepo(cfg.RepoPath, cfg.RepoURL)

	if err := repo.EnsureCloned(); err != nil {
		log.Fatalf("failed to clone repo: %v", err)
	}
	return &ConfigService{repo: repo, cfg: cfg}
}

func (c *ConfigService) UpdateRepo() error {
	log.Println("Pulling latest config...")
	return c.repo.Pull()
}

func (c *ConfigService) GetFile(env string, filename string) ([]byte, error) {
	full := filepath.Join(c.cfg.RepoPath, env, filename)
	log.Println(os.ReadFile(full))
	return os.ReadFile(full)
}

func (c *ConfigService) GetCommitHash() (string, error) {
	return c.repo.GetCommitHash()
}

func (c *ConfigService) LoadConfig(app, env, label string) (*ConfigResponse, error) {
	basePath := c.repo.Path
	var dir string
	if label != "" {
		dir = filepath.Join(basePath, label, env)
	} else {
		dir = filepath.Join(basePath, env)
	}

	// check dir exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory not found: %s", dir)
	}

	files := []string{
		fmt.Sprintf("%s-%s.yaml", app, env),
		fmt.Sprintf("application-%s.yaml", env),
		fmt.Sprintf("%s-%s.json", app, env),
		fmt.Sprintf("application-%s.json", env),
		fmt.Sprintf("%s-%s.properties", app, env),
		fmt.Sprintf("application-%s.properties", env),
	}

	var propertySources []dto.PropertySource

	for _, f := range files {
		fullPath := filepath.Join(dir, f)
		if _, err := os.Stat(fullPath); err == nil {
			props, err := helper.ParseFile(fullPath)
			if err != nil {
				return nil, err
			}
			propertySources = append(propertySources, dto.PropertySource{
				Name:   filepath.Join(filepath.Base(filepath.Dir(fullPath)), filepath.Base(fullPath)), // label/env/file
				Source: props,
			})
		}
	}

	hash, _ := c.repo.GetCommitHash()

	return &dto.ConfigResponse{
		Name:            app,
		Profiles:        []string{env},
		Label:           nullable(label),
		Version:         hash,
		PropertySources: propertySources,
	}, nil
}
