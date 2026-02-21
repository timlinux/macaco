package vim

import (
	"strings"
	"unicode"
)

// Engine processes vim commands and manages buffer state
type Engine struct {
	buffer      *Buffer
	undoStack   []*Buffer
	redoStack   []*Buffer
	pendingKeys string
	lastSearch  rune
	searchDir   int // 1 = forward, -1 = backward
	lastMotion  string
}

// NewEngine creates a new vim engine with the given text
func NewEngine(text string) *Engine {
	return &Engine{
		buffer:    NewBuffer(text),
		undoStack: make([]*Buffer, 0),
		redoStack: make([]*Buffer, 0),
	}
}

// Buffer returns the current buffer
func (e *Engine) Buffer() *Buffer {
	return e.buffer
}

// Text returns the current buffer text
func (e *Engine) Text() string {
	return e.buffer.Text()
}

// SetText sets the buffer text
func (e *Engine) SetText(text string) {
	e.saveUndo()
	e.buffer.SetText(text)
}

// CursorPosition returns the cursor position
func (e *Engine) CursorPosition() (x, y int) {
	return e.buffer.CursorPosition()
}

// CursorIndex returns the absolute cursor index
func (e *Engine) CursorIndex() int {
	return e.buffer.CursorIndex()
}

// SetCursorIndex sets the cursor position by index
func (e *Engine) SetCursorIndex(index int) {
	e.buffer.SetCursorIndex(index)
}

// Mode returns the current mode
func (e *Engine) Mode() Mode {
	return e.buffer.Mode()
}

// ProcessKey processes a single key input
func (e *Engine) ProcessKey(key string) bool {
	e.pendingKeys += key

	// Try to parse and execute the pending keys
	consumed, remaining := e.parseAndExecute(e.pendingKeys)
	e.pendingKeys = remaining

	return consumed
}

// parseAndExecute parses pending keys and executes commands
func (e *Engine) parseAndExecute(keys string) (consumed bool, remaining string) {
	if len(keys) == 0 {
		return false, ""
	}

	mode := e.buffer.Mode()

	switch mode {
	case ModeInsert:
		return e.handleInsertMode(keys)
	case ModeNormal:
		return e.handleNormalMode(keys)
	case ModeVisual, ModeVisualLine:
		return e.handleVisualMode(keys)
	default:
		return false, keys
	}
}

// handleInsertMode handles keys in insert mode
func (e *Engine) handleInsertMode(keys string) (bool, string) {
	if len(keys) == 0 {
		return false, ""
	}

	switch keys {
	case "esc", "\x1b":
		e.buffer.SetMode(ModeNormal)
		MoveLeft(e.buffer, 1)
		return true, ""
	case "backspace", "\x7f":
		if e.buffer.cursorX > 0 {
			MoveLeft(e.buffer, 1)
			e.buffer.Delete(1)
		}
		return true, ""
	case "enter", "\r", "\n":
		e.buffer.Insert("\n")
		return true, ""
	default:
		// Regular character input
		if len(keys) == 1 && keys[0] >= 32 {
			e.buffer.Insert(keys)
			return true, ""
		}
		return false, keys
	}
}

