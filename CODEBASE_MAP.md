# CODEBASE_MAP.md

Peta direktori dan berkas proyek **Conflect** disusun berdasarkan pedoman resmi [Golang Standard Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_id.md).

```
conflect/
├── .air.toml                    # Konfigurasi Live Reloading untuk pengembangan lokal (Air)
├── .dockerignore                # Daftar file/folder yang dikecualikan dari Docker build
├── .gitignore                   # Daftar file/folder yang diabaikan oleh Git
├── .pre-commit-config.yaml      # Konfigurasi hook gitleaks & golangci-lint
├── LICENSE                      # Lisensi proyek (GNU General Public License v3.0)
├── Makefile                     # Build automation & test runner commands
├── README.md                    # Dokumentasi utama proyek & petunjuk penggunaan
├── SECURITY.md                  # Kebijakan pelaporan kerentanan keamanan
├── TESTING.md                   # Panduan eksekusi & standar pengujian
├── TEST_SUMMARY.md              # Laporan ringkasan cakupan pengujian (test coverage)
├── build/                       # Berkas dan naskah pendukung build & packaging
│   ├── go-version-check.sh      # Script verifikasi kompatibilitas versi Go
│   └── package/                 # Packaging manifests & Dockerfile
├── cmd/                         # Entrypoint aplikasi utama (Main Applications)
│   └── conflect/                # Executable binary package untuk Conflect
│       └── conflect.go          # Main function: bootstrap server, DI, graceful shutdown
├── go.mod                       # Definisi modul Go & dependensi proyek
├── go.sum                       # Checksum verifikasi dependensi Go
├── internal/                    # Kode privat aplikasi (tidak ter-expose ke luar modul)
│   ├── config/                  # Modul pembacaan konfigurasi environment aplikasi
│   │   ├── config.go            # Struct Config & pemuat environment / secret files
│   │   └── config_test.go       # Unit test untuk pemuat konfigurasi
│   ├── delivery/                # Presentation Layer (Delivery Mechanism)
│   │   └── http/                # Server HTTP, routing, dan handler
│   │       ├── dto/             # Data Transfer Objects
│   │       │   └── response.go  # Struct response JSON standar untuk API
│   │       ├── middleware/      # Middleware HTTP modular
│   │       │   ├── auth.go      # Middleware autentikasi token Bearer
│   │       │   ├── chain.go     # Utility pembuat rantai (chaining) middleware
│   │       │   ├── logger.go    # Middleware pencatatan log HTTP request
│   │       │   ├── ratelimit.go # Middleware pembatas laju request (Rate Limiting)
│   │       │   └── signature.go # Middleware verifikasi signature HMAC Webhook Git
│   │       └── handler.go       # Router ServeMux & HTTP handler endpoints
│   ├── errors/                  # Modul kustom error & HTTP error mapper
│   │   ├── file.go              # Utility pembacaan & identifikasi file skip
│   │   ├── file_test.go         # Unit test file error
│   │   ├── http.go              # Struct HttpError & penanganan HTTP status code
│   │   └── http_test.go         # Unit test HTTP error mapper
│   ├── helper/                  # Modul pembantu umum (Parsing & URL)
│   │   ├── parse.go             # Parser untuk YAML, JSON, dan Properties files
│   │   ├── parse_test.go        # Unit test parser konfigurasi
│   │   ├── url.go               # Utility normalisasi URL Git
│   │   └── url_test.go          # Unit test normalisasi URL
│   ├── repository/              # Data Access Layer / Storage Abstraction
│   │   └── gitrepo.go           # Abstraksi repositori Git menggunakan go-git
│   ├── service/                 # Business Logic & Core Domain
│   │   ├── config_service.go    # Resolusi hirarki konfigurasi & penggabungan properti
│   │   ├── queue.go             # Implementasi antrean buffered channel untuk Webhook
│   │   └── queue_test.go        # Unit test antrean pembaruan
│   ├── util/                    # Utility umum terpakai ulang
│   │   ├── ratelimiter.go       # Implementasi rate limiter (Sliding Window / Token)
│   │   └── ratelimiter_test.go  # Unit test rate limiter
│   └── worker/                  # Background Jobs & Asynchronous Processing
│       └── worker.go            # Worker loop konsumen antrean sinkronisasi Git
└── loadTest/                    # Script pengujian beban (Load Testing / Benchmark)
```

---

## Ringkasan Fungsi Direktori Utama

| Direktori | Deskripsi & Fungsi |
| :--- | :--- |
| `cmd/` | Tempat penyimpanan main entrypoint aplikasi. Kode di sini harus efisien dan hanya bertugas menghubungkan komponen. |
| `internal/` | Kode inti aplikasi yang dilindungi Go compiler agar tidak bisa di-import oleh repositori/modul luar. |
| `internal/config` | Memuat variabel lingkungan (`.env`) dan secret berkas (`*_FILE`). |
| `internal/delivery` | Menangani layer komunikasi eksternal (HTTP REST API, WebSocket, Webhook). |
| `internal/service` | Menjadi tempat beradanya aturan bisnis (*business rules*), penggabungan konfigurasi, dan antrean internal. |
| `internal/repository` | Mengurus interaksi dengan repositori Git lokal dan remote via `go-git`. |
| `internal/worker` | Menjalankan *background tasks* untuk pembaruan Git tanpa menghambat proses HTTP utama. |
| `internal/util` & `internal/helper` | Berisi fungsi utilitas murni seperti parsing file, manipulasi string, dan rate limiting. |
| `build/` | Berisi skrip otomatisasi build dan spesifikasi packaging container. |
