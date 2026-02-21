# Vim Motions

Motions move the cursor. They can be used alone or combined with operators.

## Basic Motions

| Motion | Description |
|--------|-------------|
| `h` | Move left |
| `j` | Move down |
| `k` | Move up |
| `l` | Move right |

## Word Motions

| Motion | Description |
|--------|-------------|
| `w` | Next word start |
| `W` | Next WORD start (whitespace-separated) |
| `b` | Previous word start |
| `B` | Previous WORD start |
| `e` | End of word |
| `E` | End of WORD |

## Line Motions

| Motion | Description |
|--------|-------------|
| `0` | Line start |
| `^` | First non-blank character |
| `$` | Line end |
| `g_` | Last non-blank character |

## File Motions

| Motion | Description |
|--------|-------------|
| `gg` | First line |
| `G` | Last line |
| `{n}G` | Go to line n |

## Find Motions

| Motion | Description |
|--------|-------------|
| `f{char}` | Find next occurrence of char |
| `F{char}` | Find previous occurrence |
| `t{char}` | Until next occurrence |
| `T{char}` | Until previous occurrence |
| `;` | Repeat last find |
| `,` | Repeat last find, opposite direction |

## Using Counts

Most motions accept a count prefix:

- `3w` - Move forward 3 words
- `5j` - Move down 5 lines
- `2f,` - Find second comma

## Combining with Operators

Motions define the range for operators:

- `dw` - Delete to next word
- `c$` - Change to end of line
- `y2w` - Yank 2 words

## Motion Tasks in MoCaCo

Motion tasks require moving the cursor without changing text. The initial and desired text are the same - only the cursor position differs.

Example:

```
Initial: "hello world" (cursor at position 0)
Desired: "hello world" (cursor at position 6)
Solution: w
```
