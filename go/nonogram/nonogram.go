package nonogram

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

// Pattern use uint32 bitmask to represent a line pattern, 1 for filled, 0 for empty, use lowest N bits to represent a
// line, N is the length of the line. From highest bit to lowest bit we used, it represents from left to right for rows,
// and from top to bottom for columns.
type Pattern uint32

// Line of nonogram puzzle
// Contains all necessary patterns which created by clues
type Line struct {
	size       int
	candidates []Pattern
	tail       int
}

func NewLine(size int, clues []int) *Line {
	line := &Line{size: size}
	line.candidates = generateCandidates(size, clues)
	line.tail = len(line.candidates)
	return line
}

func generateCandidates(size int, clues []int) []Pattern {
	var candidates []Pattern
	if len(clues) == 0 {
		panic("clues should not be empty")
	}

	if slices.Contains(clues, 0) {
		if len(clues) > 1 {
			panic("zero clue should be only clue in a line")
		}
		return []Pattern{0}
	}

	sum := 0
	for _, clue := range clues {
		sum += clue
	}
	freeSpaces := size - sum - len(clues) + 1

	for i := 0; i <= freeSpaces; i++ {
		offset := size - clues[0] - i
		pattern := (Pattern(0x1)<<clues[0] - 1) << offset
		if len(clues) > 1 {
			// offset also means the rest size, and we need to minus 1 for a empty cell between blocks
			subCandidates := generateCandidates(offset-1, clues[1:])
			for _, subPattern := range subCandidates {
				candidates = append(candidates, pattern|subPattern)
			}
		} else {
			candidates = append(candidates, pattern)
		}
	}
	return candidates
}

func (l *Line) filterAt(position int, fill bool) tailChange {
	change := tailChange{from: l.tail, to: l.tail}
	mask := Pattern(1) << (l.size - 1 - position)
	for i := 0; i < l.tail; i++ {
		for (fill && l.candidates[i]&mask == 0) || (!fill && l.candidates[i]&mask != 0) {
			if i == l.tail-1 {
				l.tail--
				break
			}
			l.candidates[i], l.candidates[l.tail-1] = l.candidates[l.tail-1], l.candidates[i]
			l.tail--
		}
	}
	change.to = l.tail
	return change
}

type Nonogram struct {
	rows       []*Line
	cols       []*Line
	sovledRows map[int]Pattern

	Debug   bool
	Step    int
	MaxStep int
}

func New(rowClues [][]int, columnClues [][]int) *Nonogram {
	nonogram := &Nonogram{
		rows:       make([]*Line, len(rowClues)),
		cols:       make([]*Line, len(columnClues)),
		sovledRows: make(map[int]Pattern),
	}
	for i, clues := range rowClues {
		nonogram.rows[i] = NewLine(len(columnClues), clues)
	}
	for i, clues := range columnClues {
		nonogram.cols[i] = NewLine(len(rowClues), clues)
	}
	fmt.Printf("New Nonogram with %d rows and %d columns\n", len(nonogram.rows), len(nonogram.cols))
	return nonogram
}

// NewFromString create Nonogram from string in specific format:
// - use "=" separate row clues and column clues
// - use ";" separate different lines of clues
// - use "," separate different clues in a line
// for example: "3,1;2=1;3" means a nonogram with 2 rows and 2 columns,
// first row has clues [3,1], second row has clues [2],
// first column has clues [1], second column has clues [3]
func NewFromString(s string) *Nonogram {
	parts := strings.Split(s, "=")
	if len(parts) != 2 {
		panic("invalid nonogram string format")
	}
	parseClues := func(clues string) []int {
		parts := strings.Split(clues, ",")
		result := make([]int, len(parts))
		for i, part := range parts {
			clue, err := strconv.Atoi(part)
			if err != nil {
				panic("invalid clue number")
			}
			result[i] = clue
		}
		return result
	}
	parseLine := func(line string) [][]int {
		parts := strings.Split(line, ";")
		result := make([][]int, len(parts))
		for i, part := range parts {
			result[i] = parseClues(part)
		}
		return result
	}
	rowClues := parseLine(parts[0])
	columnClues := parseLine(parts[1])
	return New(rowClues, columnClues)
}

func (n *Nonogram) println(a ...any) {
	if n.Debug {
		fmt.Println(a...)
	}
}

func (n *Nonogram) printf(format string, a ...any) {
	if n.Debug {
		fmt.Printf(format, a...)
	}
}

func (n *Nonogram) GridString() string {
	builder := strings.Builder{}
	printHorizonBorder := func() {
		builder.WriteString("+")
		for range n.cols {
			builder.WriteString("-+")
		}
	}
	printHorizonBorder()
	builder.WriteString("\n")
	for row := range len(n.rows) {
		for col := 0; col < len(n.cols); col++ {
			if col == 0 {
				builder.WriteString("|")
			} else {
				builder.WriteString(" ")
			}
			filled := false
			if pattern, ok := n.sovledRows[row]; ok {
				filled = pattern&(Pattern(1)<<(len(n.cols)-1-col)) != 0
			}
			if filled {
				builder.WriteString("o")
			} else {
				builder.WriteString("•")
			}
		}
		builder.WriteString("|\n")
	}
	printHorizonBorder()
	return builder.String()
}

func (n *Nonogram) Solve() bool {
	minCandidates := math.MaxUint32
	mrvRow := -1
	for i, row := range n.rows {
		if _, ok := n.sovledRows[i]; ok {
			continue
		}
		if len(row.candidates) < minCandidates {
			minCandidates = len(row.candidates)
			mrvRow = i
		}
	}
	if mrvRow == -1 { // all rows are solved
		return true
	}

	n.println("Try to fill row", mrvRow, "with", minCandidates, "candidates")
	for _, pattern := range n.rows[mrvRow].candidates {
		if n.MaxStep > 0 && n.Step > n.MaxStep {
			panic("exceed max step " + strconv.Itoa(n.MaxStep))
		}
		fillOk, changes := n.fillRow(mrvRow, pattern)
		n.Step++
		n.println("  Try pattern", fmt.Sprintf("%0*b", len(n.cols), pattern), "fillOk:", fillOk)
		if fillOk {
			n.sovledRows[mrvRow] = pattern
			n.println(n.GridString())
			solved := n.Solve()
			if solved {
				return true
			} else {
				delete(n.sovledRows, mrvRow)
			}
		}
		n.rollbackChanges(changes)
	}
	return false
}

type tailChange struct {
	from int
	to   int
}

func (n *Nonogram) fillRow(row int, pattern Pattern) (bool, map[int]tailChange) {
	ok, changes := true, make(map[int]tailChange)
	for i, col := range n.cols {
		// pick i-th (from highest) bit in pattern, to check if need to fill or empty
		fill := (Pattern(1)<<(len(n.cols)-1-i))&pattern != 0
		// move to row-th (from-highest) position, and filter column candidates
		change := col.filterAt(row, fill)
		n.printf("    Fill %v column %d at row %d, candidates from %d to %d\n", fill, i, row, change.from, change.to)
		if change.from != change.to {
			changes[i] = change
		}
		if change.to == 0 {
			ok = false
			break
		}
	}
	return ok, changes
}

func (n *Nonogram) rollbackChanges(changes map[int]tailChange) {
	for i, change := range changes {
		col := n.cols[i]
		col.tail = change.from
	}
}
