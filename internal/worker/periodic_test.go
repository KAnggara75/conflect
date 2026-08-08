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
)

func TestStartPeriodicPull_Disabled(t *testing.T) {
	q := service.NewQueue(10)
	cfg := &config.Config{PullInterval: 0}

	done := make(chan struct{})
	go func() {
		StartPeriodicPull(cfg, q, nil)
		close(done)
	}()

	select {
	case <-done:
		// Function returned immediately because PullInterval <= 0
	case <-time.After(1 * time.Second):
		t.Fatal("StartPeriodicPull should return immediately when PullInterval <= 0")
	}
}

func TestStartPeriodicPull_NilConfig(t *testing.T) {
	q := service.NewQueue(10)

	done := make(chan struct{})
	go func() {
		StartPeriodicPull(nil, q, nil)
		close(done)
	}()

	select {
	case <-done:
		// Function returned immediately because cfg is nil
	case <-time.After(1 * time.Second):
		t.Fatal("StartPeriodicPull should return immediately when cfg is nil")
	}
}

func TestStartPeriodicPull_Active(t *testing.T) {
	q := service.NewQueue(10)
	tmpDir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(tmpDir, "main"), 0755)

	cfg := &config.Config{
		RepoPath:     tmpDir,
		PullInterval: 1, // 1 second interval
	}
	repo := repository.NewGitRepo(tmpDir, "")
	cs := service.NewConfigServiceFromRepo(repo, cfg)

	go StartPeriodicPull(cfg, q, cs)

	// Dequeue channel should receive branch within 2 seconds
	select {
	case branch := <-q.Dequeue():
		if branch != "main" {
			t.Errorf("expected branch 'main', got '%s'", branch)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("expected branch to be enqueued by periodic pull ticker")
	}
}

func TestStartPeriodicPull_ListBranchesError(t *testing.T) {
	q := service.NewQueue(10)
	nonExistentDir := filepath.Join(t.TempDir(), "nonexistent")

	cfg := &config.Config{
		RepoPath:     nonExistentDir,
		PullInterval: 1,
	}
	repo := repository.NewGitRepo(nonExistentDir, "")
	cs := service.NewConfigServiceFromRepo(repo, cfg)

	go StartPeriodicPull(cfg, q, cs)

	// Ticker should run, fail to list branches (lines 41-42), and not panic or deadlock
	time.Sleep(1200 * time.Millisecond)
}

func TestStartPeriodicPull_QueueFull(t *testing.T) {
	q := service.NewQueue(0) // Queue with capacity 0 (always full)
	tmpDir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(tmpDir, "main"), 0755)

	cfg := &config.Config{
		RepoPath:     tmpDir,
		PullInterval: 1,
	}
	repo := repository.NewGitRepo(tmpDir, "")
	cs := service.NewConfigServiceFromRepo(repo, cfg)

	go StartPeriodicPull(cfg, q, cs)

	// Ticker will trigger, list branches, attempt to enqueue, hit queue full branch (lines 49-50)
	time.Sleep(1200 * time.Millisecond)
}
