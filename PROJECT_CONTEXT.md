# PROJECT_CONTEXT.md

## 1. Visi & Tujuan Proyek

**Conflect (Config Reflect)** adalah *distributed configuration server* berkinerja tinggi yang ditulis dalam bahasa Go, terinspirasi oleh Spring Cloud Config. 

Tujuan utama Conflect adalah menyediakan layanan terpusat untuk mengelola dan mendistribusikan konfigurasi aplikasi (*application configurations*) lintas berbagai lingkungan (*environments*) dan cabang Git (*branches*). Conflect membaca berkas konfigurasi dari repositori Git, mendukung berbagai format file (YAML, JSON, Properties), menyajikan konfigurasi melalui RESTful HTTP API, serta menangani pembaruan otomatis secara *asynchronous* melalui Webhook Git.

---

## 2. Tech Stack & Framework

* **Bahasa Pemrograman**: Go (`go 1.25+`)
* **HTTP Framework / Routing**: Go Standard Library `net/http` dengan *custom middleware chaining*
* **Git Integration**: `github.com/go-git/go-git/v5`
* **Parsing Utilities**: `gopkg.in/yaml.v3` (YAML), standard `encoding/json` (JSON), *custom parser* (Properties)
* **Observability & Metrics**: `github.com/prometheus/client_golang` (Prometheus Metrics at `/metrics`)
* **DevOps & Tooling**: `Makefile`, `.air.toml` (Live Reload), `.pre-commit-config.yaml` (Code Quality & Linting)

---

## 3. Konvensi Koding (Coding Conventions)

1. **Struktur Kode**: Mematuhi [Golang Standard Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_id.md) (`cmd/`, `internal/`). Semua logika bisnis internal tersimpan di direktori `internal/` agar tidak ter-expose keluar modul.
2. **Error Handling**: 
   * Selalu tangani error secara eksplisit (`if err != nil`).
   * Dilarang mengabaikan (*ignore*) error atau menggunakan `panic()` di luar fungsi `main` / `init`.
   * Gunakan tipe error terstruktur dari `internal/errors`.
3. **Concurrency & Thread Safety**: 
   * Gunakan `sync.RWMutex` untuk akses data bersama (*shared memory*).
   * Gunakan Go *channels* untuk komunikasi antar *goroutines* (misal: antrean Webhook worker).
4. **Naming Conventions**:
   * Nama paket: *lowercase*, tunggal, tanpa *underscore* (misal: `config`, `service`, `repository`).
   * Nama interface / struct: `CamelCase` atau `PascalCase`.
   * Nama berkas: `lowercase_with_underscore.go`.
5. **Testing Standard**:
   * Setiap *package* wajib dilengkapi unit test (`*_test.go`).
   * Target coverage minimal project: `>= 80%`.

---

## 4. Core Principles

* **Git as Single Source of Truth**: Semua konfigurasi aplikasi disimpan di repositori Git terpusat.
* **Stateless & High Performance**: Cache lokal repositori dikelola secara efisien untuk menyajikan konfigurasi dengan latency rendah tanpa tergantung database eksternal.
* **Zero Security Compromise**:
  * Otentikasi berbasis *Bearer Token* untuk pengambilan konfigurasi HTTP API.
  * Verifikasi tanda tangan HMAC-SHA256 (`X-Hub-Signature-256`) untuk Webhook Git.
  * Dukungan *file-based secrets* (`*_FILE`) untuk lingkungan containerized/Kubernetes.
* **Non-blocking Event Loop**: Pembaruan repositori melalui Webhook dilakukan secara *asynchronous* melalui antrean background worker (`Queue`).

---

## 5. Hard Constraints (Aturan yang Dilarang)

* ❌ **DILARANG** mengekspos *internal error stack traces* atau path sistem lokal ke dalam *response HTTP*.
* ❌ **DILARANG** menyimpan *credentials*, *private keys*, atau *secrets* secara *hardcoded* di dalam kode program.
* ❌ **DILARANG** melakukan mutasi variabel global tanpa mekanis proteksi *concurrency* (`mutex` / `channels`).
* ❌ **DILARANG** memasukkan logika bisnis atau *git execution* langsung ke dalam HTTP Handler (Wajib melalui layer `service` dan `repository`).
* ❌ **DILARANG** membuat *blocking calls* di dalam HTTP request thread utama untuk proses *git pull/fetch* yang memakan waktu lama.
