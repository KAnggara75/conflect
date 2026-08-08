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
	"time"

	"github.com/KAnggara75/conflect/internal/config"
	"github.com/KAnggara75/conflect/internal/repository"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

func TestNewConfigService(t *testing.T) {
	originDir := t.TempDir()
	originGit, err := git.PlainInit(originDir, false)
	if err != nil {
		t.Fatalf("failed to init origin repo: %v", err)
	}

	worktree, _ := originGit.Worktree()
	_ = os.WriteFile(filepath.Join(originDir, "file.txt"), []byte("v1"), 0644)
	_, _ = worktree.Add("file.txt")
	hash1, _ := worktree.Commit("commit 1", &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@example.com", When: time.Now()},
	})
	_ = originGit.Storer.SetReference(plumbing.NewHashReference(plumbing.NewBranchReferenceName("main"), hash1))

	localRepoDir := t.TempDir()
	cfg := &config.Config{
		RepoPath: localRepoDir,
		RepoURL:  originDir,
	}

	cs := NewConfigService(cfg)
	if cs == nil || cs.repo == nil {
		t.Fatalf("NewConfigService returned nil or invalid struct")
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

	// Create a valid git repository inside tmpDir/main so GetCommitHashFromBranch succeeds (lines 101-104)
	gitRepo, err := git.PlainInit(filepath.Join(tmpDir, "main"), false)
	if err != nil {
		t.Fatalf("failed to init git repo in main: %v", err)
	}
	wt, _ := gitRepo.Worktree()
	_ = os.WriteFile(filepath.Join(tmpDir, "main", "README.md"), []byte("test"), 0644)
	_, _ = wt.Add("README.md")
	commitHash, _ := wt.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@example.com", When: time.Now()},
	})

	configFile := filepath.Join(envDir, "myapp-prod.yaml")
	content := []byte("server:\n  port: 8080\n")
	if err := os.WriteFile(configFile, content, 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	appEnvFile := filepath.Join(envDir, "application-prod.yml")
	if err := os.WriteFile(appEnvFile, []byte("env: prod"), 0644); err != nil {
		t.Fatalf("failed to write app env file: %v", err)
	}

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
	if resp.Version != commitHash.String() {
		t.Errorf("expected Version %q, got %q", commitHash.String(), resp.Version)
	}
	if len(resp.PropertySources) != 3 {
		t.Fatalf("expected 3 property sources, got %d", len(resp.PropertySources))
	}
}

func TestConfigService_LoadConfig_ReadFileAndParseErrors(t *testing.T) {
	t.Run("Permission denied on ReadFile", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("skipping permission test when running as root")
		}

		tmpDir := t.TempDir()
		envDir := filepath.Join(tmpDir, "main", "prod")
		_ = os.MkdirAll(envDir, 0755)

		configFile := filepath.Join(envDir, "myapp-prod.yaml")
		_ = os.WriteFile(configFile, []byte("test"), 0000)
		defer os.Chmod(configFile, 0644)

		cfg := &config.Config{RepoPath: tmpDir, DefaultBranch: "main"}
		repo := repository.NewGitRepo(tmpDir, "")
		cs := NewConfigServiceFromRepo(repo, cfg)

		resp := cs.LoadConfig("myapp", "prod", "main")
		if len(resp.PropertySources) != 0 {
			t.Errorf("expected empty property sources on read permission error")
		}
	})

	t.Run("Invalid file content that is skipped", func(t *testing.T) {
		tmpDir := t.TempDir()
		envDir := filepath.Join(tmpDir, "main", "prod")
		_ = os.MkdirAll(envDir, 0755)

		// Invalid YAML format -> ParseFile error -> ShouldSkipFile returns true
		configFile := filepath.Join(envDir, "myapp-prod.yaml")
		_ = os.WriteFile(configFile, []byte("invalid:\n  - item1\n item2"), 0644)

		cfg := &config.Config{RepoPath: tmpDir, DefaultBranch: "main"}
		repo := repository.NewGitRepo(tmpDir, "")
		cs := NewConfigServiceFromRepo(repo, cfg)

		resp := cs.LoadConfig("myapp", "prod", "main")
		if len(resp.PropertySources) != 0 {
			t.Errorf("expected empty property sources when file is skipped due to parse error")
		}
	})
}

func TestConfigService_LoadConfig_ReadDirFail(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		RepoPath:      tmpDir,
		DefaultBranch: "main",
	}
	repo := repository.NewGitRepo(tmpDir, "")
	cs := NewConfigServiceFromRepo(repo, cfg)

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

	err = cs.UpdateRepo("nonexistent")
	if err == nil {
		t.Errorf("expected error updating nonexistent branch repo")
	}
}