// handleNormalMode handles keys in normal mode
func (e *Engine) handleNormalMode(keys string) (bool, string) {
	if len(keys) == 0 {
		return false, ""
	}

	// Parse count prefix
	count := 1
	idx := 0
	for idx < len(keys) && keys[idx] >= '1' && keys[idx] <= '9' {
		idx++
	}
	for idx < len(keys) && keys[idx] >= '0' && keys[idx] <= '9' {
		idx++
	}
	if idx > 0 {
		count = 0
		for i := 0; i < idx; i++ {
			count = count*10 + int(keys[i]-'0')
		}
		if count == 0 {
			count = 1
		}
	}
	keys = keys[idx:]

	if len(keys) == 0 {
		return false, strings.Repeat("0", idx) // Return the count as pending
	}

	// Handle commands
	switch {
	// Mode changes
	case keys == "i":
		e.buffer.SetMode(ModeInsert)
		return true, ""
	case keys == "I":
		MoveToFirstNonBlank(e.buffer)
		e.buffer.SetMode(ModeInsert)
		return true, ""
	case keys == "a":
		MoveRight(e.buffer, 1)
		e.buffer.SetMode(ModeInsert)
		return true, ""
	case keys == "A":
		MoveToLineEnd(e.buffer)
		e.buffer.cursorX++ // Move past last character
		e.buffer.SetMode(ModeInsert)
		return true, ""
	case keys == "o":
		e.saveUndo()
		MoveToLineEnd(e.buffer)
		e.buffer.Insert("\n")
		e.buffer.SetMode(ModeInsert)
		return true, ""
	case keys == "O":
		e.saveUndo()
		MoveToLineStart(e.buffer)
		e.buffer.Insert("\n")
		MoveUp(e.buffer, 1)
		e.buffer.SetMode(ModeInsert)
		return true, ""
	case keys == "v":
		e.buffer.SetMode(ModeVisual)
		return true, ""
	case keys == "V":
		e.buffer.SetMode(ModeVisualLine)
		return true, ""

	// Basic motions
	case keys == "h":
		MoveLeft(e.buffer, count)
		return true, ""
	case keys == "l":
		MoveRight(e.buffer, count)
		return true, ""
	case keys == "j":
		MoveDown(e.buffer, count)
		return true, ""
	case keys == "k":
		MoveUp(e.buffer, count)
		return true, ""
	case keys == "0":
		MoveToLineStart(e.buffer)
		return true, ""
	case keys == "$":
		MoveToLineEnd(e.buffer)
		return true, ""
	case keys == "^":
		MoveToFirstNonBlank(e.buffer)
		return true, ""
	case keys == "w":
		MoveWordForward(e.buffer, count)
		return true, ""
	case keys == "b":
		MoveWordBackward(e.buffer, count)
		return true, ""
	case keys == "e":
		MoveWordEnd(e.buffer, count)
		return true, ""
	case keys == "gg":
		MoveToBufferStart(e.buffer)
		return true, ""
	case keys == "G":
		if count > 1 {
			MoveToLine(e.buffer, count)
		} else {
			MoveToBufferEnd(e.buffer)
		}
		return true, ""
	case keys == "%":
		MoveToMatchingBracket(e.buffer)
		return true, ""

	// Find character
	case len(keys) >= 2 && keys[0] == 'f':
		char := rune(keys[1])
		MoveToChar(e.buffer, char, count, false)
		e.lastSearch = char
		e.searchDir = 1
		return true, keys[2:]
	case len(keys) >= 2 && keys[0] == 'F':
		char := rune(keys[1])
		MoveToCharBackward(e.buffer, char, count, false)
		e.lastSearch = char
		e.searchDir = -1
		return true, keys[2:]
	case len(keys) >= 2 && keys[0] == 't':
		char := rune(keys[1])
		MoveToChar(e.buffer, char, count, true)
		e.lastSearch = char
		e.searchDir = 1
		return true, keys[2:]
	case len(keys) >= 2 && keys[0] == 'T':
		char := rune(keys[1])
		MoveToCharBackward(e.buffer, char, count, true)
		e.lastSearch = char
		e.searchDir = -1
		return true, keys[2:]
	case keys == ";":
		if e.lastSearch != 0 {
			if e.searchDir > 0 {
				MoveToChar(e.buffer, e.lastSearch, count, false)
			} else {
				MoveToCharBackward(e.buffer, e.lastSearch, count, false)
			}
		}
		return true, ""
	case keys == ",":
		if e.lastSearch != 0 {
			if e.searchDir > 0 {
				MoveToCharBackward(e.buffer, e.lastSearch, count, false)
			} else {
				MoveToChar(e.buffer, e.lastSearch, count, false)
			}
		}
		return true, ""

	// Delete operations
	case keys == "x":
		e.saveUndo()
		deleted := e.buffer.Delete(count)
		e.buffer.SetRegister(deleted)
		return true, ""
	case keys == "X":
		e.saveUndo()
		for i := 0; i < count; i++ {
			if e.buffer.cursorX > 0 {
				MoveLeft(e.buffer, 1)
				e.buffer.Delete(1)
			}
		}
		return true, ""
	case keys == "dd":
		e.saveUndo()
		for i := 0; i < count; i++ {
			deleted := e.buffer.DeleteLine()
			e.buffer.SetRegister(deleted + "\n")
		}
		return true, ""
	case keys == "D":
		e.saveUndo()
		deleted := e.buffer.DeleteToEndOfLine()
		e.buffer.SetRegister(deleted)
		return true, ""
	case strings.HasPrefix(keys, "d"):
		return e.handleOperatorPending("d", keys[1:], count)

	// Change operations
	case keys == "cc" || keys == "S":
		e.saveUndo()
		e.buffer.DeleteLine()
		e.buffer.SetMode(ModeInsert)
		return true, ""
	case keys == "C":
		e.saveUndo()
		e.buffer.DeleteToEndOfLine()
		e.buffer.SetMode(ModeInsert)
		return true, ""
	case keys == "s":
		e.saveUndo()
		e.buffer.Delete(count)
		e.buffer.SetMode(ModeInsert)
		return true, ""
	case strings.HasPrefix(keys, "c"):
		return e.handleOperatorPending("c", keys[1:], count)

	// Replace
	case len(keys) >= 2 && keys[0] == 'r':
		e.saveUndo()
		char := rune(keys[1])
		for i := 0; i < count; i++ {
			e.buffer.ReplaceChar(char)
			if i < count-1 {
				MoveRight(e.buffer, 1)
			}
		}
		return true, keys[2:]

	// Yank operations
	case keys == "yy" || keys == "Y":
		line := e.buffer.CurrentLine()
		e.buffer.SetRegister(line + "\n")
		return true, ""
	case strings.HasPrefix(keys, "y"):
		return e.handleOperatorPending("y", keys[1:], count)

	// Put
	case keys == "p":
		e.saveUndo()
		reg := e.buffer.GetRegister()
		if strings.HasSuffix(reg, "\n") {
			// Line-wise paste
			MoveToLineEnd(e.buffer)
			e.buffer.cursorX++
			e.buffer.Insert("\n" + strings.TrimSuffix(reg, "\n"))
			MoveDown(e.buffer, 1)
			MoveToFirstNonBlank(e.buffer)
		} else {
			MoveRight(e.buffer, 1)
			e.buffer.Insert(reg)
		}
		return true, ""
	case keys == "P":
		e.saveUndo()
		reg := e.buffer.GetRegister()
		if strings.HasSuffix(reg, "\n") {
			// Line-wise paste above
			MoveToLineStart(e.buffer)
			e.buffer.Insert(strings.TrimSuffix(reg, "\n") + "\n")
			MoveUp(e.buffer, 1)
			MoveToFirstNonBlank(e.buffer)
		} else {
			e.buffer.Insert(reg)
		}
		return true, ""

	// Undo/Redo
	case keys == "u":
		e.undo()
		return true, ""
	case keys == "\x12": // Ctrl-R
		e.redo()
		return true, ""

	// Pending - wait for more input
	case keys == "g" || keys == "d" || keys == "c" || keys == "y" || keys == "f" || keys == "F" || keys == "t" || keys == "T" || keys == "r":
		return false, keys

	default:
		// Unknown command, discard
		return false, ""
	}
}

