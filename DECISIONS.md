# DECISIONS.md

Dokumen ini mencatat **Architectural Decision Records (ADR)** untuk proyek **Conflect**. Setiap keputusan arsitektur penting dicatat lengkap dengan konteks, keputusan yang diambil, serta implikasi (*consequences*) atau *trade-offs*.

---

## ADR 001: Penggunaan Standard Go Project Layout & Package Boundary `internal/`

* **Status**: Accepted
* **Tanggal**: 2025-09-22

### Context
Aplikasi Conflect dikembangkan sebagai server konfigurasi terdistribusi dalam bahasa Go. Diperlukan struktur direktori yang terorganisir, standar, dan mencegah dependensi liar antar package serta pengimporan kode internal oleh modul eksternal.

### Decision
Menerapkan [Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_id.md) dengan memisahkan kode menjadi `cmd/` untuk entrypoint aplikasi dan `internal/` untuk seluruh komponen logika aplikasi (config, delivery, service, repository, worker, util).

### Consequences / Trade-offs
* **Keuntungan**:
  * Struktur mudah dipahami oleh developer Go profesional.
  * Perlindungan tingkat compiler (*Go compiler enforcement*) mencegah package luar mengimpor isi `internal/`.
  * Pembagian layer yang tegas mempermudah penulisan unit test secara terisolasi.
* **Kerugian**:
  * Jumlah direktori dan file bertambah (*indirection*), memerlukan pembacaan awal struktur bagi kontributor baru.

---

## ADR 002: In-Memory Asynchronous Worker Queue untuk Git Webhook Synchronizations

* **Status**: Accepted
* **Tanggal**: 2025-09-22

### Context
Saat Git Provider (seperti GitHub/GitLab) mengirimkan Webhook trigger akibat adanya *push* ke branch, Conflect harus melakukan pembaruan lokal repositori (*git pull/fetch*). Jika proses Git ini dilakukan secara synchronous pada HTTP thread Webhook:
1. Latency Webhook request akan sangat tinggi (tergantung kecepatan koneksi Git).
2. Jika ada lonjakan push (*burst webhooks*), server bisa mengalami kehabisan resource (*thread starvation*).

### Decision
Menggunakan **Buffered Go Channel Queue** (`internal/service/queue.go`) di memori yang dipadukan dengan **Background Worker Goroutine** (`internal/worker/worker.go`). Endpoint `/webhook` hanya bertugas memvalidasi request dan memasukkan (*enqueue*) nama branch ke antrean memori, kemudian langsung mengembalikan HTTP Status `202 Accepted`.

### Consequences / Trade-offs
* **Keuntungan**:
  * Response time endpoint `/webhook` sangat cepat (< 5ms).
  * Melindungi server dari *overload* saat menerima lonjakan webhook bersamaan.
  * Sederhana tanpa perlu dependensi message broker eksternal seperti RabbitMQ atau Redis.
* **Kerugian**:
  * Antrean tersimpan di dalam memori (*in-memory*). Jika aplikasi *crash* atau *restart* saat ada item antrean yang belum diproses, pesan dalam antrean tersebut akan hilang (dapat dipulihkan pada pemanggilan HTTP config berikutnya).

---

## ADR 003: Penggunaan Go Standard Library `net/http` dengan Custom Middleware Chain

* **Status**: Accepted
* **Tanggal**: 2025-09-22

### Context
Aplikasi membutuhkan HTTP REST API untuk menyajikan konfigurasi, menerima webhook, dan mengekspos metrik Prometheus. Pemilihan web framework menentukan besar binary, performa, dan jumlah dependensi eksternal.

### Decision
Menggunakan **Standard Library Go (`net/http`)** yang dikombinasikan dengan pembungkus *Custom Middleware Chain* (`middleware.Chain`) tanpa menggunakan framework pihak ketiga seperti Fiber, Echo, atau Gin.

### Consequences / Trade-offs
* **Keuntungan**:
  * Ukuran binary executable sangat kecil dan kompilasi cepat.
  * Performa dan efisiensi alokasi memori maksimal.
  * Tanpa dependensi pihak ketiga untuk routing, mengurangi risiko kerentanan keamanan (*supply chain vulnerability*).
  * Kompatibilitas sempurna dengan exporter Prometheus resmi (`client_golang`).
* **Kerugian**:
  * Routing parameter URL (seperti `/{app}/{env}/{label}`) harus diparse secara manual menggunakan manipulasi string pada handler.

---

## ADR 004: In-Memory Sliding Window / Token Rate Limiter

* **Status**: Accepted
* **Tanggal**: 2025-09-22

### Context
Untuk mencegah penyalahgunaan (*abuse*) atau serangan *Denial of Service (DoS)* terhadap endpoint konfigurasi dan webhook, server memerlukan mekanisme *rate limiting* berdasarkan IP atau Token.

### Decision
Mengimplementasikan **In-Memory Rate Limiter** (`internal/util/ratelimiter.go`) menggunakan mekanisme pembersihan berkala (*cleanup ticker*) untuk membatasi jumlah request per menit.

### Consequences / Trade-offs
* **Keuntungan**:
  * Performa tinggi karena tidak membutuhkan pembacaan ke network database eksternal.
  * Mudah diuji menggunakan unit test.
* **Kerugian**:
  * Pembatasan rate limiter berlaku per instance aplikasi. Jika Conflect di-deploy secara multi-instance (cluster) di belakang Load Balancer, rate limit belum terdistribusi secara terpusat (dapat ditingkatkan dengan backend Redis di masa mendatang).
