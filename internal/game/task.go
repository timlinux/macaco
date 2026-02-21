package game

import (
	"encoding/json"
	"os"
)

// TaskCategory represents the type of vim operation
type TaskCategory string

const (
	CategoryMotion  TaskCategory = "motion"
	CategoryDelete  TaskCategory = "delete"
	CategoryChange  TaskCategory = "change"
	CategoryInsert  TaskCategory = "insert"
	CategoryVisual  TaskCategory = "visual"
	CategoryComplex TaskCategory = "complex"
)

// Task represents a vim training task
type Task struct {
	ID             string       `json:"id"`
	Category       TaskCategory `json:"category"`
	Difficulty     int          `json:"difficulty"`
	Initial        string       `json:"initial"`
	Desired        string       `json:"desired"`
	CursorStart    int          `json:"cursor_start"`
	CursorEnd      int          `json:"cursor_end,omitempty"`       // For motion tasks
	HighlightStart int          `json:"highlight_start,omitempty"` // Start of text to modify
	HighlightEnd   int          `json:"highlight_end,omitempty"`   // End of text to modify
	OptimalKeys    string       `json:"optimal_keys"`
	OptimalCount   int          `json:"optimal_count"`
	Description    string       `json:"description"`
	Hint           string       `json:"hint"`
	Tags           []string     `json:"tags,omitempty"`
}

// IsMotionTask returns true if this is a motion-only task
func (t *Task) IsMotionTask() bool {
	return t.Initial == t.Desired && t.CursorEnd > 0
}

// HasHighlight returns true if the task has text to highlight
func (t *Task) HasHighlight() bool {
	return t.HighlightEnd > t.HighlightStart
}

// GetHighlightedText returns the text that should be highlighted
func (t *Task) GetHighlightedText() string {
	if !t.HasHighlight() {
		return ""
	}
	runes := []rune(t.Initial)
	if t.HighlightStart >= len(runes) || t.HighlightEnd > len(runes) {
		return ""
	}
	return string(runes[t.HighlightStart:t.HighlightEnd])
}

// TaskDatabase holds all available tasks
type TaskDatabase struct {
	Version     string                    `json:"version"`
	LastUpdated string                    `json:"last_updated"`
	Rounds      map[string]RoundDef       `json:"rounds"`
	Tasks       []Task                    `json:"tasks"`
	tasksByID   map[string]*Task          // Lookup cache
	tasksByCat  map[TaskCategory][]*Task  // Category lookup
}

// RoundDef defines a round type
type RoundDef struct {
	Name             string         `json:"name"`
	DifficultyRange  [2]int         `json:"difficulty_range"`
	TaskDistribution map[string]int `json:"task_distribution"`
	TaskIDs          []string       `json:"tasks,omitempty"`
}

// LoadTaskDatabase loads tasks from a JSON file
func LoadTaskDatabase(path string) (*TaskDatabase, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var db TaskDatabase
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, err
	}

	db.buildLookups()
	return &db, nil
}

// NewEmbeddedTaskDatabase creates a task database with built-in tasks
// Deprecated: Use NewGeneratedTaskDatabase for procedurally generated tasks
func NewEmbeddedTaskDatabase() *TaskDatabase {
	db := &TaskDatabase{
		Version:     "1.0.0",
		LastUpdated: "2026-02-21",
		Rounds:      make(map[string]RoundDef),
		Tasks:       getEmbeddedTasks(),
	}
	db.buildLookups()
	db.buildRounds()
	return db
}

// NewGeneratedTaskDatabase creates a task database with procedurally generated tasks
// using public domain texts from classic literature
func NewGeneratedTaskDatabase() *TaskDatabase {
	generator := NewTaskGenerator()

	// Generate tasks for all round types to populate the database
	var allTasks []Task
	roundTypes := []string{"beginner", "intermediate", "advanced", "expert"}

	for _, roundType := range roundTypes {
		tasks := generator.GenerateTasksForRound(roundType)
		allTasks = append(allTasks, tasks...)
	}

	db := &TaskDatabase{
		Version:     "2.0.0",
		LastUpdated: "2026-02-21",
		Rounds:      make(map[string]RoundDef),
		Tasks:       allTasks,
	}
	db.buildLookups()
	db.buildRounds()
	return db
}

