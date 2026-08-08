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
 * @project conflect service
 */

package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KAnggara75/conflect/internal/config"
	"github.com/KAnggara75/conflect/internal/repository"
)

func TestIsSafePathComponent(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"valid-name", true},
		{"main", true},
		{"app_123", true},
		{"", false},
		{"../parent", false},
		{"path/to", false},
		{"path\\to", false},
		{"..", false},
	}

	for _, tt := range tests {
		got := isSafePathComponent(tt.input)
		if got != tt.want {
			t.Errorf("isSafePathComponent(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestConfigService_LoadConfig_PathValidation(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		RepoPath:      tmpDir,
		DefaultBranch: "main",
	}
	repo := repository.NewGitRepo(tmpDir, "")
	cs := NewConfigServiceFromRepo(repo, cfg)

	t.Run("Unsafe appName", func(t *testing.T) {
		resp := cs.LoadConfig("../unsafe", "prod", "main")
		if len(resp.PropertySources) != 0 {
			t.Errorf("expected empty property sources for unsafe appName")
		}
	})

	t.Run("Unsafe env", func(t *testing.T) {
		resp := cs.LoadConfig("myapp", "../prod", "main")
		if len(resp.PropertySources) != 0 {
			t.Errorf("expected empty property sources for unsafe env")
		}
	})

	t.Run("Unsafe label", func(t *testing.T) {
		resp := cs.LoadConfig("myapp", "prod", "../label")
		if len(resp.PropertySources) != 0 {
			t.Errorf("expected empty property sources for unsafe label")
		}
	})
}

func TestConfigService_LoadConfig_SuccessAndCandidates(t *testing.T) {
	tmpDir := t.TempDir()
	envDir := filepath.Join(tmpDir, "main", "prod")
	if err := os.MkdirAll(envDir, 0755); err != nil {
		t.Fatalf("failed to create env dir: %v", err)
	}

	configFile := filepath.Join(envDir, "myapp-prod.yaml")
	content := []byte("server:\n  port: 8080\n")
	if err := os.WriteFile(configFile, content, 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Application-env config
	appEnvFile := filepath.Join(envDir, "application-prod.yml")
	if err := os.WriteFile(appEnvFile, []byte("env: prod"), 0644); err != nil {
		t.Fatalf("failed to write app env file: %v", err)
	}

	// Global config
	globalFile := filepath.Join(envDir, "application.yaml")
	globalContent := []byte("app:\n  name: test\n")
	if err := os.WriteFile(globalFile, globalContent, 0644); err != nil {
		t.Fatalf("failed to write global config file: %v", err)
	}

	// Ignored dir inside envDir
	_ = os.MkdirAll(filepath.Join(envDir, "ignored_subdir"), 0755)
	// Ignored file extension
	_ = os.WriteFile(filepath.Join(envDir, "notes.txt"), []byte("text"), 0644)

	cfg := &config.Config{
		RepoPath:      tmpDir,
		DefaultBranch: "main",
	}
	repo := repository.NewGitRepo(tmpDir, "")
	cs := NewConfigServiceFromRepo(repo, cfg)

	resp := cs.LoadConfig("myapp", "prod", "")

	if resp.Name != "myapp" {
		t.Errorf("expected Name 'myapp', got %q", resp.Name)
	}
	if resp.Label != "main" {
		t.Errorf("expected Label 'main', got %q", resp.Label)
	}
	if len(resp.PropertySources) != 3 {
		t.Fatalf("expected 3 property sources, got %d", len(resp.PropertySources))
	}
}

func TestConfigService_LoadConfig_ReadDirFail(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		RepoPath:      tmpDir,
		DefaultBranch: "main",
	}
	repo := repository.NewGitRepo(tmpDir, "")
	cs := NewConfigServiceFromRepo(repo, cfg)

	// Directory main/prod does not exist
	resp := cs.LoadConfig("myapp", "prod", "main")
	if len(resp.PropertySources) != 0 {
		t.Errorf("expected 0 property sources when directory missing")
	}
}

func TestConfigService_ListBranchesAndSHA(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(tmpDir, "main"), 0755)
	_ = os.MkdirAll(filepath.Join(tmpDir, "develop"), 0755)

	cfg := &config.Config{
		RepoPath:      tmpDir,
		DefaultBranch: "main",
	}
	repo := repository.NewGitRepo(tmpDir, "")
	cs := NewConfigServiceFromRepo(repo, cfg)

	branches, err := cs.ListBranches()
	if err != nil {
		t.Fatalf("unexpected error listing branches: %v", err)
	}

	if len(branches) != 2 {
		t.Errorf("expected 2 branches, got %d", len(branches))
	}

	_, err = cs.GetBranchSHA("nonexistent")
	if err == nil {
		t.Errorf("expected error for nonexistent branch SHA")
	}

	// UpdateRepo call for non-existent branch will fail
	err = cs.UpdateRepo("nonexistent")
	if err == nil {
		t.Errorf("expected error updating nonexistent branch repo")
	}
}
