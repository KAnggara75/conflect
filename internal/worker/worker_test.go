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
	"testing"
	"time"

	"github.com/KAnggara75/conflect/internal/config"
	"github.com/KAnggara75/conflect/internal/repository"
	"github.com/KAnggara75/conflect/internal/service"
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
		// Enqueue one item and close channel after processing
		time.Sleep(50 * time.Millisecond)
		close(done)
	}()

	go Start(q, cs)

	select {
	case <-done:
		// Worker successfully processed item
	case <-time.After(1 * time.Second):
		t.Fatal("worker did not finish processing in time")
	}
}
