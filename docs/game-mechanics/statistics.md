# Statistics

MoCaCo maintains comprehensive statistics to help you track improvement.

## Statistics Location

Statistics are stored in:

```
~/.config/macaco/stats.json
```

## Session Statistics

After each round, you see:

### Summary

- **Grade**: Overall performance rating
- **Tasks**: Completed vs attempted
- **Time**: Total round time
- **Efficiency**: Average keystroke efficiency

### Category Breakdown

For each category (motion, delete, change, insert, visual, complex):

- Average completion time
- Efficiency percentage
- Success rate
- Progress bar visualization

## Lifetime Statistics

Tracked across all sessions:

### Totals

- Total rounds completed
- Total tasks attempted/completed
- Total practice time
- Total keystrokes

### By Category

- Tasks attempted/completed
- Average time
- Best time
- Average efficiency
- Success rate

### Personal Bests

- Fastest single task
- Best efficiency on a task
- Fastest round completion

## Achievements

Unlock achievements for milestones:

| Achievement | Requirement |
|-------------|-------------|
| First Steps | Complete first round |
| Dedicated | Complete 10 rounds |
| Expert | Complete 100 rounds |
| Optimal Path | Achieve 100% efficiency |
| Flawless Victory | 95%+ efficiency on a round |

## Exporting Statistics

### JSON Export

```bash
# Via API
curl http://localhost:8080/api/v1/stats/export?format=json > stats.json
```

### CSV Export

```bash
curl http://localhost:8080/api/v1/stats/export?format=csv > stats.csv
```

## Resetting Statistics

To start fresh, delete the stats file:

```bash
rm ~/.config/macaco/stats.json
```

Or rename it to keep a backup:

```bash
mv ~/.config/macaco/stats.json ~/.config/macaco/stats.json.backup
```
