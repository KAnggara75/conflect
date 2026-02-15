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
	"reflect"
	"testing"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		ext      string
		expected map[string]interface{}
		wantErr  bool
	}{
		{
			name: "YAML simple",
			data: []byte(`
server:
  port: 8080
  host: localhost
`),
			ext: ".yaml",
			expected: map[string]interface{}{
				"server.port": int(8080),
				"server.host": "localhost",
			},
			wantErr: false,
		},
		{
			name: "JSON simple",
			data: []byte(`{
  "server": {
    "port": 8080,
    "host": "localhost"
  }
}`),
			ext: ".json",
			expected: map[string]interface{}{
				"server.port": float64(8080),
				"server.host": "localhost",
			},
			wantErr: false,
		},
		{
			name: "Properties file",
			data: []byte(`
# Comment
server.port=8080
server.host=localhost
database.enabled=true
database.timeout=30.5
`),
			ext: ".properties",
			expected: map[string]interface{}{
				"server.port":      int64(8080),
				"server.host":      "localhost",
				"database.enabled": true,
				"database.timeout": float64(30.5),
			},
			wantErr: false,
		},
		{
			name: "Properties with colon separator",
			data: []byte(`
server.port:8080
server.host:localhost
`),
			ext: ".properties",
			expected: map[string]interface{}{
				"server.port": int64(8080),
				"server.host": "localhost",
			},
			wantErr: false,
		},
		{
			name:     "Unsupported extension",
			data:     []byte("some data"),
			ext:      ".txt",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Invalid YAML",
			data:     []byte("invalid: yaml: content: ["),
			ext:      ".yaml",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Invalid JSON",
			data:     []byte("{invalid json}"),
			ext:      ".json",
			expected: nil,
			wantErr:  true,
		},
		{
			name: "YAML with arrays",
			data: []byte(`
items:
  - item1
  - item2
  - item3
`),
			ext: ".yml",
			expected: map[string]interface{}{
				"items": []interface{}{"item1", "item2", "item3"},
			},
			wantErr: false,
		},
		{
			name: "Nested YAML",
			data: []byte(`
database:
  connection:
    host: localhost
    port: 5432
  pool:
    min: 5
    max: 20
`),
			ext: ".yaml",
			expected: map[string]interface{}{
				"database.connection.host": "localhost",
				"database.connection.port": int(5432),
				"database.pool.min":        int(5),
				"database.pool.max":        int(20),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseFile(tt.data, tt.ext)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParsePrimitive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "Integer",
			input:    "123",
			expected: int64(123),
		},
		{
			name:     "Negative integer",
			input:    "-456",
			expected: int64(-456),
		},
		{
			name:     "Float",
			input:    "123.45",
			expected: float64(123.45),
		},
		{
			name:     "Boolean true",
			input:    "true",
			expected: true,
		},
		{
			name:     "Boolean false",
			input:    "false",
			expected: false,
		},
		{
			name:     "String",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePrimitive(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parsePrimitive(%q) = %v (type %T), want %v (type %T)",
					tt.input, result, result, tt.expected, tt.expected)
			}
		})
	}
}

func TestParseProperties(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected map[string]interface{}
	}{
		{
			name: "Basic properties",
			data: []byte(`
key1=value1
key2=value2
`),
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "Properties with comments",
			data: []byte(`
# This is a comment
key1=value1
! Another comment
key2=value2
`),
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "Properties with types",
			data: []byte(`
port=8080
host=localhost
enabled=true
timeout=30.5
`),
			expected: map[string]interface{}{
				"port":    int64(8080),
				"host":    "localhost",
				"enabled": true,
				"timeout": float64(30.5),
			},
		},
		{
			name: "Empty lines",
			data: []byte(`
key1=value1

key2=value2

`),
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseProperties(tt.data)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseProperties() = %v, want %v", result, tt.expected)
			}
		})
	}
}
