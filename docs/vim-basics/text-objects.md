# Text Objects

Text objects select regions of text. They're used with operators like `d`, `c`, `y`.

## Word Objects

| Object | Description |
|--------|-------------|
| `iw` | Inner word |
| `aw` | A word (includes surrounding space) |
| `iW` | Inner WORD |
| `aW` | A WORD |

**Example:**

Text: `hello world here`
Cursor on 'world':

- `diw` results in `hello  here`
- `daw` results in `hello here`

## Quote Objects

| Object | Description |
|--------|-------------|
| `i"` | Inner double quotes |
| `a"` | Including quotes |
| `i'` | Inner single quotes |
| `a'` | Including quotes |
| `` i` `` | Inner backticks |
| `` a` `` | Including backticks |

**Example:**

Text: `name = "value"`
Cursor anywhere inside quotes:

- `ci"` - Change content inside quotes
- `da"` - Delete including quotes

## Bracket Objects

| Object | Description |
|--------|-------------|
| `i(` or `ib` | Inner parentheses |
| `a(` or `ab` | Including parentheses |
| `i[` | Inner brackets |
| `a[` | Including brackets |
| `i{` or `iB` | Inner braces |
| `a{` or `aB` | Including braces |
| `i<` | Inner angle brackets |
| `a<` | Including angle brackets |

**Example:**

Text: `func(old, args)`
Cursor inside parentheses:

- `ci(` - Change everything inside ()
- `di(` - Delete everything inside ()

## Using Text Objects

Text objects are incredibly powerful because they don't depend on cursor position within the object.

With cursor anywhere in "world":

```
hello "world" here
         ^
```

- `ciw` changes just 'world'
- `ci"` changes 'world' (inside quotes)
- `da"` deletes '"world"' (including quotes)

## Best Practices

1. Use `ciw` instead of `bcw` when possible
2. Use text objects for working with quoted strings
3. Use `daw` when you want to remove the word and trailing space
4. Use `ci(` for function arguments
