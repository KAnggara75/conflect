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
	"path/filepath"

	"github.com/KAnggara75/conflect/internal/config"
	"github.com/KAnggara75/conflect/internal/delivery/http/dto"
	"github.com/KAnggara75/conflect/internal/helper"
	"github.com/KAnggara75/conflect/internal/repository"
	"github.com/go-git/go-git/v5/plumbing"
)

type ConfigService struct {
	repo *repository.GitRepo
	cfg  *config.Config
}

func NewConfigService(cfg *config.Config, q *Queue) *ConfigService {
	repo := repository.NewGitRepo(cfg.RepoPath, cfg.RepoURL)
	err := repo.InitAllBranches()
	if err != nil {
		log.Fatalf("failed to clone repo: %v", err)
	}
	return &ConfigService{repo: repo, cfg: cfg}
}

func (c *ConfigService) UpdateRepo() error {
	log.Println("Pulling latest config...")
	return c.repo.Pull()
}

func (c *ConfigService) LoadConfig(appName, env, label string) *dto.ConfigResponse {

	response := &dto.ConfigResponse{
		Name:     appName,
		Profiles: []string{env},
	}

	var refName plumbing.ReferenceName
	if label == "" {
		def, err := c.repo.DefaultBranch()
		if err != nil {
			log.Println(err)
			return response
		}
		refName = plumbing.NewRemoteReferenceName("origin", def)
	} else {
		refName = plumbing.NewRemoteReferenceName("origin", label)
	}

	candidates := c.generateConfigCandidates(appName, env)

	data, err := c.findAndReadAllConfigs(refName.String(), candidates)
	if err != nil {
		log.Println(err)
		return response
	}
	response.PropertySources = data

	hash, err := c.repo.GetCommitHash(refName)
	if err == nil {
		response.Version = hash
	}

	return response
}

func (c *ConfigService) generateConfigCandidates(appName, env string) []string {
	prefixes := []string{appName, "application"}
	extensions := []string{".yaml", ".yml", ".json", ".properties"}

	var files []string
	for _, p := range prefixes {
		files = append(files, fmt.Sprintf("%s/%s-%s", env, p, env))
	}

	var candidates []string
	for _, f := range files {
		for _, ext := range extensions {
			candidates = append(candidates, f+ext)
		}
	}

	fmt.Println("candidates:", candidates)

	return candidates
}

func (c *ConfigService) findAndReadAllConfigs(refName string, candidates []string) ([]dto.PropertySource, error) {
	var sources []dto.PropertySource

	// resolve commit sekali
	hash, err := c.repo.Repo.ResolveRevision(plumbing.Revision(refName))
	if err != nil {
		return nil, err
	}
	commit, err := c.repo.Repo.CommitObject(*hash)
	if err != nil {
		return nil, err
	}

	for _, name := range candidates {
		relPath := filepath.ToSlash(name)
		log.Printf("Loading config from repo path: %s", relPath)

		data, err := c.repo.GetFile(commit, relPath)
		if err != nil {
			log.Printf("ERROR %v", err)
			continue
		}

		ext := filepath.Ext(name)
		src, err := helper.ParseFile(data, ext)
		if err != nil {
			return nil, err
		}

		sources = append(sources, dto.PropertySource{
			Name:   relPath,
			Source: src,
		})
	}

	if len(sources) == 0 {
		return nil, fmt.Errorf("no config file found in %s (ref: %s)", c.repo.Path, refName)
	}

	return sources, nil
}
