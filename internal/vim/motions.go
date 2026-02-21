package vim

import (
	"unicode"
	"unicode/utf8"
)

// Motion represents a vim motion that moves the cursor
type Motion interface {
	Execute(b *Buffer, count int) bool
}

// CharClass represents character classification for word motions
type CharClass int

const (
	CharClassWhitespace CharClass = iota
	CharClassWord
	CharClassPunctuation
)

func classifyChar(r rune) CharClass {
	if unicode.IsSpace(r) {
		return CharClassWhitespace
	}
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
		return CharClassWord
	}
	return CharClassPunctuation
}

// MoveLeft moves cursor left
func MoveLeft(b *Buffer, count int) bool {
	moved := false
	for i := 0; i < count; i++ {
		if b.cursorX > 0 {
			b.cursorX--
			moved = true
		}
	}
	return moved
}

// MoveRight moves cursor right
func MoveRight(b *Buffer, count int) bool {
	moved := false
	line := b.CurrentLine()
	lineLen := utf8.RuneCountInString(line)

	for i := 0; i < count; i++ {
		maxX := lineLen - 1
		if b.mode == ModeInsert {
			maxX = lineLen
		}
		if b.cursorX < maxX {
			b.cursorX++
			moved = true
		}
	}
	return moved
}

// MoveUp moves cursor up
func MoveUp(b *Buffer, count int) bool {
	moved := false
	for i := 0; i < count; i++ {
		if b.cursorY > 0 {
			b.cursorY--
			moved = true
		}
	}
	b.clampCursor()
	return moved
}

// MoveDown moves cursor down
func MoveDown(b *Buffer, count int) bool {
	moved := false
	for i := 0; i < count; i++ {
		if b.cursorY < len(b.lines)-1 {
			b.cursorY++
			moved = true
		}
	}
	b.clampCursor()
	return moved
}

// MoveToLineStart moves cursor to start of line
func MoveToLineStart(b *Buffer) bool {
	if b.cursorX != 0 {
		b.cursorX = 0
		return true
	}
	return false
}

// MoveToLineEnd moves cursor to end of line
func MoveToLineEnd(b *Buffer) bool {
	line := b.CurrentLine()
	lineLen := utf8.RuneCountInString(line)
	newX := lineLen - 1
	if b.mode == ModeInsert {
		newX = lineLen
	}
	if newX < 0 {
		newX = 0
	}
	if b.cursorX != newX {
		b.cursorX = newX
		return true
	}
	return false
}

// MoveToFirstNonBlank moves cursor to first non-blank character
func MoveToFirstNonBlank(b *Buffer) bool {
	line := b.CurrentLine()
	runes := []rune(line)
	for i, r := range runes {
		if !unicode.IsSpace(r) {
			if b.cursorX != i {
				b.cursorX = i
				return true
			}
			return false
		}
	}
	return MoveToLineStart(b)
}

// MoveWordForward moves cursor to start of next word
func MoveWordForward(b *Buffer, count int) bool {
	moved := false
	for i := 0; i < count; i++ {
		if moveWordForwardOnce(b) {
			moved = true
		}
	}
	return moved
}

func moveWordForwardOnce(b *Buffer) bool {
	line := b.CurrentLine()
	runes := []rune(line)
	startX := b.cursorX

	// Skip current character class
	if b.cursorX < len(runes) {
		currentClass := classifyChar(runes[b.cursorX])
		for b.cursorX < len(runes) && classifyChar(runes[b.cursorX]) == currentClass {
			b.cursorX++
		}
	}

	// Skip whitespace
	for b.cursorX < len(runes) && unicode.IsSpace(runes[b.cursorX]) {
		b.cursorX++
	}

	// If at end of line, try next line
	if b.cursorX >= len(runes) && b.cursorY < len(b.lines)-1 {
		b.cursorY++
		b.cursorX = 0
		// Skip leading whitespace on new line
		line = b.CurrentLine()
		runes = []rune(line)
		for b.cursorX < len(runes) && unicode.IsSpace(runes[b.cursorX]) {
			b.cursorX++
		}
	}

	b.clampCursor()
	return b.cursorX != startX || b.cursorY != startX
}

// MoveWordBackward moves cursor to start of previous word
func MoveWordBackward(b *Buffer, count int) bool {
	moved := false
	for i := 0; i < count; i++ {
		if moveWordBackwardOnce(b) {
			moved = true
		}
	}
	return moved
}

func moveWordBackwardOnce(b *Buffer) bool {
	startX, startY := b.cursorX, b.cursorY

	// If at start of line, go to previous line
	if b.cursorX == 0 {
		if b.cursorY > 0 {
			b.cursorY--
			line := b.CurrentLine()
			b.cursorX = utf8.RuneCountInString(line)
		}
	}

	line := b.CurrentLine()
	runes := []rune(line)

	// Move back one character
	if b.cursorX > 0 {
		b.cursorX--
	}

	// Skip whitespace backward
	for b.cursorX > 0 && unicode.IsSpace(runes[b.cursorX]) {
		b.cursorX--
	}

	// Handle line wrap
	if b.cursorX == 0 && len(runes) > 0 && unicode.IsSpace(runes[0]) && b.cursorY > 0 {
		b.cursorY--
		line = b.CurrentLine()
		runes = []rune(line)
		b.cursorX = len(runes) - 1
		for b.cursorX > 0 && unicode.IsSpace(runes[b.cursorX]) {
			b.cursorX--
		}
	}

	// Find start of word
	if b.cursorX >= 0 && b.cursorX < len(runes) {
		currentClass := classifyChar(runes[b.cursorX])
		for b.cursorX > 0 && classifyChar(runes[b.cursorX-1]) == currentClass {
			b.cursorX--
		}
	}

	if b.cursorX < 0 {
		b.cursorX = 0
	}

	return b.cursorX != startX || b.cursorY != startY
}

