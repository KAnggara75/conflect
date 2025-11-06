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
	"strings"

	"github.com/KAnggara75/conflect/internal/config"
	"github.com/KAnggara75/conflect/internal/delivery/http/dto"
	"github.com/KAnggara75/conflect/internal/errors"
	"github.com/KAnggara75/conflect/internal/helper"
	"github.com/KAnggara75/conflect/internal/repository"
)

type ConfigService struct {
	repo *repository.GitRepo
	cfg  *config.Config
}

func NewConfigService(cfg *config.Config) *ConfigService {
	repo := repository.NewGitRepo(cfg.RepoPath, cfg.RepoURL)
	err := repo.InitAllBranches()
	if err != nil {
		log.Fatalf("failed to clone repo: %v", err)
	}
	return &ConfigService{repo: repo, cfg: cfg}
}

func (c *ConfigService) UpdateRepo(branch string) error {
	log.Printf("Pulling latest config for branch %s...", branch)
	return c.repo.Pull(branch)
}

func (c *ConfigService) LoadConfig(appName, env, label string) *dto.ConfigResponse {

	response := &dto.ConfigResponse{
		Name:            appName,
		Profiles:        []string{env},
		PropertySources: []dto.PropertySource{}, // inisialisasi slice kosong
	}

	if label == "" {
		label = c.cfg.DefaultBranch
	}

	response.Label = label

	candidates, err := c.generateConfigCandidates(appName, env, label)
	if err != nil {
		log.Println(err)
		return response
	}

	data, err := c.findAndReadAllConfigs(label, env, candidates)
	if err != nil {
		log.Println(err)
		return response
	}
	response.PropertySources = data

	hash, err := c.repo.GetCommitHashFromBranch(label)
	if err == nil {
		response.Version = hash
	}

	return response
}

func (c *ConfigService) generateConfigCandidates(appName, env, label string) ([]string, error) {
	basePath := filepath.Join(c.repo.Path, label, env)

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir %s: %w", basePath, err)
	}

	var (
		appFiles         []string
		applicationFiles []string
		globalFiles      []string
	)

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		ext := filepath.Ext(name)

		switch ext {
		case ".yaml", ".yml", ".json", ".properties":
		default:
			continue
		}

		// {appName}-{env}.*
		if strings.HasPrefix(name, appName+"-"+env) {
			appFiles = append(appFiles, name)
			continue
		}

		// application-{env}.*
		if strings.HasPrefix(name, "application-"+env) {
			applicationFiles = append(applicationFiles, name)
			continue
		}

		// application.*
		if strings.HasPrefix(name, "application.") {
			globalFiles = append(globalFiles, name)
		}
	}

	var candidates []string
	candidates = append(candidates, appFiles...)
	candidates = append(candidates, applicationFiles...)
	candidates = append(candidates, globalFiles...)

	log.Printf("Read File candidates: %v", candidates)

	return candidates, nil
}

func (c *ConfigService) findAndReadAllConfigs(label, env string, candidates []string) ([]dto.PropertySource, error) {
	var sources []dto.PropertySource

	for _, candidate := range candidates {
		filePath := filepath.Join(c.repo.Path, label, env, candidate)

		data, err := os.ReadFile(filePath)
		if err != nil {
			if skip, fileErr := errors.ShouldSkipFile(candidate, err); skip {
				continue
			} else {
				return nil, fileErr
			}
		}

		ext := filepath.Ext(filePath)
		props, err := helper.ParseFile(data, ext)
		if err != nil {
			if skip, fileErr := errors.ShouldSkipFile(candidate, err); skip {
				continue
			} else {
				return nil, fileErr
			}
		}
		sources = append(sources, dto.PropertySource{
			Name:   candidate,
			Source: props,
		})
	}

	return sources, nil
}
