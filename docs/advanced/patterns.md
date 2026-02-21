# Common Patterns

Master these patterns to handle most vim editing scenarios efficiently.

## Word Operations

### Replace Word

```
Original: hello OLD world
Goal:     hello NEW world
```

With cursor on "OLD": `ciwNEW<Esc>`

### Delete Word and Space

```
Original: hello extra world
Goal:     hello world
```

With cursor on "extra": `daw`

### Swap Words

```
Original: second first
Goal:     first second
```

Method: `dwwP` (delete word, move, put before)

## Line Operations

### Clear Line Content

```
Original: some content here
Goal:
```

Use: `S` or `cc`

### Delete to End

```
Original: keep this delete rest
Goal:     keep this
```

With cursor after "this ": `D`

### Duplicate Line

Method: `yyp` (yank line, paste below)

## Quoted Content

### Change Inside Quotes

```
Original: name = "old"
Goal:     name = "new"
```

With cursor inside quotes: `ci"new<Esc>`

### Delete Including Quotes

```
Original: text "remove" more
Goal:     text  more
```

Use: `da"`

## Bracket Content

### Change Function Arguments

```
Original: func(old, args)
Goal:     func(new, params)
```

With cursor inside parens: `ci(new, params<Esc>`

### Delete Bracketed Content

```
Original: array[remove]
Goal:     array[]
```

Use: `di[`

## Multi-Step Patterns

### Delete Until and Insert

```
Original: prefix_oldtext
Goal:     prefix_newtext
```

Method: `f_lct_newtext<Esc>` or better: `f_lcwnewtext<Esc>`

### Copy and Modify

```
Original: const name = value;
Goal:     const name = value;
          const name2 = value;
```

Method: `yyp` then modify the copy

## Search and Act

### Delete to Next Occurrence

```
Original: hello, world, end
Goal:     hello, end
```

Use: `dt,x` (delete to comma, delete comma)

Or: `df,` (delete through comma, but keeps space wrong)

### Change Through Character

```
Original: start:middle:end
Goal:     start:new:end
```

Use: `f:lcf:new<Esc>`
