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
 * @author KAnggara75 on Mon 13/10/25 08.33
 * @project conflect middleware
 * https://github.com/PakaiWA/PakaiWA/tree/main/internal/delivery/http/middleware
 */

package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
)

type AuthConfig struct {
	Token string
}

func AuthMiddleware(cfg AuthConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				resp := map[string]string{
					"error": "Unauthorized",
				}
				_ = json.NewEncoder(w).Encode(resp)
				return
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			if token != cfg.Token {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				resp := map[string]string{
					"error": "Invalid token",
				}
				_ = json.NewEncoder(w).Encode(resp)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
