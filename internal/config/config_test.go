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

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback string
		envValue string
		expected string
	}{
		{
			name:     "Environment variable set",
			key:      "TEST_VAR",
			fallback: "default",
			envValue: "custom",
			expected: "custom",
		},
		{
			name:     "Environment variable not set",
			key:      "UNSET_VAR",
			fallback: "default",
			envValue: "",
			expected: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			defer os.Unsetenv(tt.key)

			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := getEnv(tt.key, tt.fallback)
			if result != tt.expected {
				t.Errorf("getEnv(%q, %q) = %q, want %q", tt.key, tt.fallback, result, tt.expected)
			}
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback int
		envValue string
		expected int
	}{
		{
			name:     "Valid integer",
			key:      "TEST_INT",
			fallback: 10,
			envValue: "42",
			expected: 42,
		},
		{
			name:     "Invalid integer",
			key:      "TEST_INT",
			fallback: 10,
			envValue: "not-a-number",
			expected: 10,
		},
		{
			name:     "Empty value",
			key:      "TEST_INT",
			fallback: 10,
			envValue: "",
			expected: 10,
		},
		{
			name:     "Negative integer",
			key:      "TEST_INT",
			fallback: 10,
			envValue: "-5",
			expected: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			defer os.Unsetenv(tt.key)

			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := getEnvInt(tt.key, tt.fallback)
			if result != tt.expected {
				t.Errorf("getEnvInt(%q, %d) = %d, want %d", tt.key, tt.fallback, result, tt.expected)
			}
		})
	}
}

func TestReadValue(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "secret.txt")
	testContent := "file-secret-value"
	if err := os.WriteFile(testFile, []byte(testContent+"  \n"), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name         string
		envKey       string
		fileKey      string
		defaultValue string
		envValue     string
		fileValue    string
		expected     string
	}{
		{
			name:         "Environment variable takes precedence",
			envKey:       "TEST_ENV",
			fileKey:      "TEST_FILE",
			defaultValue: "default",
			envValue:     "env-value",
			fileValue:    testFile,
			expected:     "env-value",
		},
		{
			name:         "File value when env not set",
			envKey:       "TEST_ENV",
			fileKey:      "TEST_FILE",
			defaultValue: "default",
			envValue:     "",
			fileValue:    testFile,
			expected:     "file-secret-value",
		},
		{
			name:         "Default value when neither set",
			envKey:       "TEST_ENV",
			fileKey:      "TEST_FILE",
			defaultValue: "default",
			envValue:     "",
			fileValue:    "",
			expected:     "default",
		},
		{
			name:         "Trimmed whitespace from env",
			envKey:       "TEST_ENV",
			fileKey:      "TEST_FILE",
			defaultValue: "default",
			envValue:     "  env-value  ",
			fileValue:    "",
			expected:     "env-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			defer os.Unsetenv(tt.envKey)
			defer os.Unsetenv(tt.fileKey)

			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
			}
			if tt.fileValue != "" {
				os.Setenv(tt.fileKey, tt.fileValue)
			}

			result := readValue(tt.envKey, tt.fileKey, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("readValue() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// Save original env vars
	originalPort := os.Getenv("APP_PORT")
	originalLimit := os.Getenv("RATE_LIMIT")
	originalBranch := os.Getenv("DEFAULT_BRANCH")

	// Clean up after test
	defer func() {
		if originalPort != "" {
			os.Setenv("APP_PORT", originalPort)
		} else {
			os.Unsetenv("APP_PORT")
		}
		if originalLimit != "" {
			os.Setenv("RATE_LIMIT", originalLimit)
		} else {
			os.Unsetenv("RATE_LIMIT")
		}
		if originalBranch != "" {
			os.Setenv("DEFAULT_BRANCH", originalBranch)
		} else {
			os.Unsetenv("DEFAULT_BRANCH")
		}
		os.Unsetenv("REPO_URL")
		os.Unsetenv("REPO_PATH")
	}()

	// Set test environment variables
	os.Setenv("APP_PORT", "9090")
	os.Setenv("RATE_LIMIT", "20")
	os.Setenv("DEFAULT_BRANCH", "develop")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("Load() Port = %s, want 9090", cfg.Port)
	}

	if cfg.Limit != 20 {
		t.Errorf("Load() Limit = %d, want 20", cfg.Limit)
	}

	if cfg.DefaultBranch != "develop" {
		t.Errorf("Load() DefaultBranch = %s, want develop", cfg.DefaultBranch)
	}

	if cfg.RepoPath == "" {
		t.Error("Load() RepoPath should not be empty")
	}
}

func TestLoadDefaults(t *testing.T) {
	// Clean environment
	os.Unsetenv("APP_PORT")
	os.Unsetenv("RATE_LIMIT")
	os.Unsetenv("DEFAULT_BRANCH")
	os.Unsetenv("REPO_URL")
	os.Unsetenv("REPO_PATH")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("Load() default Port = %s, want 8080", cfg.Port)
	}

	if cfg.Limit != 10 {
		t.Errorf("Load() default Limit = %d, want 10", cfg.Limit)
	}

	if cfg.DefaultBranch != "main" {
		t.Errorf("Load() default DefaultBranch = %s, want main", cfg.DefaultBranch)
	}
}
