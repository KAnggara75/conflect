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
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		request: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
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
