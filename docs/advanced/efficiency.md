# Efficiency Tips

Achieve higher efficiency scores with these techniques.

## General Principles

### 1. Avoid Character-by-Character Movement

Instead of `llll`, use:

- `w` - Next word
- `e` - End of word
- `f{char}` - Find character
- `$` - End of line

### 2. Use Text Objects

Instead of `bde`:

- `diw` - Delete inner word
- `daw` - Delete word with space

### 3. Combine Operations

Instead of `dw` then `i`:

- `cw` - Change word (delete and enter insert mode)

### 4. Use Counts

Instead of `www`:

- `3w` - Move 3 words forward

## Pattern: Change Inside

For quoted strings:

```
text = "old value"
```

Instead of `f"lct"`:

- `ci"` - Change inside quotes (works from anywhere inside)

## Pattern: Delete Until

For removing text up to a character:

```
delete until (keep this)
```

- `dt(` - Delete until '('

## Pattern: Change Through

For changing text including a delimiter:

```
change through, here
```

- `cf,` - Change through comma

## Pattern: Word Replacement

For replacing a word anywhere in it:

```
hello world here
       ^
```

With cursor anywhere in "world":

- `ciw` followed by new word

## Efficiency by Category

### Motion Tasks

- Use word motions over character motions
- Use find (`f`) for jumping to specific characters
- Use `$` and `0` for line ends

### Delete Tasks

- Use `daw` for words with surrounding space
- Use `D` instead of `d$`
- Use `dt{char}` for deleting until

### Change Tasks

- Always prefer `c` over `d` followed by `i`
- Use `ci{object}` for contained text
- Use `C` instead of `c$`

### Insert Tasks

- Use `A` to append at line end
- Use `I` to insert at line start
- Use `o` and `O` for new lines

## Measuring Progress

Track your efficiency percentage:

- 50-70%: Learning the commands
- 70-85%: Good understanding
- 85-95%: Near optimal
- 95-100%: Expert level