// MoveWordEnd moves cursor to end of word
func MoveWordEnd(b *Buffer, count int) bool {
	moved := false
	for i := 0; i < count; i++ {
		if moveWordEndOnce(b) {
			moved = true
		}
	}
	return moved
}

func moveWordEndOnce(b *Buffer) bool {
	line := b.CurrentLine()
	runes := []rune(line)
	startX := b.cursorX

	// Move forward at least one character
	if b.cursorX < len(runes)-1 {
		b.cursorX++
	} else if b.cursorY < len(b.lines)-1 {
		b.cursorY++
		b.cursorX = 0
		line = b.CurrentLine()
		runes = []rune(line)
	}

	// Skip whitespace
	for b.cursorX < len(runes) && unicode.IsSpace(runes[b.cursorX]) {
		b.cursorX++
	}

	// Handle line wrap
	if b.cursorX >= len(runes) && b.cursorY < len(b.lines)-1 {
		b.cursorY++
		b.cursorX = 0
		line = b.CurrentLine()
		runes = []rune(line)
		for b.cursorX < len(runes) && unicode.IsSpace(runes[b.cursorX]) {
			b.cursorX++
		}
	}

	// Move to end of word
	if b.cursorX < len(runes) {
		currentClass := classifyChar(runes[b.cursorX])
		for b.cursorX < len(runes)-1 && classifyChar(runes[b.cursorX+1]) == currentClass {
			b.cursorX++
		}
	}

	b.clampCursor()
	return b.cursorX != startX
}

// MoveToChar moves cursor to next occurrence of character
func MoveToChar(b *Buffer, char rune, count int, before bool) bool {
	line := b.CurrentLine()
	runes := []rune(line)
	found := 0

	for i := b.cursorX + 1; i < len(runes); i++ {
		if runes[i] == char {
			found++
			if found == count {
				if before {
					b.cursorX = i - 1
				} else {
					b.cursorX = i
				}
				return true
			}
		}
	}
	return false
}

// MoveToCharBackward moves cursor to previous occurrence of character
func MoveToCharBackward(b *Buffer, char rune, count int, after bool) bool {
	line := b.CurrentLine()
	runes := []rune(line)
	found := 0

	for i := b.cursorX - 1; i >= 0; i-- {
		if runes[i] == char {
			found++
			if found == count {
				if after {
					b.cursorX = i + 1
				} else {
					b.cursorX = i
				}
				return true
			}
		}
	}
	return false
}

// MoveToBufferStart moves cursor to start of buffer
func MoveToBufferStart(b *Buffer) bool {
	if b.cursorX != 0 || b.cursorY != 0 {
		b.cursorX = 0
		b.cursorY = 0
		return true
	}
	return false
}

// MoveToBufferEnd moves cursor to end of buffer
func MoveToBufferEnd(b *Buffer) bool {
	lastLine := len(b.lines) - 1
	if lastLine < 0 {
		lastLine = 0
	}
	if b.cursorY != lastLine {
		b.cursorY = lastLine
		MoveToFirstNonBlank(b)
		return true
	}
	return false
}

// MoveToLine moves cursor to specific line number
func MoveToLine(b *Buffer, lineNum int) bool {
	// Line numbers are 1-based
	targetY := lineNum - 1
	if targetY < 0 {
		targetY = 0
	}
	if targetY >= len(b.lines) {
		targetY = len(b.lines) - 1
	}
	if b.cursorY != targetY {
		b.cursorY = targetY
		MoveToFirstNonBlank(b)
		return true
	}
	return false
}

// MoveToMatchingBracket moves cursor to matching bracket
func MoveToMatchingBracket(b *Buffer) bool {
	line := b.CurrentLine()
	runes := []rune(line)
	if b.cursorX >= len(runes) {
		return false
	}

	char := runes[b.cursorX]
	pairs := map[rune]rune{
		'(': ')', ')': '(',
		'[': ']', ']': '[',
		'{': '}', '}': '{',
		'<': '>', '>': '<',
	}

	match, ok := pairs[char]
	if !ok {
		return false
	}

	// Determine direction
	forward := char == '(' || char == '[' || char == '{' || char == '<'
	depth := 1

	if forward {
		for y := b.cursorY; y < len(b.lines); y++ {
			line := b.lines[y]
			runes := []rune(line)
			startX := 0
			if y == b.cursorY {
				startX = b.cursorX + 1
			}
			for x := startX; x < len(runes); x++ {
				if runes[x] == char {
					depth++
				} else if runes[x] == match {
					depth--
					if depth == 0 {
						b.cursorY = y
						b.cursorX = x
						return true
					}
				}
			}
		}
	} else {
		for y := b.cursorY; y >= 0; y-- {
			line := b.lines[y]
			runes := []rune(line)
			startX := len(runes) - 1
			if y == b.cursorY {
				startX = b.cursorX - 1
			}
			for x := startX; x >= 0; x-- {
				if runes[x] == char {
					depth++
				} else if runes[x] == match {
					depth--
					if depth == 0 {
						b.cursorY = y
						b.cursorX = x
						return true
					}
				}
			}
		}
	}

	return false
}
