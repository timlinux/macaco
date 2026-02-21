# Quick Start

## Running MoCaCo

Start the application:

```bash
macaco
```

You'll see the main menu:

```
MoCaCo - Motion Capture Combatant
Master vim motions through competitive practice

Select a round type:

  [1] Beginner     - Basic motions and operations
  [2] Intermediate - Counts and text objects
  [3] Advanced     - Complex combinations
  [4] Expert       - Multi-step transformations
  [5] Mixed        - Random difficulty

  [?] Help
  [q] Quit
```

## Starting a Round

Press `1` to start a Beginner round. You'll see:

- The **current task** in the center with your buffer text and the desired text
- A **previous task** above (once you complete one)
- A **next task** preview below
- **Header** showing mode, progress, and timer
- **Footer** showing keyboard shortcuts

## Completing Tasks

1. Look at the transformation: `initial -> desired`
2. Type vim commands to transform the text
3. When the buffer matches the desired text, it turns **green**
4. After a short delay, the next task appears automatically

## Controls

| Key | Action |
|-----|--------|
| `Ctrl+R` | Reset current task |
| `Ctrl+S` | Skip current task |
| `Ctrl+H` | Show/cycle hints |
| `Ctrl+P` | Pause/resume timer |
| `Ctrl+C` | Quit |
| `?` | Show help |

## After the Round

After completing 30 tasks, you'll see your statistics:

- **Grade**: S, A, B, C, D, or F based on performance
- **Summary**: Tasks completed, total time, average efficiency
- **Category breakdown**: Performance by task type
- **Improvement suggestions**: Areas to focus on

Press `Enter` to return to the main menu.