// handleOperatorPending handles operators like d, c, y with motions
func (e *Engine) handleOperatorPending(op, motion string, count int) (bool, string) {
	if len(motion) == 0 {
		return false, op // Still waiting for motion
	}

	startX, startY := e.buffer.CursorPosition()
	startIdx := e.buffer.CursorIndex()

	// Handle text objects
	if len(motion) >= 2 && (motion[0] == 'i' || motion[0] == 'a') {
		return e.handleTextObject(op, motion, count)
	}

	// Execute motion
	moved := false
	switch motion[0] {
	case 'w':
		moved = MoveWordForward(e.buffer, count)
		motion = motion[1:]
	case 'b':
		moved = MoveWordBackward(e.buffer, count)
		motion = motion[1:]
	case 'e':
		moved = MoveWordEnd(e.buffer, count)
		MoveRight(e.buffer, 1) // Include the character at end
		motion = motion[1:]
	case '$':
		moved = MoveToLineEnd(e.buffer)
		MoveRight(e.buffer, 1) // Include last char for d$
		motion = motion[1:]
	case '0':
		// d0 deletes from cursor to start of line
		endX := e.buffer.cursorX
		MoveToLineStart(e.buffer)
		startIdx = e.buffer.CursorIndex()
		e.buffer.cursorX = endX
		moved = true
		motion = motion[1:]
	case '^':
		endX := e.buffer.cursorX
		MoveToFirstNonBlank(e.buffer)
		if e.buffer.cursorX < endX {
			startIdx = e.buffer.CursorIndex()
			e.buffer.cursorX = endX
		}
		moved = true
		motion = motion[1:]
	case 'f', 't':
		if len(motion) >= 2 {
			char := rune(motion[1])
			if motion[0] == 'f' {
				moved = MoveToChar(e.buffer, char, count, false)
			} else {
				moved = MoveToChar(e.buffer, char, count, true)
			}
			MoveRight(e.buffer, 1) // Include target char
			motion = motion[2:]
		} else {
			return false, op + motion // Need more input
		}
	case 'F', 'T':
		if len(motion) >= 2 {
			char := rune(motion[1])
			if motion[0] == 'F' {
				moved = MoveToCharBackward(e.buffer, char, count, false)
			} else {
				moved = MoveToCharBackward(e.buffer, char, count, true)
			}
			// Swap start and end since we moved backward
			endIdx := e.buffer.CursorIndex()
			e.buffer.SetCursorIndex(startIdx)
			startIdx = endIdx
			motion = motion[2:]
		} else {
			return false, op + motion // Need more input
		}
	case 'G':
		if count > 1 {
			moved = MoveToLine(e.buffer, count)
		} else {
			moved = MoveToBufferEnd(e.buffer)
		}
		motion = motion[1:]
	case 'g':
		if len(motion) >= 2 && motion[1] == 'g' {
			moved = MoveToBufferStart(e.buffer)
			motion = motion[2:]
		} else {
			return false, op + motion
		}
	default:
		return false, ""
	}

	if !moved {
		e.buffer.SetCursorPosition(startX, startY)
		return true, motion
	}

	endIdx := e.buffer.CursorIndex()

	// Ensure startIdx < endIdx
	if startIdx > endIdx {
		startIdx, endIdx = endIdx, startIdx
	}

	// Perform operation
	e.saveUndo()
	e.buffer.SetCursorIndex(startIdx)

	switch op {
	case "d":
		deleted := e.buffer.Delete(endIdx - startIdx)
		e.buffer.SetRegister(deleted)
	case "c":
		deleted := e.buffer.Delete(endIdx - startIdx)
		e.buffer.SetRegister(deleted)
		e.buffer.SetMode(ModeInsert)
	case "y":
		// Yank without modifying buffer
		text := e.buffer.Text()
		if startIdx < len(text) && endIdx <= len(text) {
			e.buffer.SetRegister(text[startIdx:endIdx])
		}
		e.buffer.SetCursorPosition(startX, startY) // Return to original position
	}

	return true, motion
}

