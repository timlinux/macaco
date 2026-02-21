# MoCaCo - Motion Capture Combatant

**Master vim motions through competitive practice**

MoCaCo is a gamified vim motion training application that transforms casual vim users into motion masters through standardized practice rounds. Complete text transformations using vim commands and track your improvement over time.

## Features

- **30-Task Rounds**: Standardized practice sessions with balanced task categories
- **Real-Time Feedback**: Instant visual feedback as you type - green for correct, red for incorrect
- **Comprehensive Statistics**: Track your performance by category, see efficiency metrics, and monitor improvement
- **Multiple Difficulty Levels**: From beginner basics to expert-level transformations
- **Terminal UI**: Beautiful TUI built with Bubble Tea and Lipgloss
- **REST API**: Server mode for web clients and remote practice

## Quick Start

```bash
# Install with Go
go install github.com/timlinux/macaco/cmd/macaco@latest

# Or build from source
git clone https://github.com/timlinux/macaco
cd macaco
make build

# Run
./bin/macaco
```

## Game Loop

1. You're presented with a transformation task: `initial_text -> desired_text`
2. Use vim commands to transform the initial text into the desired state
3. System provides real-time visual feedback
4. Complete 30 tasks to finish a round
5. Review your statistics and identify areas for improvement

## Round Types

| Round | Difficulty | Description |
|-------|------------|-------------|
| Beginner | Level 1 | Basic motions, simple operations |
| Intermediate | Level 1-2 | Counts, text objects |
| Advanced | Level 2-3 | Complex combinations |
| Expert | Level 3-4 | Multi-step transformations |
| Mixed | Level 1-4 | Random difficulty |

## Task Categories

Each round includes 30 tasks distributed across:

- **Motion** (6): Cursor movement without editing
- **Delete** (6): Deletion operations (d, x, dd, D)
- **Change** (6): Change operations (c, s, C, S, r)
- **Insert** (6): Insertion operations (i, a, o, I, A, O)
- **Visual** (3): Visual mode selections
- **Complex** (3): Multi-step combinations

## Getting Started

Check out the [Installation](getting-started/installation.md) guide to get started, then work through your [First Round](getting-started/first-round.md)!
