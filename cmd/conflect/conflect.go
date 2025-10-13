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
 * @author KAnggara75 on Mon 22/09/25 07.29
 * @project conflect conflect
 * https://github.com/KAnggara75/conflect/tree/main/cmd/conflect
 */

package main

import (
	"log"
	"time"

	"github.com/KAnggara75/conflect/internal/config"
	"github.com/KAnggara75/conflect/internal/delivery/http"
	"github.com/KAnggara75/conflect/internal/service"
	"github.com/KAnggara75/conflect/internal/util"
	"github.com/KAnggara75/conflect/internal/worker"
)

func main() {
	cfg := config.Load()

	// queue and service layer
	queue := service.NewQueue(100)
	configService := service.NewConfigService(cfg, queue)

	// start worker
	go worker.Start(queue, configService)

	limiter := util.NewRateLimiter(20, time.Minute) // contoh 20 req per menit per IP

	// start HTTP server
	server := http.NewServer(cfg, queue, configService, limiter)
	log.Fatal(server.Start())
}