// handleTextObject handles inner and around text objects
func (e *Engine) handleTextObject(op, motion string, count int) (bool, string) {
	if len(motion) < 2 {
		return false, op + motion
	}

	inner := motion[0] == 'i'
	obj := motion[1]
	motion = motion[2:]

	startX, startY := e.buffer.CursorPosition()
	line := e.buffer.CurrentLine()
	runes := []rune(line)

	var startIdx, endIdx int

	switch obj {
	case 'w': // word
		// Find word boundaries
		x := e.buffer.cursorX
		if x >= len(runes) {
			return true, motion
		}

		// Find start of word
		start := x
		class := classifyChar(runes[x])
		for start > 0 && classifyChar(runes[start-1]) == class {
			start--
		}

		// Find end of word
		end := x
		for end < len(runes)-1 && classifyChar(runes[end+1]) == class {
			end++
		}

		if !inner {
			// Include trailing whitespace for 'aw'
			for end < len(runes)-1 && unicode.IsSpace(runes[end+1]) {
				end++
			}
		}

		e.buffer.cursorX = start
		startIdx = e.buffer.CursorIndex()
		e.buffer.cursorX = end + 1
		endIdx = e.buffer.CursorIndex()

	case '"', '\'', '`': // quotes
		quote := rune(obj)
		// Find opening quote
		openIdx := -1
		closeIdx := -1
		for i, r := range runes {
			if r == quote {
				if openIdx == -1 {
					openIdx = i
				} else {
					closeIdx = i
					break
				}
			}
		}
		if openIdx == -1 || closeIdx == -1 {
			return true, motion
		}
		if inner {
			e.buffer.cursorX = openIdx + 1
			startIdx = e.buffer.CursorIndex()
			e.buffer.cursorX = closeIdx
			endIdx = e.buffer.CursorIndex()
		} else {
			e.buffer.cursorX = openIdx
			startIdx = e.buffer.CursorIndex()
			e.buffer.cursorX = closeIdx + 1
			endIdx = e.buffer.CursorIndex()
		}

	case '(', ')', 'b': // parentheses
		return e.handleBracketObject(op, '(', ')', inner, motion)
	case '[', ']': // brackets
		return e.handleBracketObject(op, '[', ']', inner, motion)
	case '{', '}', 'B': // braces
		return e.handleBracketObject(op, '{', '}', inner, motion)
	case '<', '>': // angle brackets
		return e.handleBracketObject(op, '<', '>', inner, motion)

	default:
		return true, motion
	}

	if startIdx >= endIdx {
		e.buffer.SetCursorPosition(startX, startY)
		return true, motion
	}

	// Perform operation
	e.saveUndo()
	e.buffer.SetCursorIndex(startIdx)

	switch op {
	case "d":
		deleted := e.buffer.Delete(endIdx - startIdx)
		e.buffer.SetRegister(deleted)
	case "c":
		deleted := e.buffer.Delete(endIdx - startIdx)
		e.buffer.SetRegister(deleted)
		e.buffer.SetMode(ModeInsert)
	case "y":
		text := e.buffer.Text()
		if startIdx < len(text) && endIdx <= len(text) {
			e.buffer.SetRegister(text[startIdx:endIdx])
		}
		e.buffer.SetCursorPosition(startX, startY)
	}

	return true, motion
}

