# Release Process

This document details the complete release process for the Req project.

## Release Stages

### Alpha (α)
- Very early development
- Core functionality incomplete
- For internal testing or early adopters
- Format: `v0.1.0-alpha.1`, `v0.1.0-alpha.2`

### Beta (β) 
- Core functionality mostly complete
- Feature-complete but may have bugs
- Ready for broader testing
- Format: `v0.1.0-beta.1`, `v0.1.0-beta.2`

### Release Candidate (RC)
- Nearly ready for release  
- All features complete, major bugs fixed
- Final testing phase
- Format: `v0.1.0-rc.1`, `v0.1.0-rc.2`

### Stable Release
- Production ready
- All major features work
- Thoroughly tested
- Format: `v0.1.0`, `v0.2.0`

## Current Project Stage

**Status**: Alpha stage (`v0.1.0-alpha.x`)

**Why Alpha**: Core HTTP execution features are still in development. The project has working TUI components and backend infrastructure, but is not feature-complete.

**Move to Beta when**: Core HTTP execution works, request/response viewer is functional, all major features are implemented.

## Release Workflow

### 1. Development Phase
- Work on feature branches
- Create pull requests to `main`
- Ensure CI/CD passes (build, test, format)

### 2. Pre-Release Checklist
Before creating a release, ensure:
- [ ] All intended features/fixes are merged to `main`
- [ ] CI/CD pipeline passes on `main`
- [ ] Manual testing completed
- [ ] Documentation updated if needed

### 3. Creating a Release

**Step 1**: Ensure you're on the latest `main`
```bash
git checkout main
git pull origin main
```

**Step 2**: Verify everything works
```bash
go build .
go test ./...
```

**Step 3**: Create and push the tag
```bash
# For next alpha release
git tag v0.1.0-alpha.2
git push origin v0.1.0-alpha.2
```

**Step 4**: GitHub Actions automatically:
- Builds binaries for Linux, macOS, Windows (amd64/arm64)
- Creates GitHub release with release notes
- Uploads all platform binaries

### 4. Version Numbering

**Alpha/Beta/RC increments:**
- `v0.1.0-alpha.1` → `v0.1.0-alpha.2` (bug fixes, small features)
- `v0.1.0-alpha.5` → `v0.1.0-beta.1` (feature complete)
- `v0.1.0-beta.3` → `v0.1.0-rc.1` (release candidate)
- `v0.1.0-rc.2` → `v0.1.0` (stable release)

**Minor/Major version bumps:**
- `v0.1.0` → `v0.2.0` (new features, backward compatible)
- `v0.9.0` → `v1.0.0` (major release, breaking changes)

## Automated Release Process

### Multi-Platform Builds
The release workflow automatically builds:
- `req-v0.1.0-alpha.2-linux-amd64`
- `req-v0.1.0-alpha.2-linux-arm64` 
- `req-v0.1.0-alpha.2-darwin-amd64`
- `req-v0.1.0-alpha.2-darwin-arm64`
- `req-v0.1.0-alpha.2-windows-amd64.exe`

### Release Notes
Generated from template in `.github/workflows/release.yml` including:
- Project description
- Installation instructions (`go install` and binary downloads)
- Breaking change warnings (database cleanup instructions)
- Usage information
- Development stage disclaimer
- Link to full changelog

### Prerelease Detection
Versions containing hyphens (alpha, beta, rc) are automatically marked as prereleases on GitHub.

## Manual Release Steps (if needed)

If automation fails, manual release process:

1. **Build binaries manually:**
```bash
# Linux amd64
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=v0.1.0-alpha.2" -o req-v0.1.0-alpha.2-linux-amd64

# macOS arm64  
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=v0.1.0-alpha.2" -o req-v0.1.0-alpha.2-darwin-arm64

# Windows amd64
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=v0.1.0-alpha.2" -o req-v0.1.0-alpha.2-windows-amd64.exe
```

2. **Create GitHub release manually** via web interface

## Release Schedule

No fixed schedule. Releases are created when:
- Significant features are complete
- Important bug fixes are ready
- Milestone requirements are met

Release `v0.1.0-alpha.x` versions as features are completed and tested.

## Rollback Process

If a release has critical issues:

1. **Do not delete the tag** (breaks `go install`)
2. **Create a new release** with fixes: `v0.1.0-alpha.3`
3. **Mark problematic release** as prerelease if not already
4. **Update documentation** to point to newer version

## Testing Releases

After creating a release:

```bash
# Test go install works
go install github.com/maniac-en/req@v0.1.0-alpha.2

# Test binary works
req --version  # Should show: req v0.1.0-alpha.2
```

## Communication

- Release announcements in project README
- Update docs/index.html with latest version
- Consider blog post for major releases