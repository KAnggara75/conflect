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

package helper

import (
	"testing"
)

func TestNormalizeRepoURL(t *testing.T) {
	tests := []struct {
		name     string
		rawURL   string
		token    string
		expected string
	}{
		{
			name:     "URL with https prefix",
			rawURL:   "https://github.com/user/repo",
			token:    "mytoken",
			expected: "https://mytoken@github.com/user/repo.git",
		},
		{
			name:     "URL with http prefix",
			rawURL:   "http://github.com/user/repo",
			token:    "mytoken",
			expected: "https://mytoken@github.com/user/repo.git",
		},
		{
			name:     "URL without protocol",
			rawURL:   "github.com/user/repo",
			token:    "mytoken",
			expected: "https://mytoken@github.com/user/repo.git",
		},
		{
			name:     "URL already with .git suffix",
			rawURL:   "https://github.com/user/repo.git",
			token:    "mytoken",
			expected: "https://mytoken@github.com/user/repo.git",
		},
		{
			name:     "URL with special characters in token",
			rawURL:   "https://github.com/user/repo",
			token:    "token@with+special/chars",
			expected: "https://token%40with%2Bspecial%2Fchars@github.com/user/repo.git",
		},
		{
			name:     "Empty token",
			rawURL:   "https://github.com/user/repo",
			token:    "",
			expected: "https://@github.com/user/repo.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeRepoURL(tt.rawURL, tt.token)
			if result != tt.expected {
				t.Errorf("NormalizeRepoURL(%q, %q) = %q, want %q",
					tt.rawURL, tt.token, result, tt.expected)
			}
		})
	}
}
