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
 * @author KAnggara75 on Mon 13/10/25 08.37
 * @project conflect util
 * https://github.com/PakaiWA/PakaiWA/tree/main/internal/util
 */
package util

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu      sync.Mutex
	request map[string][]time.Time
	limit   int
	window  time.Duration
	stopCh  chan struct{}
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		request: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
		stopCh:  make(chan struct{}),
	}
	// Start background cleanup goroutine
	go rl.cleanupLoop()
	return rl
}

// cleanupLoop periodically removes stale entries to prevent memory leaks
func (r *RateLimiter) cleanupLoop() {
	// Cleanup interval is 2x the window duration
	ticker := time.NewTicker(r.window * 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.cleanup()
		case <-r.stopCh:
			return
		}
	}
}

// cleanup removes entries that have no recent requests
func (r *RateLimiter) cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.window)

	for key, times := range r.request {
		// Remove entries with no requests in the current window
		if len(times) == 0 {
			delete(r.request, key)
			continue
		}

		// Check if the most recent request is outside the window
		lastRequest := times[len(times)-1]
		if lastRequest.Before(windowStart) {
			delete(r.request, key)
		}
	}
}

// Stop stops the background cleanup goroutine
func (r *RateLimiter) Stop() {
	close(r.stopCh)
}

func (r *RateLimiter) IsAllow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.window)
	reqs := r.request[key]

	// Hapus request di luar jendela waktu
	validReqs := make([]time.Time, 0, len(reqs))
	for _, t := range reqs {
		if t.After(windowStart) {
			validReqs = append(validReqs, t)
		}
	}

	if len(validReqs) >= r.limit {
		r.request[key] = validReqs
		return false
	}

	r.request[key] = append(validReqs, now)
	return true
}
