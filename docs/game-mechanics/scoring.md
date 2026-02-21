# Scoring System

MoCaCo tracks multiple metrics and assigns grades based on your performance.

## Metrics

### Time

- Measured in milliseconds from task display to completion
- Paused time doesn't count
- Skipped tasks receive maximum time (60 seconds)

### Keystrokes

- Every key press counts
- Mode transitions count (`Esc`, `i`, etc.)
- Undo/redo keystrokes count
- Reset doesn't clear keystroke count

### Efficiency

Calculated as:

```
efficiency = (optimal_keystrokes / actual_keystrokes) * 100
```

- 100% = You used the optimal solution
- >100% is capped at 100%
- 0% = Task was skipped

## Grades

| Grade | Requirements |
|-------|--------------|
| **S** | 100% completion, ≥95% efficiency, under target time |
| **A** | 100% completion, ≥85% efficiency, ≤1.2x target time |
| **B** | ≥90% completion, ≥75% efficiency |
| **C** | ≥75% completion, ≥60% efficiency |
| **D** | ≥50% completion |
| **F** | <50% completion |

### Target Times

| Difficulty | Target Time |
|------------|-------------|
| Level 1 | 5 seconds |
| Level 2 | 8 seconds |
| Level 3 | 12 seconds |
| Level 4 | 20 seconds |

## Category Statistics

Each category tracks:

- Tasks attempted
- Tasks completed
- Average time
- Best time
- Average efficiency
- Success rate

## Personal Bests

The system tracks:

- Fastest task completion
- Highest efficiency
- Fastest round completion

Personal bests are highlighted in the statistics screen.

## Improvement Tips

1. **Focus on accuracy first**: Speed comes with muscle memory
2. **Learn optimal solutions**: Press `Ctrl+H` multiple times to see the solution
3. **Practice weak categories**: The stats show where to focus
4. **Use counts and text objects**: They're often more efficient
