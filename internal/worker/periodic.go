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
 * https://github.com/KAnggara75/conflect/tree/main/internal/worker
 */

package worker

import (
	"log"
	"time"

	"github.com/KAnggara75/conflect/internal/config"
	"github.com/KAnggara75/conflect/internal/service"
)

// StartPeriodicPull starts a ticker if cfg.PullInterval > 0 and periodically enqueues all branches.
func StartPeriodicPull(cfg *config.Config, q *service.Queue, s *service.ConfigService) {
	if cfg == nil || cfg.PullInterval <= 0 {
		log.Println("ℹ️  Periodic pull disabled (PULL_INTERVAL is not set or <= 0)")
		return
	}

	log.Printf("⏱️  Starting periodic pull worker with interval %d seconds...", cfg.PullInterval)

	ticker := time.NewTicker(time.Duration(cfg.PullInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		branches, err := s.ListBranches()
		if err != nil {
			log.Printf("❌ Failed to list branches for periodic pull: %v", err)
			continue
		}

		log.Printf("🔄 Periodic pull triggered for %d branch(es)...", len(branches))
		for _, branch := range branches {
			if q.Enqueue(branch) {
				log.Printf("📥 Enqueued branch %q via periodic pull", branch)
			} else {
				log.Printf("⚠️  Queue full, skipped periodic pull for branch %q", branch)
			}
		}
	}
}
