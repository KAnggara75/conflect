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
	"errors"
	"os"
	"testing"
)

func TestShouldSkipFile(t *testing.T) {
	tests := []struct {
		name         string
		candidate    string
		err          error
		expectSkip   bool
		expectErrNil bool
	}{
		{
			name:         "No error",
			candidate:    "test.yaml",
			err:          nil,
			expectSkip:   false,
			expectErrNil: true,
		},
		{
			name:         "File not found - should skip",
			candidate:    "missing.yaml",
			err:          os.ErrNotExist,
			expectSkip:   true,
			expectErrNil: true,
		},
		{
			name:         "Permission denied - should not skip",
			candidate:    "forbidden.yaml",
			err:          os.ErrPermission,
			expectSkip:   false,
			expectErrNil: false,
		},
		{
			name:         "Generic error - should not skip",
			candidate:    "error.yaml",
			err:          errors.New("some error"),
			expectSkip:   false,
			expectErrNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skip, err := ShouldSkipFile(tt.candidate, tt.err)

			if skip != tt.expectSkip {
				t.Errorf("ShouldSkipFile() skip = %v, want %v", skip, tt.expectSkip)
			}

			if tt.expectErrNil && err != nil {
				t.Errorf("ShouldSkipFile() error = %v, want nil", err)
			}

			if !tt.expectErrNil && err == nil {
				t.Error("ShouldSkipFile() error = nil, want non-nil error")
			}
		})
	}
}
