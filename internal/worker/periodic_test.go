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
