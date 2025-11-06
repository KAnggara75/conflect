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
 * @author KAnggara75 on Tue 23/09/25 10.55
 * @project conflect helper
 * https://github.com/PakaiWA/PakaiWA/tree/main/internal/helper
 */

package helper

import (
	"fmt"
	"net/url"
	"strings"
)

func NormalizeRepoURL(rawURL string, token string) string {
	clean := strings.TrimPrefix(rawURL, "https://")
	clean = strings.TrimPrefix(clean, "http://")

	if !strings.HasSuffix(clean, ".git") {
		clean = clean + ".git"
	}

	return fmt.Sprintf("https://%s@%s", url.QueryEscape(token), clean)
}