// buildLookups builds internal lookup maps
func (db *TaskDatabase) buildLookups() {
	db.tasksByID = make(map[string]*Task)
	db.tasksByCat = make(map[TaskCategory][]*Task)

	for i := range db.Tasks {
		task := &db.Tasks[i]
		db.tasksByID[task.ID] = task
		db.tasksByCat[task.Category] = append(db.tasksByCat[task.Category], task)
	}
}

// buildRounds builds default round definitions
func (db *TaskDatabase) buildRounds() {
	db.Rounds = map[string]RoundDef{
		"beginner": {
			Name:            "Beginner Round",
			DifficultyRange: [2]int{1, 1},
			TaskDistribution: map[string]int{
				"motion":  6,
				"delete":  6,
				"change":  6,
				"insert":  6,
				"visual":  3,
				"complex": 3,
			},
		},
		"intermediate": {
			Name:            "Intermediate Round",
			DifficultyRange: [2]int{1, 2},
			TaskDistribution: map[string]int{
				"motion":  6,
				"delete":  6,
				"change":  6,
				"insert":  6,
				"visual":  3,
				"complex": 3,
			},
		},
		"advanced": {
			Name:            "Advanced Round",
			DifficultyRange: [2]int{2, 3},
			TaskDistribution: map[string]int{
				"motion":  6,
				"delete":  6,
				"change":  6,
				"insert":  6,
				"visual":  3,
				"complex": 3,
			},
		},
		"expert": {
			Name:            "Expert Round",
			DifficultyRange: [2]int{3, 4},
			TaskDistribution: map[string]int{
				"motion":  6,
				"delete":  6,
				"change":  6,
				"insert":  6,
				"visual":  3,
				"complex": 3,
			},
		},
		"mixed": {
			Name:            "Mixed Round",
			DifficultyRange: [2]int{1, 4},
			TaskDistribution: map[string]int{
				"motion":  6,
				"delete":  6,
				"change":  6,
				"insert":  6,
				"visual":  3,
				"complex": 3,
			},
		},
	}
}

// GetTask returns a task by ID
func (db *TaskDatabase) GetTask(id string) *Task {
	return db.tasksByID[id]
}

// GetTasksByCategory returns all tasks in a category
func (db *TaskDatabase) GetTasksByCategory(cat TaskCategory) []*Task {
	return db.tasksByCat[cat]
}

// GetTasksForRound returns tasks for a specific round type
func (db *TaskDatabase) GetTasksForRound(roundType string) []*Task {
	roundDef, ok := db.Rounds[roundType]
	if !ok {
		roundDef = db.Rounds["beginner"]
	}

	var tasks []*Task
	minDiff, maxDiff := roundDef.DifficultyRange[0], roundDef.DifficultyRange[1]

	categories := []TaskCategory{
		CategoryMotion, CategoryDelete, CategoryChange,
		CategoryInsert, CategoryVisual, CategoryComplex,
	}

	for _, cat := range categories {
		count := roundDef.TaskDistribution[string(cat)]
		catTasks := db.GetTasksByCategory(cat)

		// Filter by difficulty
		var eligible []*Task
		for _, t := range catTasks {
			if t.Difficulty >= minDiff && t.Difficulty <= maxDiff {
				eligible = append(eligible, t)
			}
		}

		// Take up to count tasks
		for i := 0; i < count && i < len(eligible); i++ {
			tasks = append(tasks, eligible[i])
		}
	}

	return tasks
}

