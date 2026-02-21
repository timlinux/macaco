# Development Guide

Set up your development environment for MoCaCo.

## Prerequisites

- Go 1.21 or higher
- Git
- (Optional) Nix for reproducible builds

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/timlinux/macaco
cd macaco
```

### Install Dependencies

```bash
go mod download
```

### Build

```bash
make build
```

### Run

```bash
./bin/macaco
```

## Development with Nix

```bash
# Enter development shell
nix develop

# Build with Nix
nix build
```

## Project Structure

```
macaco/
├── cmd/macaco/          # Main entry point
├── internal/
│   ├── api/             # REST API server and client
│   ├── config/          # Configuration management
│   ├── game/            # Core game logic
│   ├── stats/           # Statistics tracking
│   ├── tui/             # Terminal UI
│   └── vim/             # Vim engine
├── data/                # Task database
├── docs/                # Documentation
├── scripts/             # Management scripts
└── pkg/deb/             # DEB package structure
```

## Running Tests

```bash
# All tests
make test

# With coverage
make test-cover
```

## Linting

```bash
# Format code
make fmt

# Run linter
make lint

# Vet code
make vet
```

## Adding Tasks

Tasks are defined in `internal/game/task.go` in the `getEmbeddedTasks()` function.

### Task Structure

```go
{
    ID:           "category-operation-001",
    Category:     CategoryMotion,  // motion, delete, change, insert, visual, complex
    Difficulty:   1,               // 1-4
    Initial:      "initial text",
    Desired:      "desired text",
    CursorStart:  0,
    CursorEnd:    6,               // For motion tasks
    OptimalKeys:  "w",
    OptimalCount: 1,
    Description:  "What the task teaches",
    Hint:         "How to solve it",
    Tags:         []string{"tag1", "tag2"},
}
```

### Guidelines

- Use clear, concise initial/desired text
- Provide accurate optimal solutions
- Write helpful hints
- Tag appropriately for categorization

## Building Releases

```bash
# Build with version info
VERSION=1.0.0 make build-release
```

## Documentation

```bash
# Serve locally
make docs

# Build static site
make docs-build
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions focused and small
- Write clear comments for exported items
