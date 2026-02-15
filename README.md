# Conflect

A distributed configuration server for managing application configurations across multiple environments and branches, inspired by Spring Cloud Config.

## Features

- 🔄 **Git-based Configuration**: Store configurations in Git repositories with branch support
- 🌍 **Multi-Environment**: Support for multiple environments (dev, staging, production, etc.)
- 📁 **Multiple Formats**: Support for YAML, JSON, and Properties files
- 🔐 **Authentication**: Token-based authentication and webhook signature verification
- ⚡ **Rate Limiting**: Built-in rate limiting to prevent abuse
- 🔔 **Webhook Support**: Automatic configuration updates via Git webhooks
- 📊 **Metrics**: Prometheus metrics for monitoring
- 🚀 **High Performance**: Efficient caching and concurrent request handling

## Test Coverage

![Coverage](https://img.shields.io/badge/coverage-32.0%25-yellow)

### Coverage by Package

| Package | Coverage | Status |
|---------|----------|--------|
| `internal/config` | 96.3% | ✅ Excellent |
| `internal/errors` | 100.0% | ✅ Perfect |
| `internal/helper` | 86.4% | ✅ Good |
| `internal/util` | 94.3% | ✅ Excellent |
| `internal/service` | 8.3% | ⚠️ Needs Improvement |
| `internal/repository` | 0.0% | ❌ No Tests |
| `internal/delivery/http` | 0.0% | ❌ No Tests |
| `internal/worker` | 0.0% | ❌ No Tests |

### Detailed Coverage Report

```
github.com/KAnggara75/conflect/internal/config/config.go:36:         Load            80.0%
github.com/KAnggara75/conflect/internal/config/config.go:53:         buildRepoURL    100.0%
github.com/KAnggara75/conflect/internal/config/config.go:62:         readValue       100.0%
github.com/KAnggara75/conflect/internal/config/config.go:75:         getEnv          100.0%
github.com/KAnggara75/conflect/internal/config/config.go:82:         getEnvInt       100.0%
github.com/KAnggara75/conflect/internal/errors/file.go:24:           ShouldSkipFile  100.0%
github.com/KAnggara75/conflect/internal/errors/http.go:23:           HttpError       100.0%
github.com/KAnggara75/conflect/internal/helper/parse.go:29:          ParseFile       100.0%
github.com/KAnggara75/conflect/internal/helper/parse.go:58:          flattenMap      56.2%
github.com/KAnggara75/conflect/internal/helper/parse.go:89:          parseProperties 93.8%
github.com/KAnggara75/conflect/internal/helper/parse.go:113:         parsePrimitive  100.0%
github.com/KAnggara75/conflect/internal/helper/url.go:24:            NormalizeRepoURL 100.0%
github.com/KAnggara75/conflect/internal/service/queue.go:23:         NewQueue        100.0%
github.com/KAnggara75/conflect/internal/service/queue.go:28:         Enqueue         100.0%
github.com/KAnggara75/conflect/internal/service/queue.go:38:         Dequeue         100.0%
github.com/KAnggara75/conflect/internal/util/ratelimiter.go:30:      NewRateLimiter  100.0%
github.com/KAnggara75/conflect/internal/util/ratelimiter.go:43:      cleanupLoop     100.0%
github.com/KAnggara75/conflect/internal/util/ratelimiter.go:59:      cleanup         81.8%
github.com/KAnggara75/conflect/internal/util/ratelimiter.go:82:      Stop            100.0%
github.com/KAnggara75/conflect/internal/util/ratelimiter.go:86:      IsAllow         100.0%
```

**Total Coverage: 32.0%**

## Installation

```bash
go get github.com/KAnggara75/conflect
```

## Configuration

Conflect can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | HTTP server port | `8080` |
| `RATE_LIMIT` | Rate limit (requests per minute) | `10` |
| `REPO_PATH` | Local path for git repository | `./etc/conflect/repo` |
| `REPO_URL` | Git repository URL | - |
| `DEFAULT_BRANCH` | Default git branch | `main` |
| `APP_AUTH_SECRET` | Authentication token | - |
| `GIT_AUTH_TOKEN` | Git authentication token | - |

### File-based Secrets

For sensitive values, you can use file-based configuration:

- `APP_AUTH_SECRET_FILE`: Path to file containing auth secret
- `GIT_AUTH_TOKEN_FILE`: Path to file containing git token
- `REPO_URL_FILE`: Path to file containing repository URL

## Usage

### Starting the Server

```bash
# Set required environment variables
export REPO_URL="https://github.com/your-org/config-repo"
export GIT_AUTH_TOKEN="your-github-token"
export APP_AUTH_SECRET="your-secret-token"

# Run the server
go run cmd/conflect/conflect.go
```

### API Endpoints

#### Health Check
```bash
GET /health
```

#### Get Configuration
```bash
GET /{application}/{environment}?label={branch}

# Example
GET /myapp/production?label=main
```

Response:
```json
{
  "name": "myapp",
  "profiles": ["production"],
  "label": "main",
  "version": "abc123...",
  "propertySources": [
    {
      "name": "myapp-production.yaml",
      "source": {
        "server.port": 8080,
        "database.host": "localhost"
      }
    }
  ]
}
```

#### Webhook (for automatic updates)
```bash
POST /webhook
X-Hub-Signature-256: sha256=...

# Payload
{
  "ref": "refs/heads/main"
}
```

### Configuration File Priority

Conflect loads configuration files in the following order (highest to lowest priority):

1. `{application}-{environment}.{ext}` (e.g., `myapp-production.yaml`)
2. `application-{environment}.{ext}` (e.g., `application-production.yaml`)
3. `application.{ext}` (e.g., `application.yaml`)

Supported file extensions: `.yaml`, `.yml`, `.json`, `.properties`

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out

# View coverage summary
go tool cover -func=coverage.out
```

## Development

### Project Structure

```
conflect/
├── cmd/
│   └── conflect/           # Main application entry point
├── internal/
│   ├── config/             # Configuration loading (96.3% coverage)
│   ├── delivery/
│   │   └── http/           # HTTP handlers and middleware
│   ├── errors/             # Error handling utilities (100% coverage)
│   ├── helper/             # Helper functions (86.4% coverage)
│   ├── repository/         # Git repository operations
│   ├── service/            # Business logic
│   ├── util/               # Utilities (94.3% coverage)
│   └── worker/             # Background workers
└── README.md
```

### Building

```bash
# Build binary
go build -o conflect cmd/conflect/conflect.go

# Build for production
go build -ldflags="-s -w" -o conflect cmd/conflect/conflect.go
```

## Monitoring

Conflect exposes Prometheus metrics at `/metrics`:

- HTTP request duration
- Request count by endpoint
- Rate limit hits
- Configuration load times

## License

Copyright (c) 2025 KAnggara75

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

See <https://www.gnu.org/licenses/gpl-3.0.html>.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Testing Guidelines

- Write tests for all new features
- Maintain or improve code coverage
- Run `go test ./...` before submitting PR
- Follow existing code style and conventions

## Roadmap

- [ ] Increase test coverage to >80%
- [ ] Add integration tests for HTTP handlers
- [ ] Add tests for Git repository operations
- [ ] Implement configuration encryption
- [ ] Add support for multiple Git repositories
- [ ] WebSocket support for real-time config updates
- [ ] Admin UI for configuration management
