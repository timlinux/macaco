# Your First Round

This guide walks you through your first Beginner round in MoCaCo.

## Starting Up

1. Launch MoCaCo: `macaco`
2. Press `1` to start a Beginner round

## Understanding the Display

```
MoCaCo | beginner | motion | Task 1/30 | 00:00 | [NORMAL]
────────────────────────────────────────────────────────

                      "previous task" ✓

        ┌─────────────────────────────────────────┐
        │                                         │
        │            hello world█from vim         │
        │                                         │
        │                    ↓                    │
        │                                         │
        │            hello world from vim         │
        │                                         │
        └─────────────────────────────────────────┘

                      "next task preview"

────────────────────────────────────────────────────────
Ctrl+R Reset  |  Ctrl+S Skip  |  Ctrl+H Hint  |  ? Help
```

- **Top line**: Shows round type, category, progress, timer, and current mode
- **Center**: Your buffer (top) and the goal (bottom)
- **█**: The cursor position
- **Bottom**: Available shortcuts

## Example Task: Motion

**Task**: Move the cursor from position 0 to position 6

```
Initial: "hello world from vim" (cursor at 'h')
Desired: "hello world from vim" (cursor at 'w')
```

**Solution**: Press `w` to move to the next word.

When you press `w`, the cursor moves and the task is complete!

## Example Task: Delete

**Task**: Delete the extra 'x' in the word

```
Initial: "helxlo world"
Desired: "hello world"
```

**Solution**: Move to the 'x' and delete it:

1. Press `l` three times to move to 'x'
2. Press `x` to delete it

Or more efficiently: `3lx`

## Example Task: Change

**Task**: Change "old" to "new"

```
Initial: "hello old world"
Desired: "hello new world"
```

**Solution**:

1. Move to "old": `w`
2. Change word: `cw`
3. Type: `new`
4. Press: `Esc`

Complete sequence: `wcwnew<Esc>`

## Tips for Beginners

1. **Don't rush**: Focus on accuracy over speed at first
2. **Use hints**: Press `Ctrl+H` when stuck
3. **Watch the mode**: Make sure you're in the right mode (NORMAL vs INSERT)
4. **Practice resets**: Press `Ctrl+R` to try again without penalty to your timer

## After the Round

You'll see your performance breakdown:

- **Grade**: Based on completion, time, and efficiency
- **Category stats**: See where you're strong and where to improve
- **Suggestions**: Specific tips based on your performance

## Next Steps

1. Try another Beginner round to reinforce basics
2. Move to Intermediate when you're consistently getting A grades
3. Focus on your weakest categories
4. Practice regularly for muscle memory
