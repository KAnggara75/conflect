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
 * @author KAnggara75 on Mon 13/10/25 08.32
 * @project conflect middleware
 * https://github.com/PakaiWA/PakaiWA/tree/main/internal/delivery/http/middleware
 */

package middleware

import "net/http"

// Middleware adalah fungsi yang membungkus handler berikutnya.
type Middleware func(http.Handler) http.Handler

// Chain menyusun banyak middleware menjadi satu lapisan
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
