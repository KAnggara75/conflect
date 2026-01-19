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
 * @project conflect config
 * https://github.com/KAnggara75/conflect/tree/main/internal/config
 */

package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/KAnggara75/conflect/internal/helper"
)

type Config struct {
	Port          string
	RepoPath      string
	RepoURL       string
	DefaultBranch string
	Limit         int
	Token         string
}

func Load() *Config {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	defaultRepo := filepath.Join(cwd, "/etc/conflect/repo")

	return &Config{
		Limit:         getEnvInt("RATE_LIMIT", 10), // default 10 requests
		Port:          getEnv("APP_PORT", "8080"),
		RepoPath:      getEnv("REPO_PATH", defaultRepo),
		RepoURL:       buildRepoURL(),
		DefaultBranch: getEnv("DEFAULT_BRANCH", "main"),
		Token:         readValue("APP_AUTH_SECRET", "APP_AUTH_SECRET_FILE", ""),
	}
}

func buildRepoURL() string {
	token := readValue("GIT_AUTH_TOKEN", "GIT_AUTH_TOKEN_FILE", "")
	repoURL := readValue("REPO_URL", "REPO_URL_FILE", "")
	if repoURL == "" {
		return ""
	}
	return helper.NormalizeRepoURL(repoURL, token)
}

func readValue(envKey, fileKey, defaultValue string) string {
	if val := os.Getenv(envKey); val != "" {
		return strings.TrimSpace(val)
	}
	if filePath := os.Getenv(fileKey); filePath != "" {
		data, err := os.ReadFile(filePath)
		if err == nil {
			return strings.TrimSpace(string(data))
		}
	}
	return defaultValue
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return i
}
