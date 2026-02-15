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

package errors

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpError(t *testing.T) {
	tests := []struct {
		name           string
		msg            string
		code           int
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Not Found Error",
			msg:            "Resource not found",
			code:           http.StatusNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"code":404,"error":"Resource not found"}`,
		},
		{
			name:           "Internal Server Error",
			msg:            "Internal error occurred",
			code:           http.StatusInternalServerError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"code":500,"error":"Internal error occurred"}`,
		},
		{
			name:           "Bad Request",
			msg:            "Invalid input",
			code:           http.StatusBadRequest,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":400,"error":"Invalid input"}`,
		},
		{
			name:           "Unauthorized",
			msg:            "Authentication required",
			code:           http.StatusUnauthorized,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"code":401,"error":"Authentication required"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a response recorder
			w := httptest.NewRecorder()

			// Call HttpError
			HttpError(w, tt.msg, tt.code)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("HttpError() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			// Check Content-Type header
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("HttpError() Content-Type = %s, want application/json", contentType)
			}

			// Check response body
			body := w.Body.String()
			// Trim newline that json.Encoder adds
			if len(body) > 0 && body[len(body)-1] == '\n' {
				body = body[:len(body)-1]
			}

			if body != tt.expectedBody {
				t.Errorf("HttpError() body = %s, want %s", body, tt.expectedBody)
			}
		})
	}
}
