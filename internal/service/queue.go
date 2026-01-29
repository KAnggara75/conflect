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
 * @author KAnggara75 on Mon 22/09/25 07.40
 * @project conflect service
 * https://github.com/KAnggara75/conflect/tree/main/internal/service
 */
package service

import "log"

type Queue struct {
	ch chan string
}

func NewQueue(size int) *Queue {
	return &Queue{ch: make(chan string, size)}
}

// Enqueue adds a branch to the queue. Returns true if successful, false if queue is full.
func (q *Queue) Enqueue(branch string) bool {
	select {
	case q.ch <- branch:
		return true
	default:
		log.Printf("⚠️  Queue full, dropping branch update: %s", branch)
		return false
	}
}

func (q *Queue) Dequeue() <-chan string {
	return q.ch
}
