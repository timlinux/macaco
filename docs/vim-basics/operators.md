# Vim Operators

Operators perform actions on text. They typically combine with motions or text objects.

## Delete Operator

| Command | Description |
|---------|-------------|
| `x` | Delete character under cursor |
| `X` | Delete character before cursor |
| `d{motion}` | Delete with motion |
| `dd` | Delete line |
| `D` | Delete to end of line |

**Examples:**

- `dw` - Delete word
- `d$` - Delete to end of line
- `d2w` - Delete 2 words
- `dt)` - Delete until ')'

## Change Operator

| Command | Description |
|---------|-------------|
| `s` | Substitute character (delete and insert) |
| `S` | Substitute line |
| `c{motion}` | Change with motion |
| `cc` | Change line |
| `C` | Change to end of line |

**Examples:**

- `cw` - Change word
- `ciw` - Change inner word
- `c$` - Change to end of line
- `cf,` - Change through comma

## Yank (Copy) Operator

| Command | Description |
|---------|-------------|
| `y{motion}` | Yank with motion |
| `yy` or `Y` | Yank line |

**Examples:**

- `yw` - Yank word
- `y$` - Yank to end of line
- `yiw` - Yank inner word

## Put (Paste) Operator

| Command | Description |
|---------|-------------|
| `p` | Put after cursor |
| `P` | Put before cursor |

## Replace

| Command | Description |
|---------|-------------|
| `r{char}` | Replace single character |
| `R` | Enter replace mode |

## Operator + Motion Formula

The general formula is:

```
[count] operator [count] motion
```

Examples:

- `2dw` - Delete 2 words
- `d2w` - Also delete 2 words
- `3cw` - Change 3 words