// handleBracketObject handles bracket text objects
func (e *Engine) handleBracketObject(op string, open, close rune, inner bool, remaining string) (bool, string) {
	text := e.buffer.Text()
	cursorIdx := e.buffer.CursorIndex()
	runes := []rune(text)

	// Find opening bracket (searching backward from cursor)
	openIdx := -1
	depth := 0
	for i := cursorIdx; i >= 0; i-- {
		if runes[i] == close {
			depth++
		} else if runes[i] == open {
			if depth == 0 {
				openIdx = i
				break
			}
			depth--
		}
	}

	if openIdx == -1 {
		return true, remaining
	}

	// Find closing bracket
	closeIdx := -1
	depth = 1
	for i := openIdx + 1; i < len(runes); i++ {
		if runes[i] == open {
			depth++
		} else if runes[i] == close {
			depth--
			if depth == 0 {
				closeIdx = i
				break
			}
		}
	}

	if closeIdx == -1 {
		return true, remaining
	}

	var startIdx, endIdx int
	if inner {
		startIdx = openIdx + 1
		endIdx = closeIdx
	} else {
		startIdx = openIdx
		endIdx = closeIdx + 1
	}

	if startIdx >= endIdx {
		return true, remaining
	}

	e.saveUndo()
	e.buffer.SetCursorIndex(startIdx)

	switch op {
	case "d":
		deleted := e.buffer.Delete(endIdx - startIdx)
		e.buffer.SetRegister(deleted)
	case "c":
		deleted := e.buffer.Delete(endIdx - startIdx)
		e.buffer.SetRegister(deleted)
		e.buffer.SetMode(ModeInsert)
	case "y":
		e.buffer.SetRegister(string(runes[startIdx:endIdx]))
	}

	return true, remaining
}

