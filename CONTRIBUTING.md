# Contributing to Req

## Release Process

To create a new release:

1. Ensure changes are merged to `main`
2. Create and push a Git tag:
   ```bash
   git tag v0.1.0-alpha.2
   git push origin v0.1.0-alpha.2
   ```
3. GitHub Actions automatically builds and creates the release

For detailed release information, see [RELEASE.md](RELEASE.md).

## GitHub Actions Workflows

### CI/CD Pipeline (`.github/workflows/ci-cd.yml`)
Triggers: Push to `main`, Pull requests to `main`

Jobs:
- build-test: Cross-platform testing (Ubuntu/macOS/Windows)
  - Builds project and runs tests with coverage
- format: Auto-formats Go code with `gofmt`
  - Commits formatting changes automatically
  - Skips if author is `github-actions[bot]` to prevent cascading

### Release Pipeline (`.github/workflows/release.yml`)
Triggers: Git tag push matching `v*` pattern

Jobs:
- release: Builds binaries for Linux/macOS/Windows (amd64/arm64)
- create-release: Creates GitHub release with binaries and release notes

## Development Workflow

### Building and Testing
```bash
# Build for development
go build -o ./tmp/req .

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Format code (or let CI handle it)
go fmt ./...
```

### Database Changes
```bash
# Create new migration
goose -dir db/migrations create migration_name sql

# Generate SQLC code after modifying queries
sqlc generate
```

### Installation Methods
After release, users can install via:
```bash
# Latest stable release
go install github.com/maniac-en/req@latest

# Specific version
go install github.com/maniac-en/req@v0.1.0-alpha.2
```

## Branch Strategy
- main: Primary development branch
- Feature branches: Create for new work, merge via PRs
- Releases: Use Git tags, not version branches

## Troubleshooting
- CI formatting failures: Auto-fixed by format job
- Release not triggering: Ensure Git tag starts with `v`
- GitHub Pages builds: Expected to be cancelled (we use Git tags for distribution)