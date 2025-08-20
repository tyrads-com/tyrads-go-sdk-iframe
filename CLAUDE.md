# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the TyrAds Go SDK for iframe integration - a lightweight Go wrapper that enables easy embedding of TyrAds web offerwalls. The SDK handles authentication, token generation, and iframe URL creation for TyrAds integration.

## Development Commands

Use the Makefile for streamlined development workflow:

```bash
# Display available commands
make help

# Run all CI checks (format, vet, lint, test)
make ci

# Individual commands
make build          # Build the module
make test           # Run tests with verbose output
make test-coverage  # Run tests with coverage report
make lint           # Run golangci-lint
make fmt            # Format code
make vet            # Run go vet
make deps           # Download dependencies
make tidy           # Tidy module dependencies
make clean          # Clean build artifacts

# Install development tools
make install-lint   # Install golangci-lint if not present
```

Standard Go commands are also available:
```bash
go build ./...      # Build the module
go test ./...       # Run tests
go fmt ./...        # Format code
go vet ./...        # Vet code for issues
golangci-lint run   # Run linter (after installation)
```

## Architecture

The codebase follows a clean, modular architecture:

### Core Components

- **`tyrads_sdk.go`** - Main SDK entry point containing `TyrAdsSdk` struct with two primary methods:
  - `Authenticate()` - Authenticates users and returns auth tokens
  - `IframeUrl()` - Generates iframe URLs for web integration
  
- **`config/`** - Configuration management with default API endpoints:
  - Base iframe URL: `https://sdk.tyrads.com`
  - API base URL: `https://api.tyrads.com`
  - API version: `v3.0`

- **`client/`** - HTTP client wrapper that handles:
  - API authentication via headers (`X-API-Key`, `X-API-Secret`)
  - JSON request/response marshalling
  - Error handling and status code validation

- **`contract/`** - Data models and validation:
  - `AuthenticationRequest` - User authentication payload with comprehensive validation
  - `AuthenticationSign` - Authentication response containing token and user info

- **`enum/`** - Constants for environment variables (`TYRADS_API_KEY`, `TYRADS_API_SECRET`)

### Key Patterns

1. **Functional Options Pattern** - Used in `NewConfig()` and `NewAuthenticationRequest()` for flexible initialization
2. **Interface Flexibility** - `IframeUrl()` accepts either string tokens or `AuthenticationSign` structs
3. **Environment Variable Fallbacks** - API credentials default to environment variables if not provided
4. **Comprehensive Validation** - Authentication requests validate email format, phone numbers, gender values, etc.

## Authentication Flow

1. Create `AuthenticationRequest` with required fields (publisherUserID, age, gender)
2. Call `sdk.Authenticate(request)` to get `AuthenticationSign` with token
3. Use token with `sdk.IframeUrl()` to generate iframe URL for embedding

## API Integration Details

- All API requests include SDK version and platform headers
- Language parameter is appended as query parameter (defaults to "en")
- Error responses are parsed and wrapped with descriptive messages
- HTTP client automatically handles JSON marshalling/unmarshalling

## GitHub Actions & CI/CD

The repository includes automated workflows for continuous integration and deployment:

### CI Pipeline (All Branches)
- **Build & Test**: Runs on Go 1.21 and 1.22 with race detection
- **Code Coverage**: Uploads coverage reports to Codecov
- **Linting**: Uses golangci-lint with comprehensive rule set
- **Security**: Gosec security scanner with SARIF reporting
- **Code Quality**: CodeQL analysis for security vulnerabilities

### Release Pipeline (Main Branch Only)
- **Automatic Versioning**: Uses conventional commits for semantic versioning
- **GitHub Releases**: Auto-creates releases with changelog
- **Coverage Reporting**: Enhanced coverage reporting on main branch
- **Dependency Updates**: Dependabot for Go modules and GitHub Actions

### Workflow Files
- `.github/workflows/ci.yml` - Main CI pipeline
- `.github/workflows/release.yml` - Release automation
- `.github/workflows/codeql.yml` - Security analysis
- `.github/dependabot.yml` - Dependency management

All workflows require `CODECOV_TOKEN` secret for coverage reporting.