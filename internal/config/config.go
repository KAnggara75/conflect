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
 * @author KAnggara75 on Mon 22/09/25 07.39
 * @project conflect config
 * https://github.com/KAnggara75/conflect/tree/main/internal/config
 */

package config

import (
	"os"
	"path/filepath"

	"github.com/KAnggara75/conflect/internal/helper"
)

type Config struct {
	Port          string
	RepoPath      string
	RepoURL       string
	WebhookSecret string
}

func Load() *Config {
	cwd, _ := os.Getwd()
	defaultRepo := filepath.Join(cwd, "tmp/repo")
	repoPathFile := getEnv("REPO_PATH_FILE", defaultRepo)
	webhookSecretFile := getEnv("WEBHOOK_SECRET_FILE", "")

	repoTokenFile := getEnv("TOKEN_FILE", "")
	repoToken := getEnv("TOKEN", repoTokenFile)

	repoUrlFile := getEnv("REPO_URL_FILE", "")

	repoUrl := getEnv("REPO_URL", repoUrlFile)

	url := helper.NormalizeRepoURL(repoUrl, repoToken)
	return &Config{
		Port:          getEnv("APP_PORT", "8080"),
		RepoPath:      getEnv("REPO_PATH", repoPathFile),
		RepoURL:       url,
		WebhookSecret: getEnv("WEBHOOK_SECRET", webhookSecretFile),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
