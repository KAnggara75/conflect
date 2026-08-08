# TODO.md

Rencana kerja dan daftar tugas pengembangan **Conflect (Config Reflect)** yang disusun secara bertahap berdasarkan prioritas eksekusi.

---

## Phase 1: Core Setup, Architecture & Security (Foundation)

- [x] **Project Layout Setup**: Mengatur struktur direktori berdasarkan Golang Standard Project Layout (`cmd/`, `internal/`).
- [x] **Config Management**: Implementasi pemuat konfigurasi environment dan dukungan file rahasia (`internal/config`).
- [x] **Standard Error Handling**: Membuat utility penanganan error HTTP dan file secara terpusat (`internal/errors`).
- [x] **Multi-Format Parsers**: Implementasi parser untuk file YAML, JSON, dan Properties (`internal/helper/parse.go`).
- [x] **Git Repository Abstraction**: Implementasi wrapper `go-git` untuk proses clone, fetch, dan checkout (`internal/repository/gitrepo.go`).
- [x] **Asynchronous Webhook Queue**: Membuat antrean channel memori dan background worker untuk sinkronisasi Git (`internal/service/queue.go`, `internal/worker/worker.go`).
- [x] **HTTP Middleware**: Implementasi middleware Auth Token, Webhook HMAC Signature Verification, Rate Limiter, dan Logger (`internal/delivery/http/middleware`).
- [x] **Prometheus Telemetry**: Integrasi metrik HTTP total request dan kegagalan verifikasi signature.
- [ ] **HTTP Handler Unit Tests**: Menulis unit test komprehensif untuk `internal/delivery/http/handler.go` menggunakan `net/http/httptest`.
- [ ] **Git Repository Unit Tests**: Menulis unit test dan mock test untuk `internal/repository/gitrepo.go`.

---

## Phase 2: Feature Enhancement & Robustness (Features & Reliability)

- [ ] **Multi-Repository Support**: Menambahkan dukungan untuk membaca dari banyak repositori Git sekaligus (*multi-source configuration*).
- [ ] **Property Value Encryption ({cipher})**: Implementasi enkripsi dan dekripsi nilai rahasia dalam file konfigurasi (misal `{cipher}A1B2C3...`) menggunakan AES-256.
- [ ] **Real-time Push Notifications (WebSocket/SSE)**: Menambahkan endpoint WebSocket atau Server-Sent Events (SSE) agar aplikasi client mendapatkan notifikasi pembaruan konfigurasi secara instan tanpa perlu polling.
- [ ] **Configuration Caching Layer**: Membuat in-memory cache dengan TTL untuk hasil parsing konfigurasi guna meminimalkan I/O disk lokal.
- [ ] **Distributed Rate Limiter Option**: Menyediakan opsi backend Redis untuk rate limiter ketika Conflect di-deploy dalam mode multi-instance / cluster.
- [ ] **Health Check Deep Probe**: Meningkatkan endpoint `/health` agar memeriksa ketersediaan direktori Git lokal dan konektivitas remote Git.

---

## Phase 3: Testing, Optimization, UI & Deployment (Production Readiness)

- [ ] **Code Coverage Target (>80%)**: Meningkatkan total test coverage dari 32% saat ini hingga mencapai minimal 80% di seluruh package.
- [ ] **Integration & E2E Tests**: Membuat skenario pengujian End-to-End untuk menguji alur utuh dari Git Webhook push hingga pencapaian konfigurasi di client API.
- [ ] **Admin Web Dashboard UI**: Membangun antarmuka web sederhana untuk menginspeksi status repositori Git, branch terhubung, cache status, dan log aktivitas.
- [ ] **Containerization & Helm Chart**: Menyediakan Dockerfile multi-stage build yang sangat efisien dan Helm Chart resmi untuk Kubernetes deployment.
- [ ] **CI/CD Pipeline & Automated Release**: Membuat GitHub Actions workflow untuk linting (`golangci-lint`), pengujian otomatis, dan publikasi image Docker secara otomatis saat rilis versi baru.
