# Vim Modes

Vim operates in different modes, each designed for a specific purpose. MoCaCo supports the essential modes.

## Normal Mode

The default mode. Use it for navigation and initiating commands.

- Press `Esc` from any mode to return to Normal mode
- Commands like `d`, `c`, `y` are initiated here
- Motions like `w`, `b`, `e`, `h`, `j`, `k`, `l` work here

## Insert Mode

For typing and editing text.

**Enter Insert Mode:**

| Key | Action |
|-----|--------|
| `i` | Insert before cursor |
| `I` | Insert at line start |
| `a` | Append after cursor |
| `A` | Append at line end |
| `o` | Open line below |
| `O` | Open line above |
| `s` | Substitute character |
| `S` | Substitute line |
| `c{motion}` | Change with motion |

**Exit Insert Mode:**

Press `Esc` to return to Normal mode.

## Visual Mode

For selecting text.

| Key | Action |
|-----|--------|
| `v` | Character-wise visual |
| `V` | Line-wise visual |
| `Ctrl+v` | Block visual |

In Visual mode:

- Use motions to extend selection
- Press `d` to delete selection
- Press `y` to yank (copy) selection
- Press `c` to change selection

## Command Mode

For executing commands (not fully implemented in MoCaCo).

- Press `:` to enter
- Press `Esc` to exit

## Mode Indicator

The current mode is shown in the header:

- `[NORMAL]` - Normal mode
- `[INSERT]` - Insert mode
- `[VISUAL]` - Visual mode
- `[V-LINE]` - Visual line mode
