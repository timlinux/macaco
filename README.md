# MoCaCo - Motion Capture Combatant

[![Build](https://github.com/timlinux/macaco/actions/workflows/build.yml/badge.svg)](https://github.com/timlinux/macaco/actions/workflows/build.yml)
[![Test](https://github.com/timlinux/macaco/actions/workflows/test.yml/badge.svg)](https://github.com/timlinux/macaco/actions/workflows/test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

**Master vim motions through competitive practice**

MoCaCo is a gamified vim motion training application that transforms casual vim users into motion masters through standardized practice rounds. Complete text transformations using vim commands and track your improvement over time.

## Features

- **30-Task Rounds** - Standardized practice sessions with balanced task categories
- **Real-Time Feedback** - Instant visual feedback as you type (green for correct, red for incorrect)
- **Comprehensive Statistics** - Track performance by category, efficiency metrics, and improvement over time
- **Multiple Difficulty Levels** - From beginner basics to expert-level transformations
- **Terminal UI** - Beautiful TUI built with Bubble Tea and Lipgloss
- **REST API** - Server mode for web clients and remote practice

## Quick Start

```bash
# Install with Go
go install github.com/timlinux/macaco/cmd/macaco@latest

# Or build from source
git clone https://github.com/timlinux/macaco
cd macaco
make build
./bin/macaco
```

## Usage

```bash
# Run in combined mode (default)
macaco

# Start a specific round type
macaco --round intermediate

# Run as server only
macaco --server --addr localhost:8080

# Connect to remote server
macaco --client --addr server.example.com:8080

# Show version
macaco --version
```

## Round Types

| Round | Difficulty | Description |
|-------|------------|-------------|
| Beginner | Level 1 | Basic motions, simple operations |
| Intermediate | Level 1-2 | Counts, text objects |
| Advanced | Level 2-3 | Complex combinations |
| Expert | Level 3-4 | Multi-step transformations |
| Mixed | Level 1-4 | Random difficulty |

## Controls

| Key | Action |
|-----|--------|
| `Ctrl+R` | Reset current task |
| `Ctrl+S` | Skip current task |
| `Ctrl+H` | Show/cycle hints |
| `Ctrl+P` | Pause/resume timer |
| `Ctrl+C` | Quit |
| `?` | Show help |

## Task Categories

Each round includes 30 tasks distributed across:

- **Motion** (6) - Cursor movement without editing
- **Delete** (6) - Deletion operations (d, x, dd, D)
- **Change** (6) - Change operations (c, s, C, S, r)
- **Insert** (6) - Insertion operations (i, a, o, I, A, O)
- **Visual** (3) - Visual mode selections
- **Complex** (3) - Multi-step combinations

## Installation

### From Source

```bash
git clone https://github.com/timlinux/macaco
cd macaco
make build
./bin/macaco
```

### Using Go

```bash
go install github.com/timlinux/macaco/cmd/macaco@latest
```

### Using Nix

```bash
nix run github:timlinux/macaco
```

### Pre-built Binaries

Download from [GitHub Releases](https://github.com/timlinux/macaco/releases).

## Configuration

Statistics and preferences are stored in:

```
~/.config/macaco/stats.json
```

## Documentation

Full documentation available at [https://timlinux.github.io/macaco](https://timlinux.github.io/macaco)

## Development

```bash
# Build
make build

# Test
make test

# Lint
make lint

# Run documentation server
make docs
```

## Project Structure

```
macaco/
├── cmd/macaco/          # Main entry point
├── internal/
│   ├── api/             # REST API server/client
│   ├── config/          # Configuration
│   ├── game/            # Game logic
│   ├── stats/           # Statistics tracking
│   ├── tui/             # Terminal UI
│   └── vim/             # Vim engine
├── docs/                # MkDocs documentation
└── scripts/             # Management scripts
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions welcome! Please read the [development guide](https://timlinux.github.io/macaco/contributing/development/) first.

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
