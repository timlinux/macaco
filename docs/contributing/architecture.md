# Architecture

MoCaCo follows a modular architecture with clear separation of concerns.

## Overview

```
┌─────────────────────────────────────────────────────────┐
│                      MoCaCo                              │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│  │   TUI App   │    │  REST API   │    │   Config    │ │
│  │  (tui pkg)  │    │  (api pkg)  │    │ (config pkg)│ │
│  └──────┬──────┘    └──────┬──────┘    └─────────────┘ │
│         │                  │                            │
│         └────────┬─────────┘                            │
│                  │                                      │
│         ┌────────▼────────┐                             │
│         │   Game Engine   │                             │
│         │   (game pkg)    │                             │
│         └────────┬────────┘                             │
│                  │                                      │
│    ┌─────────────┼─────────────┐                        │
│    │             │             │                        │
│  ┌─▼───┐    ┌────▼────┐   ┌────▼────┐                  │
│  │ Vim │    │  Tasks  │   │  Stats  │                  │
│  │ Eng │    │   DB    │   │ Tracker │                  │
│  └─────┘    └─────────┘   └─────────┘                  │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

## Packages

### cmd/macaco

Entry point that handles:

- Command-line flag parsing
- Mode selection (combined, server, client)
- Application initialization

### internal/config

Configuration management:

- Loading from file
- Default values
- Environment-aware paths

### internal/vim

Vim engine implementation:

- **buffer.go**: Text buffer with cursor management
- **motions.go**: Movement commands
- **engine.go**: Command parsing and execution

### internal/game

Core game logic:

- **task.go**: Task definitions and database
- **session.go**: Session state management
- **engine.go**: Game engine coordinating all components

### internal/stats

Statistics tracking:

- Per-task metrics
- Session aggregation
- Lifetime statistics
- Achievement tracking
- Persistence

### internal/api

REST API:

- **server.go**: HTTP server and handlers
- **client.go**: HTTP client for remote mode

### internal/tui

Terminal UI:

- **app.go**: Bubble Tea application model
- **styles.go**: Lipgloss styling

## Data Flow

### Combined Mode

```
User Input → TUI → Session → Vim Engine → Buffer
                      ↓
              Stats Tracker → File
```

### Server Mode

```
HTTP Request → API Handler → Game Engine → Response
                                   ↓
                            Stats Tracker
```

### Client Mode

```
User Input → TUI → API Client → HTTP → Server
                      ↓
              Display Response
```

## Key Design Decisions

### 1. Embedded Task Database

Tasks are compiled into the binary for:

- Zero configuration
- Consistent task sets
- Easier distribution

### 2. Local-First Statistics

Stats stored locally in JSON:

- Privacy-respecting
- Works offline
- Easy to export/backup

### 3. Vim Engine Simplicity

Custom vim implementation rather than embedded neovim:

- Smaller binary
- Simpler deployment
- Focused on essential operations

### 4. Mode Separation

Three running modes for flexibility:

- Combined: Simple local use
- Server: Support web clients
- Client: Connect to remote server

## Extension Points

### Adding Commands

1. Add motion/operation to `vim/` package
2. Update engine command parsing
3. Add tests

### Adding Task Categories

1. Define category constant
2. Add tasks to embedded database
3. Update round definitions

### Adding Statistics

1. Extend stats data structures
2. Update tracker aggregation
3. Update persistence format
