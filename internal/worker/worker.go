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
 * @author KAnggara75 on Mon 22/09/25 07.41
 * @project conflect worker
 * https://github.com/KAnggara75/conflect/tree/main/internal/worker
 */

package worker

import (
	"log"

	"github.com/KAnggara75/conflect/internal/service"
)

func Start(q *service.Queue, s *service.ConfigService) {
	for range q.Dequeue() {
		if err := s.UpdateRepo(); err != nil {
			log.Printf("repo update failed: %v", err)
		} else {
			log.Println("repo updated successfully")
		}
	}
}
