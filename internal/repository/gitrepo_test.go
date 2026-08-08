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
 * @project conflect repository
 */

package repository

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestNewGitRepo(t *testing.T) {
	repo := NewGitRepo("/tmp/test", "https://github.com/test/repo.git")
	if repo.Path != "/tmp/test" || repo.URL != "https://github.com/test/repo.git" {
		t.Errorf("NewGitRepo initialized incorrectly")
	}
}

func TestGitRepo_InitAllBranches_EmptyURL(t *testing.T) {
	repo := NewGitRepo("/tmp/test", "")
	err := repo.InitAllBranches()
	if err == nil {
		t.Error("expected error when URL is empty")
	}
}

func TestGitRepo_ListLocalBranches(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(tmpDir, "main"), 0755)
	_ = os.MkdirAll(filepath.Join(tmpDir, "develop"), 0755)

	repo := NewGitRepo(tmpDir, "")
	branches, err := repo.ListLocalBranches()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(branches) != 2 {
		t.Errorf("expected 2 local branches, got %d", len(branches))
	}
}

func TestGitRepo_LocalGitOperations(t *testing.T) {
	tmpDir := t.TempDir()
	branchPath := filepath.Join(tmpDir, "main")

	// Initialize a local git repository in branchPath
	gitRepo, err := git.PlainInit(branchPath, false)
	if err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	worktree, err := gitRepo.Worktree()
	if err != nil {
		t.Fatalf("failed to get worktree: %v", err)
	}

	testFile := filepath.Join(branchPath, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test Repo"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err = worktree.Add("README.md")
	if err != nil {
		t.Fatalf("failed to add file to git: %v", err)
	}

	hash, err := worktree.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	repo := NewGitRepo(tmpDir, "")

	t.Run("GetCommitHashFromBranch", func(t *testing.T) {
		gotHash, err := repo.GetCommitHashFromBranch("main")
		if err != nil {
			t.Fatalf("unexpected error getting commit hash: %v", err)
		}
		if gotHash != hash.String() {
			t.Errorf("expected commit hash %s, got %s", hash.String(), gotHash)
		}
	})

	t.Run("EnsureBranch Directory Exists", func(t *testing.T) {
		path, err := repo.EnsureBranch("main")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if path != branchPath {
			t.Errorf("expected path %s, got %s", branchPath, path)
		}
	})
}

func TestGitRepo_Pull(t *testing.T) {
	originDir := t.TempDir()
	originGit, err := git.PlainInit(originDir, false)
	if err != nil {
		t.Fatalf("failed to init origin repo: %v", err)
	}

	worktree, err := originGit.Worktree()
	if err != nil {
		t.Fatalf("failed to get worktree: %v", err)
	}

	_ = os.WriteFile(filepath.Join(originDir, "file.txt"), []byte("v1"), 0644)
	_, _ = worktree.Add("file.txt")
	hash1, _ := worktree.Commit("commit 1", &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@example.com", When: time.Now()},
	})
	_ = originGit.Storer.SetReference(plumbing.NewHashReference(plumbing.NewBranchReferenceName("main"), hash1))

	localRepoDir := t.TempDir()
	mainDir := filepath.Join(localRepoDir, "main")
	_, err = git.PlainClone(mainDir, false, &git.CloneOptions{
		URL: originDir,
	})
	if err != nil {
		t.Fatalf("failed to clone: %v", err)
	}

	_ = os.WriteFile(filepath.Join(originDir, "file.txt"), []byte("v2"), 0644)
	_, _ = worktree.Add("file.txt")
	_, _ = worktree.Commit("commit 2", &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@example.com", When: time.Now()},
	})

	repo := NewGitRepo(localRepoDir, originDir)
	err = repo.Pull("main")
	if err != nil {
		t.Fatalf("unexpected error during pull: %v", err)
	}
}

func TestGitRepo_EnsureBranch_Clone(t *testing.T) {
	originDir := t.TempDir()
	originGit, err := git.PlainInit(originDir, false)
	if err != nil {
		t.Fatalf("failed to init origin repo: %v", err)
	}

	worktree, err := originGit.Worktree()
	if err != nil {
		t.Fatalf("failed to get worktree: %v", err)
	}

	_ = os.WriteFile(filepath.Join(originDir, "file.txt"), []byte("v1"), 0644)
	_, _ = worktree.Add("file.txt")
	hash1, _ := worktree.Commit("commit 1", &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@example.com", When: time.Now()},
	})
	_ = originGit.Storer.SetReference(plumbing.NewHashReference(plumbing.NewBranchReferenceName("main"), hash1))

	localRepoDir := t.TempDir()
	repo := NewGitRepo(localRepoDir, originDir)

	targetPath, err := repo.EnsureBranch("main")
	if err != nil {
		t.Fatalf("EnsureBranch failed: %v", err)
	}

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Errorf("expected branch targetPath %s to exist", targetPath)
	}
}