// getEmbeddedTasks returns the built-in task collection
func getEmbeddedTasks() []Task {
	return []Task{
		// Motion tasks - Level 1
		{
			ID: "motion-w-001", Category: CategoryMotion, Difficulty: 1,
			Initial: "hello world from vim", Desired: "hello world from vim",
			CursorStart: 0, CursorEnd: 6,
			OptimalKeys: "w", OptimalCount: 1,
			Description: "Move to next word start",
			Hint:        "Use 'w' to move forward one word",
			Tags:        []string{"word-motion", "basic"},
		},
		{
			ID: "motion-b-001", Category: CategoryMotion, Difficulty: 1,
			Initial: "hello world from vim", Desired: "hello world from vim",
			CursorStart: 12, CursorEnd: 6,
			OptimalKeys: "b", OptimalCount: 1,
			Description: "Move to previous word start",
			Hint:        "Use 'b' to move backward one word",
			Tags:        []string{"word-motion", "basic"},
		},
		{
			ID: "motion-e-001", Category: CategoryMotion, Difficulty: 1,
			Initial: "hello world from vim", Desired: "hello world from vim",
			CursorStart: 0, CursorEnd: 4,
			OptimalKeys: "e", OptimalCount: 1,
			Description: "Move to end of word",
			Hint:        "Use 'e' to move to end of current word",
			Tags:        []string{"word-motion", "basic"},
		},
		{
			ID: "motion-0-001", Category: CategoryMotion, Difficulty: 1,
			Initial: "hello world from vim", Desired: "hello world from vim",
			CursorStart: 12, CursorEnd: 0,
			OptimalKeys: "0", OptimalCount: 1,
			Description: "Move to start of line",
			Hint:        "Use '0' to jump to line start",
			Tags:        []string{"line-motion", "basic"},
		},
		{
			ID: "motion-$-001", Category: CategoryMotion, Difficulty: 1,
			Initial: "hello world from vim", Desired: "hello world from vim",
			CursorStart: 0, CursorEnd: 19,
			OptimalKeys: "$", OptimalCount: 1,
			Description: "Move to end of line",
			Hint:        "Use '$' to jump to line end",
			Tags:        []string{"line-motion", "basic"},
		},
		{
			ID: "motion-^-001", Category: CategoryMotion, Difficulty: 1,
			Initial: "   hello world", Desired: "   hello world",
			CursorStart: 0, CursorEnd: 3,
			OptimalKeys: "^", OptimalCount: 1,
			Description: "Move to first non-blank",
			Hint:        "Use '^' to jump to first non-blank character",
			Tags:        []string{"line-motion", "basic"},
		},

		// Motion tasks - Level 2
		{
			ID: "motion-3w-001", Category: CategoryMotion, Difficulty: 2,
			Initial: "one two three four five", Desired: "one two three four five",
			CursorStart: 0, CursorEnd: 14,
			OptimalKeys: "3w", OptimalCount: 2,
			Description: "Move forward three words",
			Hint:        "Use count before motion: '3w'",
			Tags:        []string{"word-motion", "count"},
		},
		{
			ID: "motion-f-001", Category: CategoryMotion, Difficulty: 2,
			Initial: "find the letter x here", Desired: "find the letter x here",
			CursorStart: 0, CursorEnd: 16,
			OptimalKeys: "fx", OptimalCount: 2,
			Description: "Find next 'x'",
			Hint:        "Use 'f' followed by target character",
			Tags:        []string{"find-motion", "intermediate"},
		},
		{
			ID: "motion-t-001", Category: CategoryMotion, Difficulty: 2,
			Initial: "move to (parenthesis)", Desired: "move to (parenthesis)",
			CursorStart: 0, CursorEnd: 7,
			OptimalKeys: "t(", OptimalCount: 2,
			Description: "Move until '('",
			Hint:        "Use 't' to move until (before) a character",
			Tags:        []string{"find-motion", "intermediate"},
		},
		{
			ID: "motion-2b-001", Category: CategoryMotion, Difficulty: 2,
			Initial: "one two three four", Desired: "one two three four",
			CursorStart: 14, CursorEnd: 4,
			OptimalKeys: "2b", OptimalCount: 2,
			Description: "Move back two words",
			Hint:        "Use '2b' to move back two words",
			Tags:        []string{"word-motion", "count"},
		},

		// Delete tasks - Level 1
		{
			ID: "delete-x-001", Category: CategoryDelete, Difficulty: 1,
			Initial: "helxlo world", Desired: "hello world",
			CursorStart: 3, OptimalKeys: "x", OptimalCount: 1,
			Description: "Delete character under cursor",
			Hint:        "Use 'x' to delete character under cursor",
			Tags:        []string{"delete", "basic"},
		},
		{
			ID: "delete-x-002", Category: CategoryDelete, Difficulty: 1,
			Initial: "hellooo world", Desired: "hello world",
			CursorStart: 5, OptimalKeys: "xx", OptimalCount: 2,
			Description: "Delete two extra characters",
			Hint:        "Use 'x' twice or '2x' to delete two characters",
			Tags:        []string{"delete", "basic"},
		},
		{
			ID: "delete-dw-001", Category: CategoryDelete, Difficulty: 1,
			Initial: "hello extra world", Desired: "hello world",
			CursorStart: 6, OptimalKeys: "dw", OptimalCount: 2,
			Description: "Delete word",
			Hint:        "Use 'dw' to delete word",
			Tags:        []string{"delete", "word"},
		},
		{
			ID: "delete-dd-001", Category: CategoryDelete, Difficulty: 1,
			Initial: "line one\ndelete me\nline three", Desired: "line one\nline three",
			CursorStart: 9, OptimalKeys: "dd", OptimalCount: 2,
			Description: "Delete entire line",
			Hint:        "Use 'dd' to delete the current line",
			Tags:        []string{"delete", "line"},
		},
		{
			ID: "delete-D-001", Category: CategoryDelete, Difficulty: 1,
			Initial: "keep this delete rest", Desired: "keep this",
			CursorStart: 9, OptimalKeys: "D", OptimalCount: 1,
			Description: "Delete to end of line",
			Hint:        "Use 'D' to delete from cursor to end of line",
			Tags:        []string{"delete", "line"},
		},

		// Delete tasks - Level 2
		{
			ID: "delete-daw-001", Category: CategoryDelete, Difficulty: 2,
			Initial: "delete this word here", Desired: "delete word here",
			CursorStart: 9, OptimalKeys: "daw", OptimalCount: 3,
			Description: "Delete a word (including space)",
			Hint:        "Use 'daw' to delete 'a word' including surrounding space",
			Tags:        []string{"delete", "text-object"},
		},
		{
			ID: "delete-diw-001", Category: CategoryDelete, Difficulty: 2,
			Initial: "delete this word", Desired: "delete  word",
			CursorStart: 8, OptimalKeys: "diw", OptimalCount: 3,
			Description: "Delete inner word",
			Hint:        "Use 'diw' to delete just the word (not space)",
			Tags:        []string{"delete", "text-object"},
		},
		{
			ID: "delete-dt-001", Category: CategoryDelete, Difficulty: 2,
			Initial: "delete until (keep this)", Desired: "(keep this)",
			CursorStart: 0, OptimalKeys: "dt(", OptimalCount: 3,
			Description: "Delete until character",
			Hint:        "Use 'dt(' to delete until the '(' character",
			Tags:        []string{"delete", "find-motion"},
		},

		// Change tasks - Level 1
		{
			ID: "change-cw-001", Category: CategoryChange, Difficulty: 1,
			Initial: "hello old world", Desired: "hello new world",
			CursorStart: 6, OptimalKeys: "cwnew<ESC>", OptimalCount: 6,
			Description: "Change word to 'new'",
			Hint:        "Use 'cw' to change word, type 'new', press ESC",
			Tags:        []string{"change", "word"},
		},
		{
			ID: "change-r-001", Category: CategoryChange, Difficulty: 1,
			Initial: "hello warld", Desired: "hello world",
			CursorStart: 7, OptimalKeys: "ro", OptimalCount: 2,
			Description: "Replace single character",
			Hint:        "Use 'r' followed by the replacement character",
			Tags:        []string{"change", "replace"},
		},
		{
			ID: "change-s-001", Category: CategoryChange, Difficulty: 1,
			Initial: "hello xorld", Desired: "hello world",
			CursorStart: 6, OptimalKeys: "sw<ESC>", OptimalCount: 4,
			Description: "Substitute character",
			Hint:        "Use 's' to delete char and enter insert mode",
			Tags:        []string{"change", "substitute"},
		},
		{
			ID: "change-cc-001", Category: CategoryChange, Difficulty: 1,
			Initial: "wrong line", Desired: "right line",
			CursorStart: 0, OptimalKeys: "ccright line<ESC>", OptimalCount: 13,
			Description: "Change entire line",
			Hint:        "Use 'cc' to change the entire line",
			Tags:        []string{"change", "line"},
		},
		{
			ID: "change-C-001", Category: CategoryChange, Difficulty: 1,
			Initial: "keep this wrong part", Desired: "keep this right part",
			CursorStart: 10, OptimalKeys: "Cright part<ESC>", OptimalCount: 12,
			Description: "Change to end of line",
			Hint:        "Use 'C' to change from cursor to end of line",
			Tags:        []string{"change", "line"},
		},

		// Change tasks - Level 2
		{
			ID: "change-ciw-001", Category: CategoryChange, Difficulty: 2,
			Initial: "change inside word", Desired: "change outside word",
			CursorStart: 10, OptimalKeys: "ciwoutside<ESC>", OptimalCount: 11,
			Description: "Change inner word",
			Hint:        "Use 'ciw' to change word regardless of cursor position",
			Tags:        []string{"change", "text-object"},
		},
		{
			ID: "change-ci-quote-001", Category: CategoryChange, Difficulty: 2,
			Initial: `text = "old value"`, Desired: `text = "new value"`,
			CursorStart: 10, OptimalKeys: `ci"new value<ESC>`, OptimalCount: 14,
			Description: "Change inside quotes",
			Hint:        `Use 'ci"' to change text inside quotes`,
			Tags:        []string{"change", "text-object"},
		},
		{
			ID: "change-ci-paren-001", Category: CategoryChange, Difficulty: 2,
			Initial: "func(old, args)", Desired: "func(new, args)",
			CursorStart: 6, OptimalKeys: "cwnew<ESC>", OptimalCount: 6,
			Description: "Change function argument",
			Hint:        "Use 'cw' to change the word",
			Tags:        []string{"change", "word"},
		},

		// Insert tasks - Level 1
		{
			ID: "insert-i-001", Category: CategoryInsert, Difficulty: 1,
			Initial: "hello world", Desired: "hello beautiful world",
			CursorStart: 6, OptimalKeys: "ibeautiful <ESC>", OptimalCount: 12,
			Description: "Insert text before cursor",
			Hint:        "Use 'i' to insert before cursor",
			Tags:        []string{"insert", "basic"},
		},
		{
			ID: "insert-a-001", Category: CategoryInsert, Difficulty: 1,
			Initial: "hello world", Desired: "hello, world",
			CursorStart: 4, OptimalKeys: "a,<ESC>", OptimalCount: 4,
			Description: "Append after cursor",
			Hint:        "Use 'a' to append after cursor",
			Tags:        []string{"insert", "basic"},
		},
		{
			ID: "insert-A-001", Category: CategoryInsert, Difficulty: 1,
			Initial: "hello world", Desired: "hello world!",
			CursorStart: 0, OptimalKeys: "A!<ESC>", OptimalCount: 4,
			Description: "Append at end of line",
			Hint:        "Use 'A' to append at end of line",
			Tags:        []string{"insert", "basic"},
		},
		{
			ID: "insert-I-001", Category: CategoryInsert, Difficulty: 1,
			Initial: "world", Desired: "hello world",
			CursorStart: 2, OptimalKeys: "Ihello <ESC>", OptimalCount: 9,
			Description: "Insert at start of line",
			Hint:        "Use 'I' to insert at the beginning of line",
			Tags:        []string{"insert", "basic"},
		},
		{
			ID: "insert-o-001", Category: CategoryInsert, Difficulty: 1,
			Initial: "line one\nline three", Desired: "line one\nline two\nline three",
			CursorStart: 0, OptimalKeys: "oline two<ESC>", OptimalCount: 11,
			Description: "Open new line below",
			Hint:        "Use 'o' to open new line below and enter insert mode",
			Tags:        []string{"insert", "line"},
		},
		{
			ID: "insert-O-001", Category: CategoryInsert, Difficulty: 1,
			Initial: "line two\nline three", Desired: "line one\nline two\nline three",
			CursorStart: 0, OptimalKeys: "Oline one<ESC>", OptimalCount: 11,
			Description: "Open new line above",
			Hint:        "Use 'O' to open new line above and enter insert mode",
			Tags:        []string{"insert", "line"},
		},

		// Visual tasks - Level 2
		{
			ID: "visual-vwd-001", Category: CategoryVisual, Difficulty: 2,
			Initial: "select extra word here", Desired: "select word here",
			CursorStart: 7, OptimalKeys: "vewd", OptimalCount: 4,
			Description: "Visual select and delete word",
			Hint:        "Use 'v' to enter visual, 'e' to extend, 'd' to delete",
			Tags:        []string{"visual", "delete"},
		},
		{
			ID: "visual-Vy-001", Category: CategoryVisual, Difficulty: 2,
			Initial: "copy this line\npaste here", Desired: "copy this line\ncopy this line\npaste here",
			CursorStart: 0, OptimalKeys: "Vyjp", OptimalCount: 4,
			Description: "Visual line yank and paste",
			Hint:        "Use 'V' for line visual, 'y' yank, 'j' down, 'p' paste",
			Tags:        []string{"visual", "yank", "paste"},
		},
		{
			ID: "visual-viw-001", Category: CategoryVisual, Difficulty: 2,
			Initial: "change this word now", Desired: "change that word now",
			CursorStart: 9, OptimalKeys: "ciwcthat<ESC>", OptimalCount: 9,
			Description: "Change inner word",
			Hint:        "Use 'ciw' to change the word under cursor",
			Tags:        []string{"visual", "change"},
		},

		// Complex tasks - Level 3
		{
			ID: "complex-dt-001", Category: CategoryComplex, Difficulty: 3,
			Initial: "delete until (keep this)", Desired: "(keep this)",
			CursorStart: 0, OptimalKeys: "dt(", OptimalCount: 3,
			Description: "Delete until character",
			Hint:        "Use 'dt(' to delete until '('",
			Tags:        []string{"complex", "delete", "find"},
		},
		{
			ID: "complex-cf-001", Category: CategoryComplex, Difficulty: 3,
			Initial: "change until,comma here", Desired: "new text,comma here",
			CursorStart: 0, OptimalKeys: "cf,new text<ESC>", OptimalCount: 12,
			Description: "Change through character",
			Hint:        "Use 'cf,' to change through comma",
			Tags:        []string{"complex", "change", "find"},
		},
		{
			ID: "complex-yap-001", Category: CategoryComplex, Difficulty: 3,
			Initial: "func(old)", Desired: "func(old)old",
			CursorStart: 5, OptimalKeys: "yi)$p", OptimalCount: 5,
			Description: "Yank inside and paste",
			Hint:        "Use 'yi)' to yank inside parens, '$p' to paste at end",
			Tags:        []string{"complex", "yank", "paste"},
		},

		// Complex tasks - Level 4
		{
			ID: "complex-multi-001", Category: CategoryComplex, Difficulty: 4,
			Initial: "const oldName = 'value';", Desired: "const newName = 'value';",
			CursorStart: 6, OptimalKeys: "cwnewName<ESC>", OptimalCount: 11,
			Description: "Change variable name",
			Hint:        "Navigate to word and use 'cw'",
			Tags:        []string{"complex", "change"},
		},
		{
			ID: "complex-swap-001", Category: CategoryComplex, Difficulty: 4,
			Initial: "second first", Desired: "first second",
			CursorStart: 0, OptimalKeys: "dwwP", OptimalCount: 4,
			Description: "Swap two words",
			Hint:        "Delete word, move, paste before",
			Tags:        []string{"complex", "delete", "paste"},
		},
		{
			ID: "complex-dup-001", Category: CategoryComplex, Difficulty: 4,
			Initial: "duplicate", Desired: "duplicate duplicate",
			CursorStart: 0, OptimalKeys: "yiwA <ESC>p", OptimalCount: 8,
			Description: "Duplicate word at end",
			Hint:        "Yank word, append space at end, paste",
			Tags:        []string{"complex", "yank", "paste"},
		},
	}
}
