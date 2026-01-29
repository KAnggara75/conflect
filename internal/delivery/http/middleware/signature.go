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
 * @author KAnggara75 on Thu 30/10/25 23.33
 * @project conflect middleware
 * https://github.com/PakaiWA/PakaiWA/tree/main/internal/delivery/http/middleware
 */
package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func VerifySignature(cfg AuthConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const prefix = "sha256="
			signatureHeader := r.Header.Get("X-Hub-Signature-256")

			// Baca body dan simpan salinannya agar bisa dipakai ulang
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "failed to read body", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(body))

			// Validasi header
			if !strings.HasPrefix(signatureHeader, prefix) {
				writeUnauthorized(w, "Failed to Verify signatureHeader")
				return
			}

			// Hitung ulang HMAC
			mac := hmac.New(sha256.New, []byte(cfg.Token))
			mac.Write(body)
			expected := hex.EncodeToString(mac.Sum(nil))
			actual := signatureHeader[len(prefix):]

			// Timing-safe compare
			if !hmac.Equal([]byte(expected), []byte(actual)) {
				writeUnauthorized(w, "Failed to Verify Signature")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": "Unauthorized",
		"msg":   msg,
	})
}
