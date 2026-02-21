# MoCaCo - Motion Capture Combatant
## Vim Motion Training Game Specification

Version: 1.0.0
Last Updated: 2026-02-21

---

## Table of Contents

1. [Overview](#overview)
2. [User Stories](#user-stories)
3. [Functional Requirements](#functional-requirements)
4. [Technical Requirements](#technical-requirements)
5. [Business Rules](#business-rules)
6. [Task Database](#task-database)
7. [Architecture](#architecture)
8. [File Structure](#file-structure)
9. [Running Modes](#running-modes)
10. [Management Scripts](#management-scripts)
11. [Stats File Format](#stats-file-format)
12. [REST API Specification](#rest-api-specification)
13. [Visual Style Guide](#visual-style-guide)
14. [Colour Palette Reference](#colour-palette-reference)
15. [Version History](#version-history)

---

## Overview

**MoCaCo** (Motion Capture Combatant) is a competitive vim motion training application designed to transform casual vim users into motion masters through gamified, standardized practice rounds. The application features a unique visual transformation interface where users convert initial text states into desired states using only vim motions, commands, and operators.

### Technology Stack

| Component | Technology |
|-----------|-----------|
| Backend | Go 1.21+ |
| Terminal UI | Bubble Tea, Lipgloss |
| Web Frontend | React 18, Chakra UI, Framer Motion |
| Animation | Spring physics (harmonica) |
| Documentation | MkDocs Material theme |
| Builds | Nix flakes (reproducible) |
| Data Layer | JSON configuration files |

### Core Game Loop

1. User is presented with a transformation task: `initial_text → desired_text`
2. User executes vim commands to transform the initial text into the desired state
3. System provides real-time visual feedback (green for correct, red for incorrect)
4. After 30 tasks, detailed statistics are displayed showing performance metrics
5. Lifetime statistics track improvement over time

### Target Audience

- Vim beginners wanting to learn motions systematically
- Intermediate users seeking to improve speed and muscle memory
- Advanced users competing for leaderboard positions
- Development teams using it as a fun training tool
- Future esports competitors in vim motion competitions

---

## User Stories

### US-001: Basic Training Round
**As a** vim learner
**I want to** complete standardized practice rounds with visual feedback
**So that** I can systematically improve my vim motion skills

### US-002: Performance Tracking
**As a** competitive user
**I want to** see detailed statistics on my performance by motion category
**So that** I can identify weak areas and track improvement over time

### US-003: Visual Learning
**As a** visual learner
**I want to** see large, clear text transformations with immediate color-coded feedback
**So that** I can quickly understand when I'm correct or incorrect

### US-004: Standardized Assessment
**As a** vim trainer
**I want** all users to complete the same set of tasks in each round type
**So that** performance metrics are comparable across users and sessions

### US-005: Multi-Platform Access
**As a** developer
**I want to** practice vim motions in both terminal and web interfaces
**So that** I can train regardless of my current environment

### US-006: Progress Motivation
**As a** long-term user
**I want to** see my lifetime statistics and improvement graphs
**So that** I stay motivated to continue practicing

### US-007: Future Competitive Play
**As a** competitive gamer
**I want to** compete against others in real-time
**So that** I can prove my vim motion mastery (Future feature)

---

## Functional Requirements

### FR-001: Main Game Interface

The application SHALL display three text elements vertically arranged:

- **Previous Task** (top): Dimmed, small text showing the last completed transformation
- **Current Task** (center): Large, prominent display showing:
  - Initial state text (user's current buffer)
  - Separator (e.g., "→" or "|")
  - Desired state text (goal)
- **Next Task** (bottom): Dimmed, small text previewing the upcoming transformation

**Visual Hierarchy:**
- Current task SHALL be at least 3x the size of previous/next tasks
- Current task SHALL be centered both horizontally and vertically
- Previous/next tasks SHALL use 50% opacity
- The separator SHALL be clearly visible between initial and desired states

### FR-002: Text Transformation Mechanics

The application SHALL operate in a vim-like modal editing environment:

- Users start in NORMAL mode
- All standard vim motions SHALL be supported:
  - Movement: `h`, `j`, `k`, `l`, `w`, `b`, `e`, `0`, `$`, `^`, `gg`, `G`, `f`, `t`, `%`, etc.
  - Operations: `d`, `c`, `y`, `p`, `x`, `r`, `s`, etc.
  - Text objects: `iw`, `aw`, `i"`, `a(`, `it`, etc.
  - Counts: `3w`, `2dd`, `5j`, etc.
  - Combinations: `dt)`, `ci"`, `daw`, `d$`, etc.

- Mode transitions SHALL work as in vim:
  - `i`, `I`, `a`, `A`, `o`, `O` enter INSERT mode
  - `v`, `V`, `Ctrl-v` enter VISUAL modes
  - `ESC` returns to NORMAL mode
  - `:` enters COMMAND mode

- The system SHALL track the exact key sequence used to complete each task
- Timing SHALL begin when the task is displayed
- Timing SHALL end when the buffer exactly matches the desired state

### FR-003: Real-Time Visual Feedback

The application SHALL provide immediate visual feedback on the current buffer state:

- **Correct State (Complete Match):**
  - Buffer text SHALL turn bright green (#10B981)
  - A success indicator SHALL appear (e.g., checkmark icon)
  - After 500ms delay, automatically advance to next task

- **Incorrect State (No Match):**
  - Buffer text SHALL turn red (#EF4444) when modifications are made
  - Red SHALL persist until the buffer returns to correct state or initial state

- **In-Progress State (Partial Match):**
  - Text SHALL be yellow/orange (#F59E0B) to indicate changes in progress
  - Individual characters MAY be highlighted differently to show diff

**Motion Task Display:**
  - For motion-only tasks (where initial text == desired text), display SHALL include:
    - Current cursor position indicated with block character (█) in the buffer text
    - A caret (^) indicator directly below the text showing target cursor position
    - Clear instruction text: "Move your cursor to the caret (^)"
  - This provides clear visual guidance for where the cursor must be moved

**Editing Task Display:**
  - For editing tasks (where text transformation is required), display SHALL include:
    - Current buffer text with cursor indicator
    - Arrow separator (↓) between current and desired state
    - Desired text shown in success color
    - Clear instruction text: "Transform the text above to match the text below"

**Previous/Next Task Indicators:**
  - Previous task indicator SHALL be blank when on the first task of a round
  - Next task indicator SHALL show the actual next task in the queue (not the current task)
  - For motion tasks in previews, indicate "(move cursor)" instead of showing identical text

**Target Text Highlighting:**
  - For delete and change tasks, the target text to be modified SHALL be highlighted
  - Highlighted text SHALL use warning/orange color (#F59E0B) with bold and underline
  - Highlighting SHALL only appear while buffer matches initial state (not yet modified)
  - This provides visual pattern recognition for faster task comprehension
  - Task struct includes `HighlightStart` and `HighlightEnd` fields to define highlight range

- **Character-Level Diff Highlighting:**
  The application SHALL provide character-level color coding to help users quickly identify what needs to be done:

  - **Delete (Red):** Characters that need to be deleted SHALL be rendered in red (`#EF4444`)
    - Bold text for visibility
    - Indicates text that must be removed to match desired state

  - **Change (Orange):** Characters that need to be modified SHALL be rendered in orange (`#F59E0B`)
    - Bold text for visibility
    - Indicates text that must be changed (e.g., case change, character substitution)

  - **Target (Green):** The cursor target position SHALL be rendered in green (`#10B981`)
    - Bold text with underline for emphasis
    - Used for motion tasks to show where the cursor needs to move
    - Only the single target character is highlighted

  - **Highlighting Behavior:**
    - Character highlighting SHALL only appear while buffer matches initial state
    - Once user begins editing, highlighting is disabled (buffer differs from initial)
    - This provides visual guidance for pattern recognition before action
    - Highlighting is calculated by comparing initial text with desired text character-by-character

### FR-004: Standardized Round Structure

Each practice round SHALL consist of exactly 30 tasks distributed as follows:

| Category | Count | Description |
|----------|-------|-------------|
| Motion | 6 | Pure movement tasks (no editing) |
| Delete | 6 | Deletion operations (d, x, dd, D) |
| Change | 6 | Change operations (c, s, C, S, r) |
| Insert | 6 | Insertion operations (i, a, o, I, A, O) |
| Visual | 3 | Visual mode selections |
| Complex | 3 | Multi-step combinations |

**Task Selection:**
- Tasks SHALL be deterministically ordered within each round type
- Round types SHALL include: "Beginner", "Intermediate", "Advanced", "Mixed"
- Each round type SHALL have predefined task sets
- Task order within categories MAY be randomized for variety

### FR-005: Task Difficulty Progression

Tasks SHALL be categorized by difficulty level:

**Beginner (Level 1):**
- Single motion commands: `w`, `b`, `e`, `$`, `0`
- Simple deletions: `x`, `dd`, `dw`
- Basic insertions: `i`, `a`, `o`

**Intermediate (Level 2):**
- Motion with counts: `3w`, `2b`, `4j`
- Text object operations: `daw`, `ciw`, `di"`
- Line-based changes: `cc`, `C`, `D`

**Advanced (Level 3):**
- Complex combinations: `dt)`, `cf,`, `y$p`
- Search-based motions: `f`, `t`, `F`, `T`
- Multiple operations: `d2w`, `c3j`, `>G`

**Expert (Level 4):**
- Multi-step transformations requiring 3+ commands
- Macro-like efficiency challenges
- Uncommon but powerful combinations

### FR-006: Performance Metrics Tracking

The application SHALL track and persist the following metrics:

**Per-Task Metrics:**
- Task ID and category
- Time to completion (milliseconds)
- Key sequence used
- Number of keystrokes
- Optimal keystroke count (for efficiency calculation)
- Success/failure status
- Timestamp of attempt

**Per-Category Metrics:**
- Average completion time
- Best completion time
- Average keystroke efficiency (actual/optimal ratio)
- Success rate (percentage)
- Total tasks attempted
- Total tasks completed

**Session Metrics:**
- Session start/end timestamps
- Total tasks completed
- Overall average time
- Overall keystroke efficiency
- Category breakdown
- Session best times

**Lifetime Metrics:**
- All session metrics aggregated
- Personal best times per task
- Personal best times per category
- Total practice time
- Improvement trends (stored as time-series data)
- Mastery levels per category (0-100 score)

### FR-007: Statistics Display

After completing a round (30 tasks), the application SHALL display:

**Summary Screen:**
- Total time for round
- Average time per task
- Overall keystroke efficiency
- Tasks completed vs attempted
- Grade/score (S, A, B, C, D, F based on performance)

**Category Breakdown:**
- Table showing each category with:
  - Average time
  - Best time
  - Efficiency percentage
  - Success rate
  - Visual indicator (progress bar or stars)

**Detailed Task List:**
- Scrollable list of all 30 tasks showing:
  - Task number and category
  - Initial → Desired transformation
  - Time taken
  - Keystrokes used vs optimal
  - Success indicator

**Improvement Suggestions:**
- Top 3 categories needing improvement
- Specific tips for each weak category
- Recommended practice focus

**Lifetime Comparison:**
- "Personal Best" indicators on any new records
- Comparison to previous session averages
- Improvement percentage over last 5 sessions

### FR-008: Terminal User Interface

The terminal interface SHALL use Bubble Tea and Lipgloss to provide:

**Layout:**
- Full-screen terminal application
- Header showing: mode, task counter (e.g., "15/30"), timer
- Main content area with three-tier task display
- Footer showing: current mode, hints, keyboard shortcuts

**Colors:**
- 256-color mode support
- Gradient backgrounds for visual appeal
- Color-coded feedback as per FR-003
- Theme options: "Dark", "Light", "High Contrast"

**Typography:**
- Block letters using Unicode box drawing characters (█, ▀, ▄)
- Monospace font for code/text display
- Large text rendering for current task (minimum 5 rows high)

**Animations:**
- Smooth transitions between tasks (slide up effect)
- Spring physics for success indicators
- Pulsing effect on current task
- Character-by-character reveal for hints

### FR-009: Web User Interface

The web interface SHALL use React 18, Chakra UI, and Framer Motion to provide:

**Responsive Design:**
- Desktop-first design (1920x1080 optimal)
- Tablet support (768px+ width)
- Mobile support (future consideration)

**Layout Components:**
- Header: Logo, user stats summary, settings button
- Main area: Three-tier task display
- Sidebar: Progress tracker, category breakdown (optional, collapsible)
- Footer: Mode indicator, keyboard visualization (optional)

**Animations:**
- Framer Motion for task transitions
- Spring physics for success/failure feedback
- Smooth color transitions
- Confetti effect on completing round

**Accessibility:**
- Keyboard-only operation
- Screen reader support for stats
- High contrast mode
- Adjustable text sizes

### FR-010: Task Database Structure

Tasks SHALL be stored in a structured JSON format:

```json
{
  "tasks": [
    {
      "id": "motion-001",
      "category": "motion",
      "difficulty": 1,
      "initial": "Hello world from vim",
      "desired": "Hello world from vim",
      "cursor_start": 0,
      "cursor_end": 6,
      "optimal_keys": "w",
      "optimal_count": 1,
      "description": "Move to next word",
      "hint": "Use 'w' to jump to next word",
      "tags": ["word-motion", "basic"]
    }
  ]
}
```

**Required Fields:**
- `id`: Unique identifier
- `category`: One of [motion, delete, change, insert, visual, complex]
- `difficulty`: Integer 1-4
- `initial`: Starting text
- `desired`: Goal text
- `cursor_start`: Initial cursor position (character index)
- `cursor_end`: Expected final cursor position (for motion tasks)
- `optimal_keys`: Most efficient key sequence
- `optimal_count`: Number of keystrokes in optimal solution

**Optional Fields:**
- `description`: Human-readable task description
- `hint`: Help text for user
- `tags`: Array of tags for categorization
- `alternative_solutions`: Array of other valid solutions

### FR-011: Vim Engine Integration

The application SHALL embed a vim-compatible text editing engine:

**Options:**
1. **Neovim Embedded:** Use neovim's embedded mode via RPC
2. **Vim.js:** JavaScript vim implementation for web
3. **Custom Engine:** Minimal vim motion parser in Go

**Requirements:**
- Must support all standard vim motions and operations
- Must provide buffer state access for comparison
- Must support mode tracking (NORMAL, INSERT, VISUAL, etc.)
- Must provide keystroke capture
- Must support undo/redo for task resets
- Should support visual mode for visual tasks

**Integration:**
- Backend SHALL manage vim engine instance per session
- Frontend SHALL send keystrokes to backend
- Backend SHALL return buffer state and cursor position
- Comparison logic SHALL run server-side

### FR-012: Session Management

The application SHALL support multiple concurrent sessions:

**Session Lifecycle:**
1. User starts new session → Session ID generated
2. Round type selected → Tasks loaded
3. Tasks completed sequentially → Progress tracked
4. Round completes → Statistics saved
5. Session ends → Resources cleaned up

**Session Data:**
- Session ID (UUID)
- User ID (for future multi-user support)
- Round type (beginner/intermediate/advanced/mixed)
- Current task index
- Start timestamp
- Buffer state
- Keystroke history
- Task completion status array

**Persistence:**
- Sessions SHALL be saved to disk every 5 tasks
- Incomplete sessions SHALL be recoverable on restart
- Sessions older than 24 hours MAY be auto-deleted

### FR-013: Keyboard Shortcuts

The application SHALL support the following global shortcuts:

| Shortcut | Action |
|----------|--------|
| `Ctrl+C` | Quit application (with confirmation) |
| `Ctrl+R` | Reset current task |
| `Ctrl+S` | Skip current task (counts as failure) |
| `Ctrl+H` | Show hint for current task |
| `Ctrl+P` | Pause/resume timer |
| `?` | Show help overlay |
| `ESC` | Close overlays/modals |

**In-Task Shortcuts:**
- All vim commands SHALL work as expected
- `u` SHALL undo last change
- `Ctrl+r` SHALL redo
- `:q` SHALL exit task (with confirmation)

### FR-014: Help System

The application SHALL provide contextual help:

**Help Overlay:**
- Keyboard shortcuts reference
- Current mode explanation
- Available vim commands for current task
- Link to full documentation

**Task Hints:**
- Progressive hints (3 levels):
  1. Category hint: "This is a word motion task"
  2. Motion hint: "Use the 'w' command"
  3. Full solution: Show optimal key sequence

**Vim Command Reference:**
- Searchable command database
- Categorized by type (motion/operator/text-object)
- Examples for each command
- Interactive practice mode

### FR-015: Settings and Preferences

Users SHALL be able to configure:

**Appearance:**
- Color theme (dark/light/high-contrast)
- Font size multiplier
- Animation speed
- Show/hide hints

**Gameplay:**
- Auto-advance delay (0-2000ms)
- Enable/disable timer pressure
- Show/hide optimal solution
- Task shuffle within categories

**Accessibility:**
- Screen reader mode
- Reduced motion
- High contrast mode
- Keyboard sound effects (on/off)

**Data:**
- Reset statistics (with confirmation)
- Export statistics (JSON/CSV)
- Import statistics
- Delete session history

### FR-016: Round Types

The application SHALL support predefined round types:

**Beginner Round:**
- 30 tasks, all difficulty level 1
- Focuses on: basic motions (h/j/k/l/w/b/e), simple deletions (x/dd), basic insertions (i/a/o)
- No time pressure
- Hints enabled by default

**Intermediate Round:**
- 30 tasks, difficulty levels 1-2
- Introduces: counts (3w/2dd), text objects (iw/aw/i"), line operations (C/D)
- Moderate time pressure
- Hints available on request

**Advanced Round:**
- 30 tasks, difficulty levels 2-3
- Complex operations: dt)/cf,/ci", search motions (f/t/F/T), combinations
- High time pressure
- Limited hints

**Expert Round:**
- 30 tasks, difficulty levels 3-4
- Multi-step transformations, efficiency focus
- Timed with penalties
- No hints

**Mixed Round:**
- 30 tasks, random difficulty 1-4
- Balanced distribution
- Realistic practice
- Adaptive hints

**Daily Challenge:**
- New round generated daily
- Same tasks for all users globally
- Leaderboard support
- Fixed difficulty curve

### FR-017: Leaderboards (Future)

The application SHALL support competitive leaderboards:

**Global Leaderboards:**
- Top 100 by round type
- Sorted by: speed, efficiency, or combined score
- Display: rank, username, time, efficiency, date

**Category Leaderboards:**
- Best times per category
- Efficiency rankings
- Consistency rankings (low variance)

**Daily Challenge Leaderboard:**
- Resets daily at UTC midnight
- Requires completion within 24 hours
- Prevents replays (one attempt per day)

### FR-018: Multiplayer Support (Future)

The application SHALL support real-time competitive play:

**Match Types:**
- 1v1: Head-to-head race through same 30 tasks
- Tournament: Bracket-style elimination
- Co-op: Team completion (alternating tasks)

**Match Flow:**
1. Players join lobby
2. Round type selected
3. Countdown timer (3-2-1-GO)
4. Both players complete same tasks
5. Live progress displayed for both players
6. Winner determined by first to complete or best time
7. Post-match statistics comparison

**Network Protocol:**
- WebSocket connections for real-time state
- Server authoritative task validation
- Anti-cheat: server-side buffer comparison
- Latency compensation

### FR-019: Achievement System

The application SHALL award achievements for milestones:

**Achievement Categories:**

**Mastery:**
- "Motion Master": Complete 100 motion tasks
- "Delete Dominator": Complete 100 delete tasks
- "Change Champion": Complete 100 change tasks
- "Insert Expert": Complete 100 insert tasks
- "Visual Virtuoso": Complete 50 visual tasks
- "Combo King": Complete 50 complex tasks

**Speed:**
- "Speed Demon": Complete task in under 1 second
- "Lightning Fast": Complete round in under 60 seconds
- "Consistent": Complete 10 tasks in a row under 3s each

**Efficiency:**
- "Optimal Path": Use exact optimal keystrokes 10 times
- "No Waste": Complete 20 tasks with ≤110% optimal keystrokes
- "Perfect Efficiency": Complete round with 100% average efficiency

**Progression:**
- "First Steps": Complete first round
- "Dedicated": Complete 10 rounds
- "Expert": Complete 100 rounds
- "Legend": Complete 1000 rounds

**Special:**
- "Flawless Victory": Complete round without any mistakes
- "Comeback": Complete round after failing 50% of tasks
- "Persistent": Reset task 10 times before succeeding

### FR-020: Documentation Site

A MkDocs documentation site SHALL be maintained with:

**Content Structure:**
- Getting Started
  - Installation
  - Quick Start Tutorial
  - First Round Walkthrough
- Vim Basics
  - Modes Overview
  - Essential Motions
  - Operators and Text Objects
  - Combining Commands
- Game Mechanics
  - Round Types
  - Scoring System
  - Statistics Explained
- Advanced Topics
  - Efficiency Optimization
  - Common Patterns
  - Expert Strategies
- API Documentation
  - REST API Reference
  - WebSocket Protocol
  - Data Formats
- Contributing
  - Development Setup
  - Architecture Overview
  - Adding New Tasks

**Features:**
- Search functionality
- Dark/light theme toggle
- Code syntax highlighting
- Interactive vim command examples
- Embedded video tutorials (future)

---

## Technical Requirements

### TR-001: Programming Language and Frameworks

**Backend:**
- The application SHALL be written in Go 1.21 or higher
- The backend SHALL use the following libraries:
  - Bubble Tea (https://github.com/charmbracelet/bubbletea) for TUI
  - Lipgloss (https://github.com/charmbracelet/lipgloss) for styling
  - Chi or Gin for REST API routing
  - Gorilla WebSocket for real-time communication (future)

**Frontend:**
- The web application SHALL be written in TypeScript
- The web application SHALL use React 18+
- The web application SHALL use Chakra UI for components
- The web application SHALL use Framer Motion for animations
- The web application SHALL use React Query for API state management

**Documentation:**
- The documentation site SHALL use MkDocs with Material theme
- Documentation SHALL be written in Markdown

### TR-002: Cross-Platform Support

The application SHALL support:

**Operating Systems:**
- Linux (Ubuntu 20.04+, Fedora 35+, Arch)
- macOS (11.0 Big Sur or higher)
- Windows (10 or higher)

**Package Formats:**
- Standalone binaries (Linux, macOS, Windows)
- DEB packages (Debian/Ubuntu)
- RPM packages (Fedora/RHEL)
- Flatpak (universal Linux)
- Homebrew formula (macOS)
- Scoop manifest (Windows)
- Docker image

### TR-003: Build System

The application SHALL use:

- **Nix Flakes** for reproducible builds
- **GitHub Actions** for CI/CD
- **GoReleaser** for release automation
- **Vite** for frontend bundling

**Build Requirements:**
- Single-command build for all targets
- Deterministic builds (same inputs → same outputs)
- Version embedding at build time
- Automated testing before release

### TR-004: Performance Requirements

The application SHALL meet:

**Response Times:**
- Keystroke to visual feedback: <16ms (60 FPS)
- Task transition: <200ms
- Statistics calculation: <100ms
- API response time: <50ms (95th percentile)

**Resource Usage:**
- Terminal: <50MB RAM, <5% CPU (idle)
- Web backend: <100MB RAM per session
- Web frontend: <200MB browser memory

**Concurrency:**
- Support minimum 100 concurrent sessions (server mode)
- Session isolation (no cross-session interference)
- Graceful degradation under load

---

## Business Rules

### BR-001: Task Completion Validation

A task SHALL be considered complete if and only if:
- The current buffer text exactly matches the desired text (character-for-character)
- The comparison is case-sensitive
- Whitespace is significant (spaces, tabs, newlines must match exactly)
- Cursor position is irrelevant for completion (except motion-only tasks)

**For motion-only tasks** (where initial == desired):
- Cursor position SHALL match the expected final position
- Buffer text SHALL remain unchanged

### BR-002: Time Measurement

Timing SHALL be measured as follows:
- Start time: When task is displayed and user can begin typing
- End time: When buffer matches desired state
- Elapsed time = End time - Start time (in milliseconds)
- Paused time SHALL NOT count toward elapsed time
- Timer SHALL pause when help overlay is shown
- Timer SHALL pause when user presses Ctrl+P

### BR-003: Efficiency Calculation

Keystroke efficiency SHALL be calculated as:

```
efficiency_percentage = (optimal_keystrokes / actual_keystrokes) * 100
```

Where:
- `optimal_keystrokes`: Minimum keystrokes needed (from task definition)
- `actual_keystrokes`: Keystrokes user actually used

**Special Cases:**
- If user resets task, keystrokes before reset count toward total
- Undo/redo keystrokes count as regular keystrokes
- Mode transitions count (ESC, i, v, etc.)
- Failed attempts carry over keystroke counts

### BR-004: Grading System

Round performance SHALL be graded using this scale:

| Grade | Requirements |
|-------|-------------|
| S | 100% tasks completed, avg efficiency ≥95%, avg time ≤ target time |
| A | 100% tasks completed, avg efficiency ≥85%, avg time ≤ target time * 1.2 |
| B | ≥90% tasks completed, avg efficiency ≥75% |
| C | ≥75% tasks completed, avg efficiency ≥60% |
| D | ≥50% tasks completed |
| F | <50% tasks completed |

**Target Times** (by difficulty):
- Level 1: 5 seconds
- Level 2: 8 seconds
- Level 3: 12 seconds
- Level 4: 20 seconds

### BR-005: Statistics Persistence

Statistics SHALL be persisted as follows:

- After every completed task, session data SHALL be updated in memory
- Every 5 completed tasks, session SHALL be written to disk
- Upon round completion, full statistics SHALL be written to disk
- Upon application exit, current session SHALL be saved
- Corrupted statistics files SHALL be backed up and reset
- Statistics older than 2 years MAY be archived

**Data Integrity:**
- If average efficiency is >200%, data SHALL be flagged as suspicious
- If best time is <500ms for level 2+ tasks, data SHALL be flagged
- Flagged data SHALL be marked but not deleted

### BR-006: Task Randomization

Within each category, tasks MAY be randomized, but:
- The same round type SHALL always contain the same task IDs
- Task order within a category MAY vary between sessions
- Difficulty progression SHALL be maintained (no level 4 before level 1)
- Each task SHALL appear exactly once per round
- Randomization seed SHALL be stored for round replay

### BR-007: Skip and Reset Behavior

**Skip (Ctrl+S):**
- Marks task as failed
- Records time as "maximum time" (60 seconds)
- Records efficiency as 0%
- Advances to next task immediately
- User can skip maximum 5 tasks per round

**Reset (Ctrl+R):**
- Reverts buffer to initial state
- Resets cursor to starting position
- Does NOT reset timer
- Does NOT reset keystroke count
- Unlimited resets allowed

---

## Task Database

### Procedural Task Generation

Tasks are procedurally generated using sample text from **public domain literature**. This approach provides:

- **Variety:** Each session presents unique text combinations
- **Scalability:** Unlimited task generation without manual authoring
- **Legal Compliance:** All source texts are in the public domain (published pre-1928)

#### Text Sources

All texts are sourced from [Project Gutenberg](https://www.gutenberg.org):

| Work | Author | Year | License |
|------|--------|------|---------|
| Pride and Prejudice | Jane Austen | 1813 | Public Domain |
| A Tale of Two Cities | Charles Dickens | 1859 | Public Domain |
| The Adventures of Sherlock Holmes | Arthur Conan Doyle | 1892 | Public Domain |
| Moby Dick | Herman Melville | 1851 | Public Domain |
| Alice's Adventures in Wonderland | Lewis Carroll | 1865 | Public Domain |

#### Generation Algorithm

The task generator:
1. Selects random sentences from the public domain corpus
2. Applies vim operation templates based on task category
3. Calculates optimal keystroke sequences
4. Assigns difficulty based on operation complexity
5. Ensures 30 tasks per round (6 motion, 6 delete, 6 change, 6 insert, 3 visual, 3 complex)

See `LICENSE_TEXTS.md` for full attribution and licensing details.

### Task Categories and Examples

#### Motion Tasks (6 per round)

Motion tasks require moving the cursor without modifying text.

**Example Tasks:**

```json
{
  "id": "motion-w-001",
  "category": "motion",
  "difficulty": 1,
  "initial": "hello world from vim",
  "desired": "hello world from vim",
  "cursor_start": 0,
  "cursor_end": 6,
  "optimal_keys": "w",
  "optimal_count": 1,
  "description": "Move to next word start",
  "hint": "Use 'w' to move forward one word"
}
```

```json
{
  "id": "motion-3w-001",
  "category": "motion",
  "difficulty": 2,
  "initial": "one two three four five",
  "desired": "one two three four five",
  "cursor_start": 0,
  "cursor_end": 14,
  "optimal_keys": "3w",
  "optimal_count": 2,
  "description": "Move forward three words",
  "hint": "Use a count before motion: '3w'"
}
```

```json
{
  "id": "motion-f-001",
  "category": "motion",
  "difficulty": 2,
  "initial": "find the letter x in this line",
  "cursor_start": 0,
  "cursor_end": 16,
  "optimal_keys": "fx",
  "optimal_count": 2,
  "description": "Find next occurrence of 'x'",
  "hint": "Use 'f' followed by target character"
}
```

#### Delete Tasks (6 per round)

Delete tasks require removing text using vim delete operations.

```json
{
  "id": "delete-x-001",
  "category": "delete",
  "difficulty": 1,
  "initial": "helxlo world",
  "desired": "hello world",
  "cursor_start": 3,
  "optimal_keys": "x",
  "optimal_count": 1,
  "description": "Delete single character under cursor",
  "hint": "Use 'x' to delete character under cursor"
}
```

```json
{
  "id": "delete-dw-001",
  "category": "delete",
  "difficulty": 1,
  "initial": "hello extra world",
  "desired": "hello world",
  "cursor_start": 6,
  "optimal_keys": "dw",
  "optimal_count": 2,
  "description": "Delete word forward",
  "hint": "Use 'dw' to delete word"
}
```

```json
{
  "id": "delete-daw-001",
  "category": "delete",
  "difficulty": 2,
  "initial": "delete this word here",
  "desired": "delete word here",
  "cursor_start": 10,
  "optimal_keys": "daw",
  "optimal_count": 3,
  "description": "Delete a word (text object)",
  "hint": "Use 'daw' to delete 'a word' including surrounding space"
}
```

#### Change Tasks (6 per round)

Change tasks require modifying text using vim change operations.

```json
{
  "id": "change-cw-001",
  "category": "change",
  "difficulty": 1,
  "initial": "hello old world",
  "desired": "hello new world",
  "cursor_start": 6,
  "optimal_keys": "cwnew<ESC>",
  "optimal_count": 6,
  "description": "Change word to 'new'",
  "hint": "Use 'cw' to change word, type new text, press ESC"
}
```

```json
{
  "id": "change-ciw-001",
  "category": "change",
  "difficulty": 2,
  "initial": "change inside word",
  "desired": "change outside word",
  "cursor_start": 10,
  "optimal_keys": "ciwoutside<ESC>",
  "optimal_count": 11,
  "description": "Change inner word",
  "hint": "Use 'ciw' to change inside word (doesn't require cursor at start)"
}
```

```json
{
  "id": "change-ci-001",
  "category": "change",
  "difficulty": 2,
  "initial": "text = \"old value\"",
  "desired": "text = \"new value\"",
  "cursor_start": 10,
  "optimal_keys": "ci\"new value<ESC>",
  "optimal_count": 14,
  "description": "Change inside quotes",
  "hint": "Use 'ci\"' to change text inside double quotes"
}
```

#### Insert Tasks (6 per round)

Insert tasks require adding text using vim insert operations.

```json
{
  "id": "insert-i-001",
  "category": "insert",
  "difficulty": 1,
  "initial": "hello world",
  "desired": "hello beautiful world",
  "cursor_start": 6,
  "optimal_keys": "ibeautiful <ESC>",
  "optimal_count": 12,
  "description": "Insert before cursor",
  "hint": "Use 'i' to enter insert mode before cursor"
}
```

```json
{
  "id": "insert-A-001",
  "category": "insert",
  "difficulty": 1,
  "initial": "hello world",
  "desired": "hello world!",
  "cursor_start": 0,
  "optimal_keys": "A!<ESC>",
  "optimal_count": 4,
  "description": "Append at end of line",
  "hint": "Use 'A' to append at end of line (capital A)"
}
```

```json
{
  "id": "insert-o-001",
  "category": "insert",
  "difficulty": 1,
  "initial": "line one\nline three",
  "desired": "line one\nline two\nline three",
  "cursor_start": 8,
  "optimal_keys": "oline two<ESC>",
  "optimal_count": 11,
  "description": "Open new line below",
  "hint": "Use 'o' to open a new line below and enter insert mode"
}
```

#### Visual Tasks (3 per round)

Visual tasks require using visual mode for selection and operations.

```json
{
  "id": "visual-vd-001",
  "category": "visual",
  "difficulty": 2,
  "initial": "select and delete this text here",
  "desired": "select and delete here",
  "cursor_start": 19,
  "optimal_keys": "viwwd",
  "optimal_count": 5,
  "description": "Visually select word and delete",
  "hint": "Use 'viw' to visually select inner word, then 'd' to delete"
}
```

```json
{
  "id": "visual-Vy-001",
  "category": "visual",
  "difficulty": 2,
  "initial": "copy this line\n\npaste here",
  "desired": "copy this line\ncopy this line\npaste here",
  "cursor_start": 0,
  "optimal_keys": "Vyjp",
  "optimal_count": 4,
  "description": "Line-visual yank and paste",
  "hint": "Use 'V' for line-visual, 'y' to yank, 'j' to move down, 'p' to paste"
}
```

#### Complex Tasks (3 per round)

Complex tasks require multiple operations or advanced combinations.

```json
{
  "id": "complex-dt-001",
  "category": "complex",
  "difficulty": 3,
  "initial": "delete until (keep this)",
  "desired": "(keep this)",
  "cursor_start": 0,
  "optimal_keys": "dt(",
  "optimal_count": 3,
  "description": "Delete until character",
  "hint": "Use 'dt(' to delete until the '(' character"
}
```

```json
{
  "id": "complex-multi-001",
  "category": "complex",
  "difficulty": 4,
  "initial": "const oldName = 'value';",
  "desired": "const newName = 'value';",
  "cursor_start": 6,
  "optimal_keys": "cewnewName<ESC>",
  "optimal_count": 12,
  "description": "Change variable name",
  "hint": "Navigate to the word and use 'cw' to change it"
}
```

### Task Database Structure

The task database SHALL be stored in `/data/tasks.json`:

```json
{
  "version": "1.0.0",
  "last_updated": "2026-02-21",
  "rounds": {
    "beginner": {
      "name": "Beginner Round",
      "difficulty_range": [1, 1],
      "task_distribution": {
        "motion": 6,
        "delete": 6,
        "change": 6,
        "insert": 6,
        "visual": 3,
        "complex": 3
      },
      "tasks": [
        "motion-w-001",
        "motion-b-001",
        "delete-x-001",
        "..."
      ]
    },
    "intermediate": {
      "name": "Intermediate Round",
      "difficulty_range": [1, 2],
      "tasks": ["..."]
    },
    "advanced": {
      "name": "Advanced Round",
      "difficulty_range": [2, 3],
      "tasks": ["..."]
    },
    "expert": {
      "name": "Expert Round",
      "difficulty_range": [3, 4],
      "tasks": ["..."]
    }
  },
  "tasks": [
    {
      "id": "motion-w-001",
      "category": "motion",
      "difficulty": 1,
      "initial": "hello world from vim",
      "desired": "hello world from vim",
      "cursor_start": 0,
      "cursor_end": 6,
      "optimal_keys": "w",
      "optimal_count": 1,
      "description": "Move to next word start",
      "hint": "Use 'w' to move forward one word",
      "tags": ["word-motion", "basic", "essential"]
    }
  ]
}
```

---

## Architecture

### System Architecture

MoCaCo follows a client-server architecture with three operational modes:

```
┌─────────────────────────────────────────────────────────────┐
│                     MoCaCo Architecture                      │
└─────────────────────────────────────────────────────────────┘

┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│   Terminal   │         │     Web      │         │    Mobile    │
│   Client     │◄───────►│   Client     │◄───────►│   Client     │
│ (Bubble Tea) │   HTTP  │   (React)    │   HTTP  │   (Future)   │
└──────────────┘   REST  └──────────────┘   REST  └──────────────┘
       │                        │                        │
       │                        │                        │
       └────────────────────────┼────────────────────────┘
                                │
                                ▼
                    ┌───────────────────────┐
                    │    Game API Layer     │
                    │   (REST + WebSocket)  │
                    └───────────────────────┘
                                │
                ┌───────────────┼───────────────┐
                ▼               ▼               ▼
       ┌────────────┐  ┌────────────┐  ┌────────────┐
       │  Session   │  │    Vim     │  │   Stats    │
       │  Manager   │  │   Engine   │  │  Tracker   │
       └────────────┘  └────────────┘  └────────────┘
                ▼               ▼               ▼
       ┌────────────┐  ┌────────────┐  ┌────────────┐
       │  Session   │  │    Task    │  │   Stats    │
       │   Store    │  │  Database  │  │    File    │
       │  (JSON)    │  │  (JSON)    │  │  (JSON)    │
       └────────────┘  └────────────┘  └────────────┘
```

### Component Descriptions

**Terminal Client (Bubble Tea):**
- Renders three-tier task display
- Captures keyboard input
- Sends keystrokes to backend
- Receives buffer state updates
- Displays statistics

**Web Client (React):**
- Responsive UI with Chakra components
- Sends keystrokes via REST API
- Polls for buffer state (or WebSocket)
- Animated transitions with Framer Motion
- Interactive statistics dashboard

**Game API Layer:**
- RESTful HTTP endpoints for session management
- WebSocket support for real-time updates (future)
- Authentication and session validation
- Rate limiting and abuse prevention

**Session Manager:**
- Creates and destroys sessions
- Manages concurrent sessions
- Tracks session state
- Enforces session timeouts

**Vim Engine:**
- Executes vim commands on buffer
- Maintains buffer state
- Tracks cursor position
- Supports all vim motions/operations
- Provides undo/redo

**Stats Tracker:**
- Records per-task metrics
- Aggregates session statistics
- Calculates lifetime metrics
- Generates improvement trends

**Data Stores:**
- Session Store: Active session data (in-memory + disk)
- Task Database: All task definitions (read-only)
- Stats File: User statistics (append-only)

### Data Flow

**New Session Flow:**
1. Client requests new session with round type
2. Session Manager creates session with UUID
3. Task Database loads 30 tasks for round type
4. Session initialized with first task
5. Client receives session ID and first task

**Keystroke Processing Flow:**
1. Client captures keystroke
2. Client sends keystroke to backend via API
3. Vim Engine processes keystroke
4. Buffer state updated
5. Comparison logic checks if buffer == desired
6. Response includes: buffer state, cursor position, match status
7. Client renders updated buffer with color coding

**Task Completion Flow:**
1. Buffer matches desired state
2. Stats Tracker records: time, keystrokes, efficiency
3. Session advances to next task
4. Client shows success animation (500ms)
5. Next task displayed automatically

**Round Completion Flow:**
1. Task 30 completed
2. Stats Tracker calculates session statistics
3. Stats persisted to disk
4. Lifetime stats updated
5. Client receives full statistics
6. Statistics screen displayed
7. Session marked as complete

### Running Modes

**Mode 1: Combined (Default)**
```
┌──────────────────────────────┐
│  Single Process              │
│  ┌────────────┐              │
│  │  Terminal  │              │
│  │   Client   │              │
│  └─────┬──────┘              │
│        │                     │
│        ▼                     │
│  ┌────────────┐              │
│  │  Backend   │              │
│  │  (In-proc) │              │
│  └────────────┘              │
└──────────────────────────────┘
```

**Mode 2: Server Only**
```
┌──────────────────────────────┐
│  Background Process          │
│  ┌────────────┐              │
│  │  Backend   │              │
│  │  Server    │◄─────────────┼─── HTTP :8080
│  └────────────┘              │
└──────────────────────────────┘
```

**Mode 3: Client Only**
```
┌──────────────────────────────┐
│  Terminal/Web Client         │
│  ┌────────────┐              │
│  │   Client   │──────────────┼─── HTTP → Remote Server
│  └────────────┘              │
└──────────────────────────────┘
```

---

## File Structure

```
macaco/
├── cmd/
│   ├── macaco/              # Main CLI entry point
│   │   └── main.go
│   ├── server/              # Server-only mode
│   │   └── main.go
│   └── client/              # Client-only mode
│       └── main.go
├── internal/
│   ├── game/                # Core game logic
│   │   ├── session.go       # Session management
│   │   ├── task.go          # Task definitions
│   │   └── round.go         # Round logic
│   ├── vim/                 # Vim engine integration
│   │   ├── engine.go        # Vim command processor
│   │   ├── buffer.go        # Buffer state management
│   │   └── motions.go       # Motion implementations
│   ├── stats/               # Statistics tracking
│   │   ├── tracker.go       # Metrics collection
│   │   ├── aggregator.go    # Statistics aggregation
│   │   └── persistence.go   # File I/O
│   ├── api/                 # REST API
│   │   ├── server.go        # HTTP server
│   │   ├── handlers.go      # Request handlers
│   │   └── middleware.go    # Auth, logging, etc.
│   ├── tui/                 # Terminal UI (Bubble Tea)
│   │   ├── app.go           # Main TUI application
│   │   ├── game_view.go     # Game screen
│   │   ├── stats_view.go    # Statistics screen
│   │   └── styles.go        # Lipgloss styles
│   └── config/              # Configuration
│       ├── config.go        # Config loading
│       └── defaults.go      # Default values
├── web/                     # React frontend
│   ├── src/
│   │   ├── components/
│   │   │   ├── GameView/
│   │   │   │   ├── TaskDisplay.tsx
│   │   │   │   ├── BufferEditor.tsx
│   │   │   │   └── ProgressBar.tsx
│   │   │   ├── StatsView/
│   │   │   │   ├── SessionStats.tsx
│   │   │   │   ├── LifetimeStats.tsx
│   │   │   │   └── CategoryBreakdown.tsx
│   │   │   └── Shared/
│   │   │       ├── Header.tsx
│   │   │       └── Footer.tsx
│   │   ├── hooks/
│   │   │   ├── useGameSession.ts
│   │   │   ├── useKeyboard.ts
│   │   │   └── useStats.ts
│   │   ├── services/
│   │   │   └── api.ts       # API client
│   │   ├── theme/
│   │   │   └── index.ts     # Chakra theme
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
├── data/
│   ├── tasks.json           # Task database
│   └── stats.json           # User statistics (generated)
├── docs/                    # MkDocs documentation
│   ├── index.md
│   ├── getting-started/
│   ├── vim-basics/
│   ├── game-mechanics/
│   ├── advanced/
│   └── api/
├── scripts/
│   ├── build.sh             # Build script
│   ├── test.sh              # Test runner
│   ├── dev.sh               # Development server
│   └── release.sh           # Release automation
├── flake.nix                # Nix flake for reproducible builds
├── flake.lock
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── SPECIFICATION.md         # This file
├── LICENSE
└── mkdocs.yml               # Documentation config
```

---

## Running Modes

### Mode 1: Combined (Default)

Start the application with both backend and TUI client in a single process:

```bash
macaco
```

**Use Cases:**
- Local practice sessions
- No network required
- Simplest setup
- Fastest performance (no HTTP overhead)

**Behavior:**
- Backend runs in-process
- TUI client communicates via Go interfaces (not HTTP)
- Stats saved to `~/.config/macaco/stats.json`
- PID file: `~/.config/macaco/macaco.pid`

### Mode 2: Server Only

Start the backend server without any client:

```bash
macaco server --port 8080
```

**Use Cases:**
- Serve web clients
- Remote practice sessions
- Multi-user support
- Centralized statistics

**Behavior:**
- HTTP server on specified port
- REST API available for clients
- Runs as background daemon
- PID file: `~/.config/macaco/server.pid`
- Logs: `~/.config/macaco/server.log`

### Mode 3: Client Only (Terminal)

Connect TUI client to remote server:

```bash
macaco client --server http://localhost:8080
```

**Use Cases:**
- Connect to remote server
- Practice on shared infrastructure
- Compete with others

**Behavior:**
- TUI client only
- Communicates via REST API
- No local backend
- Stats stored on server

### Mode 4: Web Client

Serve the React web application:

```bash
# Development
cd web && npm run dev

# Production
macaco server --port 8080 --serve-web
```

Access at: `http://localhost:8080`

**Use Cases:**
- Browser-based practice
- No terminal required
- Mobile support (future)
- Easier onboarding

---

## Management Scripts

### Build Script

`scripts/build.sh` - Build all components

```bash
#!/usr/bin/env bash
set -e

echo "Building MoCaCo..."

# Build backend
go build -o bin/macaco cmd/macaco/main.go
go build -o bin/macaco-server cmd/server/main.go
go build -o bin/macaco-client cmd/client/main.go

# Build frontend
cd web
npm run build
cd ..

echo "Build complete!"
echo "Binaries: bin/"
echo "Web dist: web/dist/"
```

### Development Script

`scripts/dev.sh` - Run in development mode

```bash
#!/usr/bin/env bash

# Run backend with hot reload
air -c .air.toml &

# Run frontend dev server
cd web && npm run dev &

wait
```

### Test Script

`scripts/test.sh` - Run all tests

```bash
#!/usr/bin/env bash
set -e

# Backend tests
go test ./... -v -cover

# Frontend tests
cd web
npm test
cd ..

# Integration tests
go test ./tests/integration -v
```

### Release Script

`scripts/release.sh` - Create release packages

```bash
#!/usr/bin/env bash
set -e

VERSION=$1
if [ -z "$VERSION" ]; then
  echo "Usage: ./release.sh <version>"
  exit 1
fi

# Use GoReleaser
goreleaser release --clean

echo "Release $VERSION created!"
```

---

## Stats File Format

User statistics are stored in `~/.config/macaco/stats.json`:

```json
{
  "version": "1.0.0",
  "user_id": "uuid-here",
  "created_at": "2026-02-21T10:00:00Z",
  "last_updated": "2026-02-21T11:30:00Z",

  "lifetime": {
    "total_rounds": 25,
    "total_tasks": 750,
    "total_time_ms": 450000,
    "total_keystrokes": 3500,
    "total_practice_time_ms": 1800000,

    "by_category": {
      "motion": {
        "tasks_attempted": 150,
        "tasks_completed": 145,
        "total_time_ms": 90000,
        "total_keystrokes": 500,
        "best_time_ms": 800,
        "avg_time_ms": 620,
        "avg_efficiency": 92.5,
        "success_rate": 96.7
      },
      "delete": { "..." },
      "change": { "..." },
      "insert": { "..." },
      "visual": { "..." },
      "complex": { "..." }
    },

    "personal_bests": {
      "fastest_task": {
        "task_id": "motion-w-001",
        "time_ms": 450,
        "date": "2026-02-20T14:30:00Z"
      },
      "best_efficiency": {
        "task_id": "delete-daw-001",
        "efficiency": 100.0,
        "date": "2026-02-21T09:15:00Z"
      },
      "fastest_round": {
        "round_type": "beginner",
        "time_ms": 120000,
        "date": "2026-02-21T10:00:00Z"
      }
    }
  },

  "sessions": [
    {
      "session_id": "session-uuid",
      "round_type": "beginner",
      "started_at": "2026-02-21T10:00:00Z",
      "completed_at": "2026-02-21T10:15:00Z",
      "total_time_ms": 900000,
      "tasks_completed": 30,
      "tasks_attempted": 30,
      "grade": "A",
      "avg_efficiency": 87.3,

      "tasks": [
        {
          "task_id": "motion-w-001",
          "category": "motion",
          "difficulty": 1,
          "time_ms": 1200,
          "keystrokes": 1,
          "optimal_keystrokes": 1,
          "efficiency": 100.0,
          "success": true,
          "keys_used": "w",
          "resets": 0,
          "hints_used": 0
        },
        {
          "task_id": "delete-x-001",
          "category": "delete",
          "difficulty": 1,
          "time_ms": 2500,
          "keystrokes": 3,
          "optimal_keystrokes": 1,
          "efficiency": 33.3,
          "success": true,
          "keys_used": "llx",
          "resets": 1,
          "hints_used": 0
        }
      ],

      "category_summary": {
        "motion": {
          "avg_time_ms": 1500,
          "avg_efficiency": 95.0,
          "success_rate": 100.0
        }
      }
    }
  ],

  "achievements": [
    {
      "id": "first-steps",
      "name": "First Steps",
      "description": "Complete your first round",
      "unlocked_at": "2026-02-21T10:15:00Z"
    },
    {
      "id": "speed-demon",
      "name": "Speed Demon",
      "description": "Complete a task in under 1 second",
      "unlocked_at": "2026-02-21T10:30:00Z"
    }
  ],

  "preferences": {
    "theme": "dark",
    "auto_advance_delay_ms": 500,
    "show_hints": true,
    "enable_sounds": false,
    "animation_speed": 1.0
  }
}
```

### Statistics Validation

Upon loading stats file, the application SHALL validate:

- `version` field matches supported versions
- All timestamps are valid ISO 8601 format
- Numeric values are within reasonable ranges
- Session task counts match actual task arrays
- Category summaries match individual task data

**If validation fails:**
- Back up corrupted file to `stats.json.backup.<timestamp>`
- Create fresh stats file
- Log error with details
- Notify user of data reset

---

## REST API Specification

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

**Future Implementation:** JWT tokens

**Current:** No authentication (local-only)

### Endpoints

#### Sessions

**Create New Session**

```
POST /sessions
```

Request:
```json
{
  "round_type": "beginner",
  "user_id": "optional-user-id"
}
```

Response (201 Created):
```json
{
  "session_id": "uuid",
  "round_type": "beginner",
  "total_tasks": 30,
  "current_task_index": 0,
  "started_at": "2026-02-21T10:00:00Z",
  "current_task": {
    "task_id": "motion-w-001",
    "category": "motion",
    "difficulty": 1,
    "initial": "hello world from vim",
    "desired": "hello world from vim",
    "cursor_start": 0,
    "description": "Move to next word start",
    "hint": "Use 'w' to move forward one word"
  }
}
```

**Get Session**

```
GET /sessions/:session_id
```

Response (200 OK):
```json
{
  "session_id": "uuid",
  "round_type": "beginner",
  "total_tasks": 30,
  "current_task_index": 5,
  "started_at": "2026-02-21T10:00:00Z",
  "current_task": { "..." },
  "buffer_state": "hello world from vim",
  "cursor_position": 6,
  "elapsed_time_ms": 30000
}
```

**Delete Session**

```
DELETE /sessions/:session_id
```

Response (204 No Content)

#### Keystrokes

**Send Keystroke**

```
POST /sessions/:session_id/keystroke
```

Request:
```json
{
  "key": "w",
  "modifiers": [],
  "timestamp": "2026-02-21T10:00:01.234Z"
}
```

Response (200 OK):
```json
{
  "buffer_state": "hello world from vim",
  "cursor_position": 6,
  "current_mode": "normal",
  "match_status": "complete",
  "task_completed": false,
  "elapsed_time_ms": 1234
}
```

**Match Status Values:**
- `"none"`: Buffer hasn't been modified yet
- `"in_progress"`: Buffer modified but doesn't match desired
- `"complete"`: Buffer matches desired state (for motion tasks, cursor position must also match)

**Send Keystroke Batch**

```
POST /sessions/:session_id/keystrokes
```

Request:
```json
{
  "keys": ["d", "d"],
  "timestamp": "2026-02-21T10:00:01.234Z"
}
```

#### Tasks

**Complete Current Task**

```
POST /sessions/:session_id/complete
```

Request:
```json
{
  "time_ms": 1234,
  "keystrokes": 1,
  "keys_used": "w"
}
```

Response (200 OK):
```json
{
  "task_completed": true,
  "next_task": {
    "task_id": "delete-x-001",
    "category": "delete",
    "..."
  },
  "tasks_remaining": 29,
  "round_complete": false
}
```

**Skip Current Task**

```
POST /sessions/:session_id/skip
```

Response (200 OK):
```json
{
  "task_skipped": true,
  "next_task": { "..." },
  "tasks_remaining": 29
}
```

**Reset Current Task**

```
POST /sessions/:session_id/reset
```

Response (200 OK):
```json
{
  "task_reset": true,
  "buffer_state": "initial state",
  "cursor_position": 0,
  "elapsed_time_ms": 5432
}
```

#### Statistics

**Get Session Statistics**

```
GET /sessions/:session_id/stats
```

Response (200 OK):
```json
{
  "session_id": "uuid",
  "round_type": "beginner",
  "total_time_ms": 90000,
  "tasks_completed": 30,
  "tasks_attempted": 30,
  "grade": "A",
  "avg_efficiency": 87.3,
  "category_summary": { "..." },
  "tasks": [ "..." ]
}
```

**Get Lifetime Statistics**

```
GET /stats/lifetime
```

Response (200 OK):
```json
{
  "total_rounds": 25,
  "total_tasks": 750,
  "by_category": { "..." },
  "personal_bests": { "..." },
  "achievements": [ "..." ]
}
```

**Export Statistics**

```
GET /stats/export?format=json
GET /stats/export?format=csv
```

Response: File download

#### Tasks Database

**Get All Tasks**

```
GET /tasks
```

Response (200 OK):
```json
{
  "tasks": [ "..." ],
  "total": 180
}
```

**Get Task by ID**

```
GET /tasks/:task_id
```

Response (200 OK):
```json
{
  "task_id": "motion-w-001",
  "category": "motion",
  "..."
}
```

**Get Round Definition**

```
GET /rounds/:round_type
```

Response (200 OK):
```json
{
  "round_type": "beginner",
  "name": "Beginner Round",
  "difficulty_range": [1, 1],
  "task_distribution": { "..." },
  "tasks": [ "..." ]
}
```

#### Health

**Health Check**

```
GET /health
```

Response (200 OK):
```json
{
  "status": "ok",
  "version": "1.0.0",
  "uptime_seconds": 3600
}
```

### Error Responses

All errors follow this format:

```json
{
  "error": {
    "code": "SESSION_NOT_FOUND",
    "message": "Session with ID 'xyz' not found",
    "details": {}
  }
}
```

**Common Error Codes:**
- `SESSION_NOT_FOUND` (404)
- `INVALID_ROUND_TYPE` (400)
- `INVALID_KEYSTROKE` (400)
- `TASK_NOT_FOUND` (404)
- `SESSION_EXPIRED` (410)
- `RATE_LIMIT_EXCEEDED` (429)
- `INTERNAL_ERROR` (500)

---

## Visual Style Guide

### Terminal UI (TUI) Style Guide

#### Color Scheme

**Dark Theme (Default):**

| Element | Foreground | Background | Notes |
|---------|-----------|------------|-------|
| Background | - | `#0a0e14` | Deep dark blue-gray |
| Normal Text | `#b3b1ad` | - | Light gray |
| Current Task | `#ffffff` | `#1a1f29` | Bright white on dark |
| Previous/Next Task | `#6b7280` | - | Dimmed gray (50% opacity) |
| Correct (Match) | `#10b981` | - | Bright green |
| Incorrect (No Match) | `#ef4444` | - | Bright red |
| In Progress | `#f59e0b` | - | Amber/orange |
| Separator | `#6366f1` | - | Indigo |
| Header/Footer | `#9ca3af` | `#111827` | Gray on darker gray |
| Highlight/Selection | `#fbbf24` | `#1f2937` | Yellow on dark |

**Light Theme:**

| Element | Foreground | Background | Notes |
|---------|-----------|------------|-------|
| Background | - | `#f9fafb` | Very light gray |
| Normal Text | `#1f2937` | - | Dark gray |
| Current Task | `#000000` | `#ffffff` | Black on white |
| Correct | `#059669` | - | Dark green |
| Incorrect | `#dc2626` | - | Dark red |
| In Progress | `#d97706` | - | Dark amber |

**High Contrast Theme:**

| Element | Foreground | Background | Notes |
|---------|-----------|------------|-------|
| Background | - | `#000000` | Pure black |
| Normal Text | `#ffffff` | - | Pure white |
| Correct | `#00ff00` | - | Pure green |
| Incorrect | `#ff0000` | - | Pure red |

#### Typography

**Block Letters (Current Task):**
- Use Unicode box drawing: `█ ▀ ▄ ▌ ▐ ░ ▒ ▓`
- Minimum 5 rows tall for readability
- 3D effect with shading

Example:
```
███████╗ ██╗  ██╗ ███████╗
██╔════╝ ██║  ██║ ██╔════╝
█████╗   ███████║ █████╗
██╔══╝   ██╔══██║ ██╔══╝
███████╗ ██║  ██║ ███████╗
╚══════╝ ╚═╝  ╚═╝ ╚══════╝
```

**Monospace Font:**
- All text content uses terminal default monospace
- Code examples: `JetBrains Mono`, `Fira Code`, or `Cascadia Code` recommended
- Consistent spacing for alignment

#### Layout Structure

```
╔════════════════════════════════════════════════════════════════╗
║  MoCaCo | Beginner Round | Task 15/30 | Time: 00:45 | [NORMAL] ║
╠════════════════════════════════════════════════════════════════╣
║                                                                ║
║          Previous: "hello old" → "hello new" ✓                 ║
║                                                                ║
║                                                                ║
║         ██████╗ ██╗   ██╗ ██████╗  ██████╗ ███████╗          ║
║        ██╔════╝ ██║   ██║ ██╔══██╗ ██╔══██╗██╔════╝          ║
║        ██║      ██║   ██║ ██████╔╝ ██████╔╝█████╗            ║
║        ██║      ██║   ██║ ██╔══██╗ ██╔══██╗██╔══╝            ║
║        ╚██████╗ ╚██████╔╝ ██║  ██║ ██║  ██║███████╗          ║
║         ╚═════╝  ╚═════╝  ╚═╝  ╚═╝ ╚═╝  ╚═╝╚══════╝          ║
║                                                                ║
║                    "cursor here" → "cursor move"               ║
║                                                                ║
║                                                                ║
║            Next: "delete this" → "delete that"                 ║
║                                                                ║
╠════════════════════════════════════════════════════════════════╣
║  Mode: NORMAL | Ctrl+H: Help | Ctrl+R: Reset | Ctrl+S: Skip   ║
╚════════════════════════════════════════════════════════════════╝
```

**Proportions:**
- Header: 1 row
- Top margin: 2 rows
- Previous task: 1 row
- Spacing: 2 rows
- Current task: 6 rows (large text)
- Separator/transformation: 1 row
- Spacing: 2 rows
- Next task: 1 row
- Bottom margin: 2 rows
- Footer: 1 row (MUST be anchored to bottom of terminal)
- **Total:** Minimum 19 rows (fits in 24-row terminal)

**Footer Anchoring:**
- The footer SHALL always be positioned at the bottom of the terminal
- Content area SHALL flex to fill available space
- Footer position is calculated dynamically based on terminal height

#### Animations

**Task Transition:**
- Current task slides up to become previous (200ms ease-out)
- Next task slides up to become current (200ms ease-out)
- New next task fades in from bottom (100ms)

**Success Animation:**
- Text pulses green (2 pulses, 250ms each)
- Checkmark ✓ appears and bounces (spring physics)
- Subtle confetti burst (optional)
- Auto-advance after 500ms

**Failure Feedback:**
- Text flashes red (single flash, 100ms)
- Shake animation (3 shakes, 50ms each)
- Returns to orange/yellow after 200ms

**Typing Feedback:**
- Character appears with slight scale animation
- Cursor blinks at 530ms interval (vim default)
- Smooth color transitions (100ms)

#### Visual Feedback States

**Idle State (No Input Yet):**
```
Initial text    │    Desired text
cursor here     →    cursor move
```
Colors: Normal text (gray), Separator (indigo)

**In Progress (Typing):**
```
cursor here     │    cursor move
      ^
```
Colors: Modified text (amber), Unchanged (gray)

**Correct (Complete Match):**
```
cursor move ✓
```
Colors: All text (bright green), Checkmark animation

**Incorrect (Wrong State):**
```
cursre move ✗
    ^
```
Colors: All text (red), Error indicator under wrong character

#### Statistics Screen Layout

```
╔════════════════════════════════════════════════════════════════╗
║                    ROUND COMPLETE!                             ║
║                                                                ║
║              Grade: A  |  Time: 02:30  |  Score: 87            ║
╠════════════════════════════════════════════════════════════════╣
║                                                                ║
║  SUMMARY                                                       ║
║  ────────────────────────────────────────────────────          ║
║    Tasks Completed:    30/30  (100%)                          ║
║    Average Time:       5.0s                                    ║
║    Average Efficiency: 87.3%                                   ║
║    Personal Bests:     3 🏆                                    ║
║                                                                ║
║  CATEGORY BREAKDOWN                                            ║
║  ────────────────────────────────────────────────────          ║
║    Motion      [████████████████░░] 91%  (4.2s avg)          ║
║    Delete      [███████████████░░░] 85%  (5.1s avg)          ║
║    Change      [█████████████░░░░░] 82%  (5.8s avg) 🎯       ║
║    Insert      [██████████████████] 95%  (3.9s avg)          ║
║    Visual      [██████████████░░░░] 88%  (6.2s avg)          ║
║    Complex     [███████████░░░░░░░] 75%  (8.5s avg) 🎯       ║
║                                                                ║
║  IMPROVEMENT SUGGESTIONS                                       ║
║  ────────────────────────────────────────────────────          ║
║    🎯 Focus on: Change operations (82% efficiency)            ║
║    💡 Tip: Practice "ci" commands with text objects           ║
║    📚 Learn: dt/df motions for complex tasks                  ║
║                                                                ║
║  LIFETIME STATS                                                ║
║  ────────────────────────────────────────────────────          ║
║    Total Rounds:   25                                          ║
║    Best Grade:     S (Beginner Round)                         ║
║    Fastest Round:  02:15                                       ║
║    Improvement:    +12% from last session                     ║
║                                                                ║
╠════════════════════════════════════════════════════════════════╣
║  Press ENTER to continue | Ctrl+E to export stats              ║
╚════════════════════════════════════════════════════════════════╝
```

---

### Web UI Style Guide

#### Design System

**Typography:**
- Headings: `Inter`, sans-serif
  - H1: 48px, 700 weight
  - H2: 36px, 600 weight
  - H3: 24px, 600 weight
  - H4: 20px, 500 weight
- Body: `Inter`, 16px, 400 weight
- Code/Monospace: `JetBrains Mono`, 14px

**Spacing Scale:**
- 4px, 8px, 12px, 16px, 24px, 32px, 48px, 64px, 96px

**Border Radius:**
- Small: 4px (buttons, inputs)
- Medium: 8px (cards, modals)
- Large: 16px (major sections)
- Full: 9999px (pills, badges)

**Shadows:**
```css
--shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
--shadow-md: 0 4px 6px rgba(0, 0, 0, 0.07);
--shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.1);
--shadow-xl: 0 20px 25px rgba(0, 0, 0, 0.15);
```

#### Color Palette (Chakra Theme)

**Brand Colors:**
```javascript
const colors = {
  brand: {
    50: '#f0f9ff',
    100: '#e0f2fe',
    200: '#bae6fd',
    300: '#7dd3fc',
    400: '#38bdf8',
    500: '#0ea5e9',  // Primary
    600: '#0284c7',
    700: '#0369a1',
    800: '#075985',
    900: '#0c4a6e',
  },
  success: {
    500: '#10b981',
    600: '#059669',
  },
  error: {
    500: '#ef4444',
    600: '#dc2626',
  },
  warning: {
    500: '#f59e0b',
    600: '#d97706',
  }
}
```

#### Component Styles

**Button Variants:**
```jsx
// Primary
<Button colorScheme="brand" size="lg">
  Start Round
</Button>

// Secondary
<Button variant="outline" colorScheme="brand">
  View Stats
</Button>

// Ghost
<Button variant="ghost">
  Skip
</Button>
```

**Card:**
```jsx
<Card
  bg="white"
  borderRadius="lg"
  boxShadow="md"
  p={6}
>
  <CardHeader>
    <Heading size="md">Task 15/30</Heading>
  </CardHeader>
  <CardBody>
    {/* Content */}
  </CardBody>
</Card>
```

**Task Display:**
```jsx
<VStack spacing={8} align="stretch">
  {/* Previous Task */}
  <Text
    fontSize="lg"
    color="gray.500"
    opacity={0.5}
    textAlign="center"
  >
    "hello old" → "hello new" ✓
  </Text>

  {/* Current Task */}
  <Box
    bg="gray.50"
    borderRadius="xl"
    p={12}
    boxShadow="xl"
  >
    <Text
      fontSize="6xl"
      fontFamily="JetBrains Mono"
      fontWeight="bold"
      textAlign="center"
      color={getTaskColor(status)}
      transition="color 0.2s"
    >
      {bufferText}
    </Text>
    <Divider my={4} />
    <Text
      fontSize="2xl"
      color="gray.600"
      textAlign="center"
    >
      → {desiredText}
    </Text>
  </Box>

  {/* Next Task */}
  <Text
    fontSize="lg"
    color="gray.500"
    opacity={0.5}
    textAlign="center"
  >
    "delete this" → "delete that"
  </Text>
</VStack>
```

#### Layout Structure

**Desktop (1920x1080):**
```
┌────────────────────────────────────────────────────┐
│  Header: Logo, Stats, Settings          [User] ⚙  │
├────────────────────────────────────────────────────┤
│                                                    │
│                  Previous Task                     │
│                     (dimmed)                       │
│                                                    │
│  ┌──────────────────────────────────────────────┐ │
│  │                                              │ │
│  │              CURRENT TASK                    │ │
│  │           (large, centered)                  │ │
│  │                                              │ │
│  │         buffer text → desired text           │ │
│  │                                              │ │
│  └──────────────────────────────────────────────┘ │
│                                                    │
│                    Next Task                       │
│                     (dimmed)                       │
│                                                    │
│  Progress: [████████░░░░░░░░░░] 15/30            │
│                                                    │
├────────────────────────────────────────────────────┤
│  Mode: NORMAL | Time: 00:45 | Help (?) | Reset    │
└────────────────────────────────────────────────────┘
```

**Tablet (768px - 1024px):**
- Single column layout
- Slightly smaller current task
- Collapsible sidebar for stats

**Mobile (< 768px) - Future:**
- Vertical stack
- Smaller text sizes
- Touch-optimized controls

#### Animations (Framer Motion)

**Task Transition:**
```jsx
<motion.div
  initial={{ y: 100, opacity: 0 }}
  animate={{ y: 0, opacity: 1 }}
  exit={{ y: -100, opacity: 0 }}
  transition={{
    type: "spring",
    stiffness: 260,
    damping: 20
  }}
>
  {currentTask}
</motion.div>
```

**Success Animation:**
```jsx
<motion.div
  animate={{
    scale: [1, 1.05, 1],
    color: ["#000", "#10b981", "#10b981"]
  }}
  transition={{
    duration: 0.5,
    times: [0, 0.5, 1]
  }}
>
  {taskText} ✓
</motion.div>
```

**Confetti Effect:**
```jsx
import Confetti from 'react-confetti'

<Confetti
  width={windowWidth}
  height={windowHeight}
  recycle={false}
  numberOfPieces={200}
  gravity={0.3}
/>
```

#### Stats Dashboard

**Chart Styles:**
- Use Recharts library
- Consistent color scheme matching brand
- Interactive tooltips
- Responsive sizing

**Progress Bars:**
```jsx
<Progress
  value={efficiency}
  colorScheme={efficiency >= 90 ? "green" : efficiency >= 75 ? "yellow" : "red"}
  size="lg"
  borderRadius="full"
  hasStripe
  isAnimated
/>
```

**Category Cards:**
```jsx
<SimpleGrid columns={{ base: 1, md: 2, lg: 3 }} spacing={6}>
  {categories.map(cat => (
    <Card key={cat.name}>
      <CardHeader>
        <HStack justify="space-between">
          <Heading size="sm">{cat.name}</Heading>
          <Badge colorScheme={cat.badge}>{cat.success}%</Badge>
        </HStack>
      </CardHeader>
      <CardBody>
        <Stat>
          <StatLabel>Avg Time</StatLabel>
          <StatNumber>{cat.avgTime}s</StatNumber>
          <StatHelpText>
            <StatArrow type={cat.trend} />
            {cat.change}%
          </StatHelpText>
        </Stat>
      </CardBody>
    </Card>
  ))}
</SimpleGrid>
```

---

### Documentation Site Style Guide (MkDocs)

#### MkDocs Configuration

```yaml
# mkdocs.yml
site_name: MoCaCo Documentation
site_description: Master vim motions through competitive practice
site_url: https://macaco.dev

theme:
  name: material
  palette:
    # Light mode
    - media: "(prefers-color-scheme: light)"
      scheme: default
      primary: blue
      accent: indigo
      toggle:
        icon: material/brightness-7
        name: Switch to dark mode
    # Dark mode
    - media: "(prefers-color-scheme: dark)"
      scheme: slate
      primary: blue
      accent: indigo
      toggle:
        icon: material/brightness-4
        name: Switch to light mode

  features:
    - navigation.instant
    - navigation.tracking
    - navigation.tabs
    - navigation.sections
    - navigation.expand
    - navigation.top
    - search.suggest
    - search.highlight
    - content.code.copy
    - content.code.annotate

  font:
    text: Inter
    code: JetBrains Mono

  icon:
    repo: fontawesome/brands/github

extra_css:
  - stylesheets/extra.css

extra_javascript:
  - javascripts/extra.js

markdown_extensions:
  - pymdownx.highlight:
      anchor_linenums: true
  - pymdownx.superfences
  - pymdownx.tabbed:
      alternate_style: true
  - pymdownx.keys
  - admonition
  - pymdownx.details
  - attr_list
  - md_in_html
  - pymdownx.emoji:
      emoji_index: !!python/name:material.extensions.emoji.twemoji
      emoji_generator: !!python/name:material.extensions.emoji.to_svg

nav:
  - Home: index.md
  - Getting Started:
    - Installation: getting-started/installation.md
    - Quick Start: getting-started/quick-start.md
    - First Round: getting-started/first-round.md
  - Vim Basics:
    - Modes: vim-basics/modes.md
    - Motions: vim-basics/motions.md
    - Operators: vim-basics/operators.md
    - Text Objects: vim-basics/text-objects.md
  - Game Mechanics:
    - Round Types: game-mechanics/rounds.md
    - Scoring: game-mechanics/scoring.md
    - Statistics: game-mechanics/statistics.md
  - Advanced:
    - Efficiency: advanced/efficiency.md
    - Patterns: advanced/patterns.md
    - Strategies: advanced/strategies.md
  - API:
    - REST API: api/rest.md
    - WebSocket: api/websocket.md
  - Contributing:
    - Development: contributing/development.md
    - Architecture: contributing/architecture.md
```

#### Custom CSS (`docs/stylesheets/extra.css`)

```css
:root {
  --md-primary-fg-color: #0ea5e9;
  --md-accent-fg-color: #6366f1;
}

/* Code blocks with vim command highlighting */
.vim-command {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 0.2em 0.4em;
  border-radius: 4px;
  font-family: 'JetBrains Mono', monospace;
  font-weight: 600;
}

/* Task example boxes */
.task-example {
  border-left: 4px solid var(--md-primary-fg-color);
  background: rgba(14, 165, 233, 0.1);
  padding: 1em;
  margin: 1em 0;
  border-radius: 4px;
}

.task-example .initial {
  color: #6b7280;
  font-family: monospace;
}

.task-example .arrow {
  color: var(--md-primary-fg-color);
  font-weight: bold;
  margin: 0 0.5em;
}

.task-example .desired {
  color: #10b981;
  font-family: monospace;
  font-weight: 600;
}

/* Keyboard key styling */
kbd {
  background: linear-gradient(180deg, #f9fafb 0%, #e5e7eb 100%);
  border: 1px solid #d1d5db;
  border-radius: 4px;
  box-shadow: 0 2px 0 rgba(0, 0, 0, 0.1);
  color: #1f2937;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.9em;
  padding: 0.2em 0.5em;
  white-space: nowrap;
}

/* Stats table styling */
.stats-table {
  width: 100%;
  margin: 1em 0;
}

.stats-table th {
  background: var(--md-primary-fg-color);
  color: white;
  padding: 0.75em;
  text-align: left;
}

.stats-table td {
  padding: 0.75em;
  border-bottom: 1px solid #e5e7eb;
}

/* Achievement badges */
.achievement {
  display: inline-flex;
  align-items: center;
  background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
  color: white;
  padding: 0.5em 1em;
  border-radius: 999px;
  font-weight: 600;
  margin: 0.25em;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.achievement::before {
  content: "🏆";
  margin-right: 0.5em;
}
```

#### Content Style Guide

**Headings:**
- Use sentence case
- Be concise and descriptive
- Use action verbs for tutorials

**Code Examples:**
```markdown
Use vim commands with special formatting:

Press ++w++ to move forward one word.

The command <span class="vim-command">ciw</span> changes the inner word.
```

**Task Examples:**
```markdown
<div class="task-example">
  <strong>Task:</strong> Change the word "old" to "new"<br>
  <span class="initial">hello old world</span>
  <span class="arrow">→</span>
  <span class="desired">hello new world</span><br>
  <strong>Solution:</strong> <code>cwnew<ESC></code>
</div>
```

**Admonitions:**
```markdown
!!! tip "Pro Tip"
    Use counts before motions for faster navigation: `3w` moves forward 3 words.

!!! warning "Common Mistake"
    Don't confuse `d` (delete) with `x` (delete character). Use `d` with motions!

!!! info "Did You Know?"
    The most efficient vim users rarely use arrow keys!
```

**Interactive Elements:**
```markdown
=== "Beginner"
    Start with basic motions: `h`, `j`, `k`, `l`, `w`, `b`

=== "Intermediate"
    Add counts and text objects: `3w`, `ciw`, `daw`

=== "Advanced"
    Master search motions: `f`, `t`, `dt)`, `cf,`
```

---

## Colour Palette Reference

### Terminal Colours (256-color mode)

| Color Name | 256-Color Code | Hex Equivalent | Usage |
|------------|---------------|----------------|-------|
| Background | 234 | `#0a0e14` | Main background |
| Foreground | 250 | `#b3b1ad` | Normal text |
| Bright White | 15 | `#ffffff` | Current task |
| Bright Green | 10 | `#10b981` | Correct/success |
| Bright Red | 9 | `#ef4444` | Incorrect/error |
| Bright Yellow | 11 | `#f59e0b` | In progress |
| Bright Blue | 12 | `#6366f1` | Separator/accent |
| Dim Gray | 242 | `#6b7280` | Previous/next tasks |
| Dark Gray | 236 | `#1a1f29` | Current task bg |

### Brand Palette (Web & Docs)

**Primary (Blue):**
- 50: `#f0f9ff`
- 100: `#e0f2fe`
- 200: `#bae6fd`
- 300: `#7dd3fc`
- 400: `#38bdf8`
- **500: `#0ea5e9`** ← Primary brand color
- 600: `#0284c7`
- 700: `#0369a1`
- 800: `#075985`
- 900: `#0c4a6e`

**Accent (Indigo):**
- 50: `#eef2ff`
- 100: `#e0e7ff`
- 200: `#c7d2fe`
- 300: `#a5b4fc`
- 400: `#818cf8`
- **500: `#6366f1`** ← Accent color
- 600: `#4f46e5`
- 700: `#4338ca`
- 800: `#3730a3`
- 900: `#312e81`

**Semantic Colors:**

| Purpose | Color | Hex |
|---------|-------|-----|
| Success | Green 500 | `#10b981` |
| Error | Red 500 | `#ef4444` |
| Warning | Amber 500 | `#f59e0b` |
| Info | Blue 500 | `#0ea5e9` |

**Character Highlighting Colors:**

| Highlight Type | Color | Hex | Usage |
|----------------|-------|-----|-------|
| Delete | Red 500 | `#ef4444` | Characters to be deleted |
| Change | Amber 500 | `#f59e0b` | Characters to be modified |
| Target | Green 500 | `#10b981` | Cursor target position (motion tasks) |

**Grayscale:**
- 50: `#f9fafb`
- 100: `#f3f4f6`
- 200: `#e5e7eb`
- 300: `#d1d5db`
- 400: `#9ca3af`
- 500: `#6b7280`
- 600: `#4b5563`
- 700: `#374151`
- 800: `#1f2937`
- 900: `#111827`

---

## Version History

### v1.0.0 (Planned)

**Initial Release:**
- Basic game loop with 30-task rounds
- Terminal UI using Bubble Tea + Lipgloss
- Web UI using React + Chakra UI
- Four round types: Beginner, Intermediate, Advanced, Expert
- Six task categories: Motion, Delete, Change, Insert, Visual, Complex
- Real-time visual feedback (green/red/amber)
- Comprehensive statistics tracking
- Session persistence
- Statistics export (JSON/CSV)
- Dark/Light/High-Contrast themes
- MkDocs documentation site
- Cross-platform binaries (Linux, macOS, Windows)
- Package formats: DEB, RPM, Flatpak, Homebrew

**Task Database:**
- 180 total tasks (30 per difficulty level across 6 categories)
- Difficulty progression 1-4
- Optimal solution tracking
- Hints and descriptions

**Future Features (Post-1.0):**
- Multiplayer competitive mode (1v1, tournaments)
- Global leaderboards
- Daily challenges
- Achievement system
- WebSocket real-time updates
- Mobile app (iOS/Android)
- Monetization (premium features, cosmetics)
- Esports tournament platform

---

## License

MIT License - See LICENSE file for details

---

## Contact

- GitHub: https://github.com/timlinux/macaco
- Issues: https://github.com/timlinux/macaco/issues
- Documentation: https://macaco.dev

---

**End of Specification Document**