// handleVisualMode handles keys in visual mode
func (e *Engine) handleVisualMode(keys string) (bool, string) {
	if len(keys) == 0 {
		return false, ""
	}

	switch keys {
	case "esc", "\x1b", "v", "V":
		e.buffer.SetMode(ModeNormal)
		return true, ""
	case "d", "x":
		// Delete visual selection (simplified - just exits visual mode for now)
		e.buffer.SetMode(ModeNormal)
		return true, ""
	case "y":
		// Yank visual selection
		e.buffer.SetMode(ModeNormal)
		return true, ""
	// Add motion support in visual mode
	case "h":
		MoveLeft(e.buffer, 1)
		return true, ""
	case "l":
		MoveRight(e.buffer, 1)
		return true, ""
	case "j":
		MoveDown(e.buffer, 1)
		return true, ""
	case "k":
		MoveUp(e.buffer, 1)
		return true, ""
	case "w":
		MoveWordForward(e.buffer, 1)
		return true, ""
	case "b":
		MoveWordBackward(e.buffer, 1)
		return true, ""
	case "e":
		MoveWordEnd(e.buffer, 1)
		return true, ""
	case "$":
		MoveToLineEnd(e.buffer)
		return true, ""
	case "0":
		MoveToLineStart(e.buffer)
		return true, ""
	default:
		return false, ""
	}
}

// saveUndo saves current state for undo
func (e *Engine) saveUndo() {
	e.undoStack = append(e.undoStack, e.buffer.Clone())
	e.redoStack = nil // Clear redo stack on new change
	// Limit undo stack size
	if len(e.undoStack) > 100 {
		e.undoStack = e.undoStack[1:]
	}
}

// undo reverts to previous state
func (e *Engine) undo() bool {
	if len(e.undoStack) == 0 {
		return false
	}
	// Save current state to redo
	e.redoStack = append(e.redoStack, e.buffer.Clone())
	// Pop from undo stack
	e.buffer = e.undoStack[len(e.undoStack)-1]
	e.undoStack = e.undoStack[:len(e.undoStack)-1]
	return true
}

// redo reapplies undone changes
func (e *Engine) redo() bool {
	if len(e.redoStack) == 0 {
		return false
	}
	// Save current state to undo
	e.undoStack = append(e.undoStack, e.buffer.Clone())
	// Pop from redo stack
	e.buffer = e.redoStack[len(e.redoStack)-1]
	e.redoStack = e.redoStack[:len(e.redoStack)-1]
	return true
}

// Reset resets the engine with new text
func (e *Engine) Reset(text string, cursorPos int) {
	e.buffer = NewBuffer(text)
	e.buffer.SetCursorIndex(cursorPos)
	e.undoStack = nil
	e.redoStack = nil
	e.pendingKeys = ""
}

// GetPendingKeys returns any pending key sequence
func (e *Engine) GetPendingKeys() string {
	return e.pendingKeys
}

// LineCount returns the number of lines in the buffer
func (e *Engine) LineCount() int {
	return len(e.buffer.lines)
}

// GetLine returns a specific line
func (e *Engine) GetLine(n int) string {
	if n >= 0 && n < len(e.buffer.lines) {
		return e.buffer.lines[n]
	}
	return ""
}
