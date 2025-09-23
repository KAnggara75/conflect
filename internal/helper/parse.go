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
 * @author KAnggara75 on Mon 22/09/25 09.01
 * @project conflect helper
 * https://github.com/PakaiWA/PakaiWA/tree/main/internal/helper
 */

package helper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseFile(data []byte, ext string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	switch ext {
	case ".yaml", ".yml":
		var m map[string]interface{}
		if err := yaml.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("yaml unmarshal: %w", err)
		}
		flattenMap("", m, out)
	case ".json":
		var m map[string]interface{}
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("json unmarshal: %w", err)
		}
		flattenMap("", m, out)
	case ".properties":
		m := parseProperties(data)
		// properties are usually flat keys already (like "a.b.c"), copy as-is
		for k, v := range m {
			out[k] = v
		}
	default:
		return nil, fmt.Errorf("unsupported ext: %s", ext)
	}

	return out, nil
}

// flattenMap flattens nested maps into dot.notation keys
func flattenMap(prefix string, cur interface{}, out map[string]interface{}) {
	switch t := cur.(type) {
	case map[string]interface{}:
		for k, v := range t {
			var key string
			if prefix == "" {
				key = k
			} else {
				key = prefix + "." + k
			}
			flattenMap(key, v, out)
		}
	case map[interface{}]interface{}: // in case yaml produced interface{} keys
		for kk, vv := range t {
			k := fmt.Sprintf("%v", kk)
			var key string
			if prefix == "" {
				key = k
			} else {
				key = prefix + "." + k
			}
			flattenMap(key, vv, out)
		}
	case []interface{}:
		// keep slices as-is (consumer can interpret), placed at prefix key
		out[prefix] = t
	default:
		out[prefix] = t
	}
}

// parseProperties parses key=value properties, tries to coerce to int/float/bool when sensible
func parseProperties(data []byte) map[string]interface{} {
	res := make(map[string]interface{})
	s := bufio.NewScanner(bytes.NewReader(data))
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}
		// support key:value and key=value
		var kv []string
		if strings.Contains(line, "=") {
			kv = strings.SplitN(line, "=", 2)
		} else if strings.Contains(line, ":") {
			kv = strings.SplitN(line, ":", 2)
		} else {
			continue
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		res[k] = parsePrimitive(v)
	}
	return res
}

func parsePrimitive(s string) interface{} {
	// try int
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	// try float
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	// try bool
	if b, err := strconv.ParseBool(strings.ToLower(s)); err == nil {
		return b
	}
	// fallback string
	return s
}
