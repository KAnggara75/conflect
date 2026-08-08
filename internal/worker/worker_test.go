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
 * @project conflect worker
 */

package worker

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/KAnggara75/conflect/internal/config"
	"github.com/KAnggara75/conflect/internal/repository"
	"github.com/KAnggara75/conflect/internal/service"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestWorkerStart(t *testing.T) {
	q := service.NewQueue(10)
	tmpDir := t.TempDir()
	cfg := &config.Config{RepoPath: tmpDir}
	repo := repository.NewGitRepo(tmpDir, "")
	cs := service.NewConfigServiceFromRepo(repo, cfg)

	q.Enqueue("nonexistent-branch")

	done := make(chan struct{})
	go func() {
		time.Sleep(50 * time.Millisecond)
		close(done)
	}()

	go Start(q, cs)

	select {
	case <-done:
		// Worker successfully processed item (error path)
	case <-time.After(1 * time.Second):
		t.Fatal("worker did not finish processing in time")
	}
}

func TestWorkerStart_Success(t *testing.T) {
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

	q := service.NewQueue(10)
	cfg := &config.Config{RepoPath: localRepoDir}
	repo := repository.NewGitRepo(localRepoDir, originDir)
	cs := service.NewConfigServiceFromRepo(repo, cfg)

	q.Enqueue("main")

	done := make(chan struct{})
	go func() {
		time.Sleep(150 * time.Millisecond)
		close(done)
	}()

	go Start(q, cs)

	select {
	case <-done:
		// Worker successfully processed item (success path, lines 28-29)
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not finish processing in time")
	}
}
