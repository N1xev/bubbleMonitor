# CI/CD & Linting Guide

This document explains the CI/CD pipeline and golangci-lint configuration added to this project.

## Table of Contents

- [golangci-lint](#golangci-lint)
- [GitHub Actions CI/CD](#github-actions-cicd)
- [Adding New Stuff](#adding-new-stuff)

---

## golangci-lint

### What is it?

golangci-lint is a fast Go linter that runs multiple linters in parallel. It catches bugs, style issues, and suspicious code patterns before they reach production.

### Why use it?

- Catches common mistakes early (unchecked errors, unused variables, etc.)
- Enforces consistent code style across the project
- Reduces code review nitpicks
- Runs 10+ linters simultaneously (much faster than running them one by one)

### How to run

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Run with specific config
golangci-lint run -c .golangci.yml
```

### Configuration (.golangci.yml)

The project uses these linters:

| Linter | Purpose |
|--------|---------|
| errcheck | Catches unchecked errors |
| govet | Reports suspicious constructs |
| staticcheck | Static analysis checks |
| gocritic | Bug and style diagnostics |
| gofmt | Checks code formatting |
| goimports | Validates import statements |
| revive | Modern replacement for golint |
| gocyclo | Checks function complexity |
| dupl | Detects duplicate code |
| misspell | Catches spelling mistakes |

### Adding new linters

To enable a new linter, add it to the `linters.enable` section:

```yaml
linters:
  enable:
    - your-new-linter  # Add here
```

### Customizing linter settings

Each linter has settings in `linters-settings`:

```yaml
linters-settings:
  gocyclo:
    min-complexity: 15  # Increase for more lenient complexity check
```

---

## GitHub Actions CI/CD

### What is it?

GitHub Actions automates your workflow. The CI pipeline runs on every push and PR.

### Why use it?

- Catches issues before merging
- Tests on multiple platforms (Linux, macOS, Windows)
- Builds release binaries automatically
- No manual testing needed

### How it works

The workflow (`.github/workflows/ci.yml`) has 3 jobs:

1. **Lint** - Runs golangci-lint
2. **Test** - Runs tests with coverage
3. **Build** - Compiles for Linux, macOS, Windows

### Adding new jobs

Add a new job to the workflow:

```yaml
  security:
    name: Security Scan
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      
      - name: Run security scan
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...
```

### Adding new build targets

Extend the matrix in the build job:

```yaml
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [linux, darwin, windows, freebsd]
        goarch: [amd64, arm64, armv6]
```

### Running workflows locally

```bash
# Install act (GitHub Actions local runner)
brew install act

# Run workflow locally
act -l  # List workflows
act     # Run default workflow
```

---

## Adding New Stuff

### Adding a new linter

1. Find the linter on [golangci-lint linters page](https://golangci-lint.run/usage/linters/)
2. Add to `.golangci.yml`:

```yaml
linters:
  enable:
    - linter-name

linters-settings:
  linter-name:
    setting: value
```

3. Test locally: `golangci-lint run`

### Adding a new CI job

1. Edit `.github/workflows/ci.yml`
2. Add a new job:

```yaml
  job-name:
    name: Job Display Name
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - run: your-command
```

3. Test with `act` locally

### Adding a new build platform

Edit the matrix in the build job:

```yaml
strategy:
  matrix:
    goos: [linux, darwin, windows, freebsd]
    goarch: [amd64, arm64]
```

### Adding a new dependency scan

```yaml
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      
      - name: Run govulncheck
        run: go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...
      
      - name: Run dependency review
        uses: actions/dependency-review-action@v4
```

---

## Common Issues

### Linter fails on existing code

Fix issues manually or temporarily exclude files:

```yaml
issues:
  exclude-rules:
    - path: old-file.go
      linters:
        - gocritic
```

### CI takes too long

Enable caching:

```yaml
- uses: actions/setup-go@v5
  with:
    go-version: '1.25'
    cache: true
```

### Build matrix too large

Limit combinations:

```yaml
exclude:
  - goos: windows
    goarch: arm64
```
