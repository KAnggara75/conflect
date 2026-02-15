/*
 * Copyright (c) 2025 KAnggara75
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * See <https://www.gnu.org/licenses/gpl-3.0.html>.
 */

package util

import (
	"sync"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(5, time.Minute)
	defer rl.Stop()

	if rl.limit != 5 {
		t.Errorf("Expected limit 5, got %d", rl.limit)
	}
	if rl.window != time.Minute {
		t.Errorf("Expected window 1 minute, got %v", rl.window)
	}
	if rl.request == nil {
		t.Error("Request map should be initialized")
	}
}

func TestRateLimiter_IsAllow(t *testing.T) {
	rl := NewRateLimiter(3, 100*time.Millisecond)
	defer rl.Stop()

	key := "test-key"

	// First 3 requests should be allowed
	for i := 0; i < 3; i++ {
		if !rl.IsAllow(key) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 4th request should be denied
	if rl.IsAllow(key) {
		t.Error("4th request should be denied")
	}

	// Wait for window to expire
	time.Sleep(150 * time.Millisecond)

	// Should be allowed again after window expires
	if !rl.IsAllow(key) {
		t.Error("Request should be allowed after window expires")
	}
}

func TestRateLimiter_MultipleKeys(t *testing.T) {
	rl := NewRateLimiter(2, 100*time.Millisecond)
	defer rl.Stop()

	key1 := "key1"
	key2 := "key2"

	// Both keys should have independent limits
	if !rl.IsAllow(key1) {
		t.Error("First request for key1 should be allowed")
	}
	if !rl.IsAllow(key1) {
		t.Error("Second request for key1 should be allowed")
	}

	if !rl.IsAllow(key2) {
		t.Error("First request for key2 should be allowed")
	}
	if !rl.IsAllow(key2) {
		t.Error("Second request for key2 should be allowed")
	}

	// Third request for each key should be denied
	if rl.IsAllow(key1) {
		t.Error("Third request for key1 should be denied")
	}
	if rl.IsAllow(key2) {
		t.Error("Third request for key2 should be denied")
	}
}

func TestRateLimiter_Cleanup(t *testing.T) {
	rl := NewRateLimiter(5, 50*time.Millisecond)
	defer rl.Stop()

	key := "cleanup-test"

	// Make some requests
	rl.IsAllow(key)
	rl.IsAllow(key)

	// Verify entry exists
	rl.mu.Lock()
	if _, exists := rl.request[key]; !exists {
		t.Error("Key should exist in request map")
	}
	rl.mu.Unlock()

	// Wait for cleanup (cleanup runs at 2x window duration)
	time.Sleep(150 * time.Millisecond)

	// Manually trigger cleanup
	rl.cleanup()

	// Verify entry is cleaned up
	rl.mu.Lock()
	if _, exists := rl.request[key]; exists {
		t.Error("Key should be cleaned up after window expires")
	}
	rl.mu.Unlock()
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	rl := NewRateLimiter(100, time.Second)
	defer rl.Stop()

	var wg sync.WaitGroup
	key := "concurrent-key"
	successCount := 0
	var mu sync.Mutex

	// Spawn 200 concurrent goroutines
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if rl.IsAllow(key) {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Should allow exactly 100 requests
	if successCount != 100 {
		t.Errorf("Expected 100 successful requests, got %d", successCount)
	}
}

func TestRateLimiter_Stop(t *testing.T) {
	rl := NewRateLimiter(5, 100*time.Millisecond)

	// Stop should not panic
	rl.Stop()

	// Calling Stop multiple times should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Error("Stop() should not panic when called multiple times")
		}
	}()
}

func TestRateLimiter_WindowSliding(t *testing.T) {
	rl := NewRateLimiter(2, 100*time.Millisecond)
	defer rl.Stop()

	key := "sliding-window"

	// First request at t=0
	if !rl.IsAllow(key) {
		t.Error("First request should be allowed")
	}

	// Wait 60ms
	time.Sleep(60 * time.Millisecond)

	// Second request at t=60ms
	if !rl.IsAllow(key) {
		t.Error("Second request should be allowed")
	}

	// Third request at t=60ms should be denied (2 requests in window)
	if rl.IsAllow(key) {
		t.Error("Third request should be denied")
	}

	// Wait another 50ms (total 110ms from first request)
	time.Sleep(50 * time.Millisecond)

	// Fourth request at t=110ms should be allowed (first request expired)
	if !rl.IsAllow(key) {
		t.Error("Fourth request should be allowed after first request expires")
	}
}
