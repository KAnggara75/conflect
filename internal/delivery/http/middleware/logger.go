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
 * @author KAnggara75 on Mon 13/10/25 08.45
 * @project conflect middleware
 * https://github.com/PakaiWA/PakaiWA/tree/main/internal/delivery/http/middleware
 */

package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging middleware
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s in %v", r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
	})
}
