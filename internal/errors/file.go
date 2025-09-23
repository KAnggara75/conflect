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
 * @author KAnggara75 on Tue 23/09/25 22.47
 * @project conflect errors
 * https://github.com/PakaiWA/PakaiWA/tree/main/internal/errors
 */

package errors

import (
	"errors"
	"fmt"
	"os"
)

func ShouldSkipFile(candidate string, err error) (skip bool, finalErr error) {
	if err == nil {
		return false, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("skip file %s: not found\n", candidate)
		return true, nil
	}

	return false, fmt.Errorf("failed to process file %s: %w", candidate, err)
}
