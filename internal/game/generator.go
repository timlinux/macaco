package game

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

// TextSource represents a source of public domain text
type TextSource struct {
	Name        string
	Author      string
	Year        int
	License     string
	Attribution string
	Sentences   []string
}

// GetPublicDomainSources returns text sources from public domain works
// All texts are from works published before 1928 and are in the public domain
func GetPublicDomainSources() []TextSource {
	return []TextSource{
		{
			Name:        "Pride and Prejudice",
			Author:      "Jane Austen",
			Year:        1813,
			License:     "Public Domain",
			Attribution: "Text from 'Pride and Prejudice' by Jane Austen (1813), sourced from Project Gutenberg",
			Sentences: []string{
				"It is a truth universally acknowledged",
				"that a single man in possession of a good fortune",
				"must be in want of a wife",
				"Mr Bennet was so odd a mixture",
				"of quick parts and sarcastic humour",
				"reserve and caprice",
				"She was a woman of mean understanding",
				"little information and uncertain temper",
				"The business of her life was to get her daughters married",
				"He was an intelligent and handsome man",
				"Elizabeth had been obliged to accept him",
				"Their visit afforded was the sight",
				"Mr Darcy soon drew the attention",
				"of the room by his fine person",
				"He was the proudest man in the world",
				"She is tolerable but not handsome enough",
				"I could easily forgive his pride",
				"The evening altogether passed off pleasantly",
				"I have been used to consider poetry",
				"as the food of love",
			},
		},
		{
			Name:        "A Tale of Two Cities",
			Author:      "Charles Dickens",
			Year:        1859,
			License:     "Public Domain",
			Attribution: "Text from 'A Tale of Two Cities' by Charles Dickens (1859), sourced from Project Gutenberg",
			Sentences: []string{
				"It was the best of times",
				"it was the worst of times",
				"it was the age of wisdom",
				"it was the age of foolishness",
				"it was the epoch of belief",
				"it was the epoch of incredulity",
				"it was the season of Light",
				"it was the season of Darkness",
				"it was the spring of hope",
				"it was the winter of despair",
				"we had everything before us",
				"we had nothing before us",
				"we were all going direct to Heaven",
				"we were all going direct the other way",
				"There were a king with a large jaw",
				"and a queen with a plain face",
				"In both countries it was clearer than crystal",
				"the state of public feeling",
				"France received the news by mail",
				"England had a certain authority",
			},
		},
		{
			Name:        "The Adventures of Sherlock Holmes",
			Author:      "Arthur Conan Doyle",
			Year:        1892,
			License:     "Public Domain",
			Attribution: "Text from 'The Adventures of Sherlock Holmes' by Arthur Conan Doyle (1892), sourced from Project Gutenberg",
			Sentences: []string{
				"To Sherlock Holmes she is always the woman",
				"I have seldom heard him mention her",
				"In his eyes she eclipses the whole of her sex",
				"He never spoke of softer passions",
				"They were admirable things for the observer",
				"But for the trained reasoner",
				"to admit such intrusions into his own mind",
				"was to introduce a distracting factor",
				"I had seen little of Holmes lately",
				"My marriage had drifted us away",
				"My own complete happiness",
				"absorbed all my attention",
				"He was buried in his chair",
				"reading and rereading a letter",
				"The note was undated",
				"and without signature or address",
				"There will call upon you tonight",
				"a gentleman who desires to consult",
				"Your recent services to the crown",
				"have shown that you may be trusted",
			},
		},
		{
			Name:        "Moby Dick",
			Author:      "Herman Melville",
			Year:        1851,
			License:     "Public Domain",
			Attribution: "Text from 'Moby Dick' by Herman Melville (1851), sourced from Project Gutenberg",
			Sentences: []string{
				"Call me Ishmael",
				"Some years ago never mind how long",
				"having little or no money in my purse",
				"and nothing particular to interest me",
				"I thought I would sail about a little",
				"and see the watery part of the world",
				"whenever I find myself growing grim",
				"whenever it is a damp drizzly November",
				"I account it high time to get to sea",
				"This is my substitute for pistol and ball",
				"There now is your insular city",
				"belted round by wharves as Indian isles",
				"Commerce surrounds it with her surf",
				"Right and left the streets take you",
				"Its extreme downtown is the battery",
				"where that noble mole is washed",
				"Look at the crowds of water gazers",
				"Circumambulate the city on a dreamy day",
				"Go from Corlears Hook to Coenties Slip",
				"What do you see there",
			},
		},
		{
			Name:        "Alice's Adventures in Wonderland",
			Author:      "Lewis Carroll",
			Year:        1865,
			License:     "Public Domain",
			Attribution: "Text from 'Alice's Adventures in Wonderland' by Lewis Carroll (1865), sourced from Project Gutenberg",
			Sentences: []string{
				"Alice was beginning to get very tired",
				"of sitting by her sister on the bank",
				"and of having nothing to do",
				"once or twice she had peeped",
				"into the book her sister was reading",
				"but it had no pictures or conversations",
				"what is the use of a book",
				"without pictures or conversations",
				"So she was considering in her own mind",
				"whether the pleasure of making a daisy chain",
				"would be worth the trouble of getting up",
				"when suddenly a White Rabbit ran close by",
				"There was nothing so very remarkable",
				"nor did Alice think it so very odd",
				"to hear the Rabbit say to itself",
				"Oh dear Oh dear I shall be late",
				"but when the Rabbit took a watch",
				"out of its waistcoat pocket",
				"Alice started to her feet",
				"burning with curiosity she ran across",
			},
		},
	}
}

