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

package service

import (
	"testing"
)

func TestNewQueue(t *testing.T) {
	tests := []struct {
		name     string
		size     int
		wantSize int
	}{
		{
			name:     "Standard queue size",
			size:     100,
			wantSize: 100,
		},
		{
			name:     "Small queue",
			size:     10,
			wantSize: 10,
		},
		{
			name:     "Large queue",
			size:     1000,
			wantSize: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewQueue(tt.size)
			if q == nil {
				t.Fatal("NewQueue() returned nil")
			}
			if q.ch == nil {
				t.Error("Queue channel should be initialized")
			}
			if cap(q.ch) != tt.wantSize {
				t.Errorf("NewQueue(%d) channel capacity = %d, want %d", tt.size, cap(q.ch), tt.wantSize)
			}
		})
	}
}

func TestQueue_Enqueue(t *testing.T) {
	q := NewQueue(3)

	// Enqueue should succeed when queue is not full
	if !q.Enqueue("branch1") {
		t.Error("Enqueue should succeed on empty queue")
	}

	if !q.Enqueue("branch2") {
		t.Error("Enqueue should succeed when queue has space")
	}

	if !q.Enqueue("branch3") {
		t.Error("Enqueue should succeed when queue has space")
	}

	// Queue is now full, next enqueue should fail
	if q.Enqueue("branch4") {
		t.Error("Enqueue should fail when queue is full")
	}
}

func TestQueue_Dequeue(t *testing.T) {
	q := NewQueue(10)

	testBranch := "main"
	q.Enqueue(testBranch)

	// Dequeue returns a channel
	ch := q.Dequeue()
	if ch == nil {
		t.Fatal("Dequeue() should return a channel")
	}

	// Read from the channel
	select {
	case branch := <-ch:
		if branch != testBranch {
			t.Errorf("Dequeued branch = %s, want %s", branch, testBranch)
		}
	default:
		t.Error("Should be able to read from dequeue channel")
	}
}

func TestQueue_EnqueueDequeue(t *testing.T) {
	q := NewQueue(5)

	branches := []string{"main", "develop", "feature-1", "feature-2"}

	// Enqueue all branches
	for _, branch := range branches {
		if !q.Enqueue(branch) {
			t.Errorf("Failed to enqueue branch: %s", branch)
		}
	}

	// Dequeue and verify order (FIFO)
	ch := q.Dequeue()
	for i, expected := range branches {
		select {
		case branch := <-ch:
			if branch != expected {
				t.Errorf("Dequeue[%d] = %s, want %s", i, branch, expected)
			}
		default:
			t.Errorf("Failed to dequeue branch at index %d", i)
		}
	}
}