// TaskGenerator generates procedural vim training tasks
type TaskGenerator struct {
	sources []TextSource
	rng     *rand.Rand
}

// NewTaskGenerator creates a new task generator
func NewTaskGenerator() *TaskGenerator {
	return &TaskGenerator{
		sources: GetPublicDomainSources(),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewSeededTaskGenerator creates a task generator with a specific seed for reproducibility
func NewSeededTaskGenerator(seed int64) *TaskGenerator {
	return &TaskGenerator{
		sources: GetPublicDomainSources(),
		rng:     rand.New(rand.NewSource(seed)),
	}
}

// GetAttribution returns attribution text for all sources used
func (g *TaskGenerator) GetAttribution() string {
	var lines []string
	lines = append(lines, "Text sources used in MoCaCo (all Public Domain):")
	lines = append(lines, "")
	for _, src := range g.sources {
		lines = append(lines, fmt.Sprintf("- '%s' by %s (%d)", src.Name, src.Author, src.Year))
	}
	lines = append(lines, "")
	lines = append(lines, "All texts sourced from Project Gutenberg (https://www.gutenberg.org)")
	lines = append(lines, "These works are in the public domain in the United States.")
	return strings.Join(lines, "\n")
}

// randomSentence returns a random sentence from the sources
func (g *TaskGenerator) randomSentence() string {
	source := g.sources[g.rng.Intn(len(g.sources))]
	return source.Sentences[g.rng.Intn(len(source.Sentences))]
}

// randomWord returns a random word from a sentence
func (g *TaskGenerator) randomWord(sentence string) (word string, startIdx int) {
	words := strings.Fields(sentence)
	if len(words) == 0 {
		return "", 0
	}
	wordIdx := g.rng.Intn(len(words))
	word = words[wordIdx]

	// Find the start index of this word in the sentence
	startIdx = 0
	for i := 0; i < wordIdx; i++ {
		startIdx = strings.Index(sentence[startIdx:], words[i]) + len(words[i]) + startIdx
		// Skip whitespace
		for startIdx < len(sentence) && sentence[startIdx] == ' ' {
			startIdx++
		}
	}
	startIdx = strings.Index(sentence, word)

	return word, startIdx
}

// GenerateMotionTask generates a motion task
func (g *TaskGenerator) GenerateMotionTask(difficulty int) Task {
	sentence := g.randomSentence()
	words := strings.Fields(sentence)

	var task Task
	task.Category = CategoryMotion
	task.Difficulty = difficulty
	task.Initial = sentence
	task.Desired = sentence // Same for motion tasks

	switch difficulty {
	case 1:
		// Simple word motion: w, b, e, 0, $
		motions := []struct {
			name        string
			keys        string
			description string
			hint        string
			setup       func() (cursorStart, cursorEnd int)
		}{
			{
				name: "w", keys: "w", description: "Move to next word",
				hint: "Use 'w' to move to the start of the next word",
				setup: func() (int, int) {
					if len(words) < 2 {
						return 0, len(sentence) - 1
					}
					startIdx := 0
					endIdx := strings.Index(sentence, words[1])
					return startIdx, endIdx
				},
			},
			{
				name: "e", keys: "e", description: "Move to end of word",
				hint: "Use 'e' to move to the end of the current word",
				setup: func() (int, int) {
					if len(words) < 1 {
						return 0, 0
					}
					return 0, len(words[0]) - 1
				},
			},
			{
				name: "0", keys: "0", description: "Move to start of line",
				hint: "Use '0' to move to the beginning of the line",
				setup: func() (int, int) {
					startPos := len(sentence) / 2
					if startPos < 1 {
						startPos = 1
					}
					return startPos, 0
				},
			},
			{
				name: "$", keys: "$", description: "Move to end of line",
				hint: "Use '$' to move to the end of the line",
				setup: func() (int, int) {
					return 0, len(sentence) - 1
				},
			},
			{
				name: "b", keys: "b", description: "Move to previous word",
				hint: "Use 'b' to move back to the start of the previous word",
				setup: func() (int, int) {
					if len(words) < 2 {
						return len(sentence) - 1, 0
					}
					// Start at second word, move to first
					startIdx := strings.Index(sentence, words[1])
					return startIdx, 0
				},
			},
		}

		choice := motions[g.rng.Intn(len(motions))]
		task.CursorStart, task.CursorEnd = choice.setup()
		task.OptimalKeys = choice.keys
		task.OptimalCount = len(choice.keys)
		task.Description = choice.description
		task.Hint = choice.hint
		task.ID = fmt.Sprintf("gen-motion-%s-%d", choice.name, g.rng.Int())

	case 2:
		// Motion with counts or find
		if len(words) >= 3 && g.rng.Float32() < 0.5 {
			// Count motion: 2w, 3w
			count := 2
			if len(words) > 3 {
				count = 2 + g.rng.Intn(2) // 2 or 3
			}
			task.CursorStart = 0
			targetWord := words[min(count, len(words)-1)]
			task.CursorEnd = strings.Index(sentence, targetWord)
			task.OptimalKeys = fmt.Sprintf("%dw", count)
			task.OptimalCount = 2
			task.Description = fmt.Sprintf("Move forward %d words", count)
			task.Hint = fmt.Sprintf("Use '%dw' to move forward %d words", count, count)
			task.ID = fmt.Sprintf("gen-motion-%dw-%d", count, g.rng.Int())
		} else {
			// Find motion: f{char}
			word, wordStart := g.randomWord(sentence)
			if len(word) > 0 && wordStart > 0 {
				targetChar := word[0]
				task.CursorStart = 0
				task.CursorEnd = wordStart
				task.OptimalKeys = fmt.Sprintf("f%c", targetChar)
				task.OptimalCount = 2
				task.Description = fmt.Sprintf("Find '%c'", targetChar)
				task.Hint = fmt.Sprintf("Use 'f%c' to jump to the next '%c'", targetChar, targetChar)
				task.ID = fmt.Sprintf("gen-motion-f%c-%d", targetChar, g.rng.Int())
			} else {
				// Fallback to $ motion
				task.CursorStart = 0
				task.CursorEnd = len(sentence) - 1
				task.OptimalKeys = "$"
				task.OptimalCount = 1
				task.Description = "Move to end of line"
				task.Hint = "Use '$' to move to the end of the line"
				task.ID = fmt.Sprintf("gen-motion-dollar-%d", g.rng.Int())
			}
		}

	default:
		// Advanced: t motion, multiple finds
		word, wordStart := g.randomWord(sentence)
		if len(word) > 0 && wordStart > 1 {
			targetChar := word[0]
			task.CursorStart = 0
			task.CursorEnd = wordStart - 1
			task.OptimalKeys = fmt.Sprintf("t%c", targetChar)
			task.OptimalCount = 2
			task.Description = fmt.Sprintf("Move until '%c'", targetChar)
			task.Hint = fmt.Sprintf("Use 't%c' to move to just before '%c'", targetChar, targetChar)
			task.ID = fmt.Sprintf("gen-motion-t%c-%d", targetChar, g.rng.Int())
		} else {
			task.CursorStart = len(sentence) - 1
			task.CursorEnd = 0
			task.OptimalKeys = "0"
			task.OptimalCount = 1
			task.Description = "Move to start of line"
			task.Hint = "Use '0' to move to the beginning"
			task.ID = fmt.Sprintf("gen-motion-zero-%d", g.rng.Int())
		}
	}

	task.Tags = []string{"motion", "procedural"}
	return task
}

// GenerateDeleteTask generates a delete task
func (g *TaskGenerator) GenerateDeleteTask(difficulty int) Task {
	sentence := g.randomSentence()
	words := strings.Fields(sentence)

	var task Task
	task.Category = CategoryDelete
	task.Difficulty = difficulty
	task.Tags = []string{"delete", "procedural"}

	switch difficulty {
	case 1:
		// Simple delete: x, dw, dd
		if len(words) >= 2 {
			// Delete a word with dw
			wordIdx := g.rng.Intn(len(words) - 1) // Not the last word
			wordToDelete := words[wordIdx]
			startIdx := strings.Index(sentence, wordToDelete)

			task.Initial = sentence
			// Remove the word and the space after it
			task.Desired = strings.Replace(sentence, wordToDelete+" ", "", 1)
			task.CursorStart = startIdx
			// Highlight the word to be deleted (including trailing space)
			task.HighlightStart = startIdx
			task.HighlightEnd = startIdx + len(wordToDelete) + 1 // +1 for space
			task.OptimalKeys = "dw"
			task.OptimalCount = 2
			task.Description = "Delete word"
			task.Hint = "Use 'dw' to delete the word under the cursor"
			task.ID = fmt.Sprintf("gen-delete-dw-%d", g.rng.Int())
		} else {
			// Delete character with x
			if len(sentence) > 1 {
				pos := g.rng.Intn(len(sentence))
				task.Initial = sentence
				task.Desired = sentence[:pos] + sentence[pos+1:]
				task.CursorStart = pos
				// Highlight the character to be deleted
				task.HighlightStart = pos
				task.HighlightEnd = pos + 1
				task.OptimalKeys = "x"
				task.OptimalCount = 1
				task.Description = "Delete character"
				task.Hint = "Use 'x' to delete the character under the cursor"
				task.ID = fmt.Sprintf("gen-delete-x-%d", g.rng.Int())
			}
		}

	case 2:
		// Text object delete: daw, diw
		if len(words) >= 2 {
			wordIdx := g.rng.Intn(len(words))
			wordToDelete := words[wordIdx]
			startIdx := strings.Index(sentence, wordToDelete)
			// Position cursor in middle of word
			cursorPos := startIdx + len(wordToDelete)/2

			task.Initial = sentence
			// daw removes word and surrounding space
			if wordIdx == 0 {
				task.Desired = strings.TrimPrefix(sentence, wordToDelete+" ")
				task.HighlightStart = 0
				task.HighlightEnd = len(wordToDelete) + 1
			} else if wordIdx == len(words)-1 {
				task.Desired = strings.TrimSuffix(sentence, " "+wordToDelete)
				task.HighlightStart = startIdx - 1 // Include leading space
				task.HighlightEnd = startIdx + len(wordToDelete)
			} else {
				task.Desired = strings.Replace(sentence, " "+wordToDelete, "", 1)
				task.HighlightStart = startIdx - 1 // Include leading space
				task.HighlightEnd = startIdx + len(wordToDelete)
			}
			task.CursorStart = cursorPos
			task.OptimalKeys = "daw"
			task.OptimalCount = 3
			task.Description = "Delete a word"
			task.Hint = "Use 'daw' to delete 'a word' including surrounding space"
			task.ID = fmt.Sprintf("gen-delete-daw-%d", g.rng.Int())
		}

	default:
		// Advanced: dt, df
		// Find a punctuation or specific character to delete until
		word, wordStart := g.randomWord(sentence)
		if len(word) > 0 && wordStart > 0 {
			targetChar := word[0]
			task.Initial = sentence
			task.Desired = sentence[wordStart:]
			task.CursorStart = 0
			// Highlight from cursor to target
			task.HighlightStart = 0
			task.HighlightEnd = wordStart
			task.OptimalKeys = fmt.Sprintf("dt%c", targetChar)
			task.OptimalCount = 3
			task.Description = fmt.Sprintf("Delete until '%c'", targetChar)
			task.Hint = fmt.Sprintf("Use 'dt%c' to delete until '%c'", targetChar, targetChar)
			task.ID = fmt.Sprintf("gen-delete-dt-%d", g.rng.Int())
		}
	}

	return task
}

// GenerateChangeTask generates a change task
func (g *TaskGenerator) GenerateChangeTask(difficulty int) Task {
	sentence := g.randomSentence()
	words := strings.Fields(sentence)
	replacementWords := []string{"new", "changed", "updated", "modified", "different"}

	var task Task
	task.Category = CategoryChange
	task.Difficulty = difficulty
	task.Tags = []string{"change", "procedural"}

	replacement := replacementWords[g.rng.Intn(len(replacementWords))]

	switch difficulty {
	case 1:
		// Simple change: cw, r
		if len(words) >= 1 {
			wordIdx := g.rng.Intn(len(words))
			oldWord := words[wordIdx]
			startIdx := strings.Index(sentence, oldWord)

			task.Initial = sentence
			task.Desired = strings.Replace(sentence, oldWord, replacement, 1)
			task.CursorStart = startIdx
			// Highlight the word to be changed
			task.HighlightStart = startIdx
			task.HighlightEnd = startIdx + len(oldWord)
			task.OptimalKeys = fmt.Sprintf("cw%s<ESC>", replacement)
			task.OptimalCount = 2 + len(replacement) + 1 // cw + word + ESC
			task.Description = fmt.Sprintf("Change word to '%s'", replacement)
			task.Hint = "Use 'cw' to change the word, type the new word, press ESC"
			task.ID = fmt.Sprintf("gen-change-cw-%d", g.rng.Int())
		}

	case 2:
		// Text object change: ciw
		if len(words) >= 1 {
			wordIdx := g.rng.Intn(len(words))
			oldWord := words[wordIdx]
			startIdx := strings.Index(sentence, oldWord)
			// Position cursor in middle of word
			cursorPos := startIdx + len(oldWord)/2

			task.Initial = sentence
			task.Desired = strings.Replace(sentence, oldWord, replacement, 1)
			task.CursorStart = cursorPos
			// Highlight the word to be changed
			task.HighlightStart = startIdx
			task.HighlightEnd = startIdx + len(oldWord)
			task.OptimalKeys = fmt.Sprintf("ciw%s<ESC>", replacement)
			task.OptimalCount = 3 + len(replacement) + 1
			task.Description = fmt.Sprintf("Change inner word to '%s'", replacement)
			task.Hint = "Use 'ciw' to change the word regardless of cursor position"
			task.ID = fmt.Sprintf("gen-change-ciw-%d", g.rng.Int())
		}

	default:
		// Advanced: C, cc - entire line
		task.Initial = sentence
		task.Desired = replacement
		task.CursorStart = 0
		// Highlight entire line
		task.HighlightStart = 0
		task.HighlightEnd = len(sentence)
		task.OptimalKeys = fmt.Sprintf("cc%s<ESC>", replacement)
		task.OptimalCount = 2 + len(replacement) + 1
		task.Description = "Change entire line"
		task.Hint = "Use 'cc' to change the entire line"
		task.ID = fmt.Sprintf("gen-change-cc-%d", g.rng.Int())
	}

	return task
}

// GenerateInsertTask generates an insert task
func (g *TaskGenerator) GenerateInsertTask(difficulty int) Task {
	sentence := g.randomSentence()
	insertWords := []string{"very", "quite", "rather", "extremely", "somewhat"}

	var task Task
	task.Category = CategoryInsert
	task.Difficulty = difficulty
	task.Tags = []string{"insert", "procedural"}

	insertion := insertWords[g.rng.Intn(len(insertWords))]

	switch difficulty {
	case 1:
		// Simple insert: i, a, A
		words := strings.Fields(sentence)
		if len(words) >= 2 {
			// Insert before second word
			secondWord := words[1]
			insertPos := strings.Index(sentence, secondWord)

			task.Initial = sentence
			task.Desired = sentence[:insertPos] + insertion + " " + sentence[insertPos:]
			task.CursorStart = insertPos
			task.OptimalKeys = fmt.Sprintf("i%s <ESC>", insertion)
			task.OptimalCount = 1 + len(insertion) + 1 + 1
			task.Description = fmt.Sprintf("Insert '%s' before cursor", insertion)
			task.Hint = "Use 'i' to insert before the cursor"
			task.ID = fmt.Sprintf("gen-insert-i-%d", g.rng.Int())
		}

	case 2:
		// Append: A (end of line)
		task.Initial = sentence
		task.Desired = sentence + " " + insertion
		task.CursorStart = 0
		task.OptimalKeys = fmt.Sprintf("A %s<ESC>", insertion)
		task.OptimalCount = 1 + 1 + len(insertion) + 1
		task.Description = "Append at end of line"
		task.Hint = "Use 'A' to append at the end of the line"
		task.ID = fmt.Sprintf("gen-insert-A-%d", g.rng.Int())

	default:
		// Open line: o, O
		task.Initial = sentence
		task.Desired = sentence + "\n" + insertion
		task.CursorStart = 0
		task.OptimalKeys = fmt.Sprintf("o%s<ESC>", insertion)
		task.OptimalCount = 1 + len(insertion) + 1
		task.Description = "Open new line below"
		task.Hint = "Use 'o' to open a new line below and enter insert mode"
		task.ID = fmt.Sprintf("gen-insert-o-%d", g.rng.Int())
	}

	return task
}

// GenerateVisualTask generates a visual mode task
func (g *TaskGenerator) GenerateVisualTask(difficulty int) Task {
	sentence := g.randomSentence()
	words := strings.Fields(sentence)

	var task Task
	task.Category = CategoryVisual
	task.Difficulty = max(difficulty, 2) // Visual is at least level 2
	task.Tags = []string{"visual", "procedural"}

	if len(words) >= 2 {
		// Visual delete word
		wordIdx := g.rng.Intn(len(words))
		wordToDelete := words[wordIdx]
		startIdx := strings.Index(sentence, wordToDelete)

		task.Initial = sentence
		// Remove word
		if wordIdx < len(words)-1 {
			task.Desired = strings.Replace(sentence, wordToDelete+" ", "", 1)
		} else {
			task.Desired = strings.Replace(sentence, " "+wordToDelete, "", 1)
		}
		task.CursorStart = startIdx
		task.OptimalKeys = "viwd"
		task.OptimalCount = 4
		task.Description = "Visually select and delete word"
		task.Hint = "Use 'viw' to visually select the word, then 'd' to delete"
		task.ID = fmt.Sprintf("gen-visual-viwd-%d", g.rng.Int())
	}

	return task
}

// GenerateComplexTask generates a complex multi-step task
func (g *TaskGenerator) GenerateComplexTask(difficulty int) Task {
	sentence := g.randomSentence()
	words := strings.Fields(sentence)

	var task Task
	task.Category = CategoryComplex
	task.Difficulty = max(difficulty, 3) // Complex is at least level 3
	task.Tags = []string{"complex", "procedural"}

	if len(words) >= 2 {
		// Swap first two words
		task.Initial = sentence
		newWords := make([]string, len(words))
		copy(newWords, words)
		newWords[0], newWords[1] = newWords[1], newWords[0]
		task.Desired = strings.Join(newWords, " ")
		task.CursorStart = 0
		task.OptimalKeys = "dwwP"
		task.OptimalCount = 4
		task.Description = "Swap first two words"
		task.Hint = "Delete first word, move to next word, paste before"
		task.ID = fmt.Sprintf("gen-complex-swap-%d", g.rng.Int())
	}

	return task
}

// GenerateTasksForRound generates all tasks for a round
func (g *TaskGenerator) GenerateTasksForRound(roundType string) []Task {
	var tasks []Task

	// Distribution per round: 6 motion, 6 delete, 6 change, 6 insert, 3 visual, 3 complex = 30
	distribution := map[TaskCategory]int{
		CategoryMotion:  6,
		CategoryDelete:  6,
		CategoryChange:  6,
		CategoryInsert:  6,
		CategoryVisual:  3,
		CategoryComplex: 3,
	}

	// Difficulty range based on round type
	minDiff, maxDiff := 1, 1
	switch roundType {
	case "beginner":
		minDiff, maxDiff = 1, 1
	case "intermediate":
		minDiff, maxDiff = 1, 2
	case "advanced":
		minDiff, maxDiff = 2, 3
	case "expert":
		minDiff, maxDiff = 3, 4
	case "mixed":
		minDiff, maxDiff = 1, 4
	}

	for cat, count := range distribution {
		for i := 0; i < count; i++ {
			diff := minDiff
			if maxDiff > minDiff {
				diff = minDiff + g.rng.Intn(maxDiff-minDiff+1)
			}

			var task Task
			switch cat {
			case CategoryMotion:
				task = g.GenerateMotionTask(diff)
			case CategoryDelete:
				task = g.GenerateDeleteTask(diff)
			case CategoryChange:
				task = g.GenerateChangeTask(diff)
			case CategoryInsert:
				task = g.GenerateInsertTask(diff)
			case CategoryVisual:
				task = g.GenerateVisualTask(diff)
			case CategoryComplex:
				task = g.GenerateComplexTask(diff)
			}
			tasks = append(tasks, task)
		}
	}

	// Shuffle tasks
	g.rng.Shuffle(len(tasks), func(i, j int) {
		tasks[i], tasks[j] = tasks[j], tasks[i]
	})

	return tasks
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper to check if character is a letter
func isLetter(c byte) bool {
	return unicode.IsLetter(rune(c))
}
