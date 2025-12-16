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

type Direction int

const (
	Row Direction = iota
	Column
)

func (d Direction) String() string {
	switch d {
	case Row:
		return "Row"
	case Column:
		return "Column"
	default:
		panic("unknown direction")
	}
}

// Line of nonogram puzzle
// Contains all necessary patterns which created by clues
type Line struct {
	direction  Direction
	index      int
	size       int
	candidates []Pattern
	tail       int
}

func NewLine(dir Direction, index, size int, clues []int) *Line {
	line := &Line{
		direction: dir,
		index:     index,
		size:      size,
	}
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

// forcedValues returns two patterns, forcedFilled and forcedEmpty
// each 1 in forcedFilled means that cell must be filled in all candidates
// each 1 in forcedEmpty means that cell must be empty in all candidates
func (l *Line) forcedValues() (Pattern, Pattern) {
	lineMask := Pattern(1)<<l.size - 1
	forcedFilled, forcedEmpty := lineMask, Pattern(0)
	for _, pattern := range l.candidates[:l.tail] {
		forcedFilled &= pattern
		forcedEmpty |= pattern
	}
	forcedEmpty = ^forcedEmpty & lineMask
	return forcedFilled, forcedEmpty
}

type Nonogram struct {
	rows       []*Line
	cols       []*Line
	solvedRows map[int]Pattern

	Debug   bool
	Step    int
	MaxStep int
}

func New(rowClues [][]int, columnClues [][]int) *Nonogram {
	nonogram := &Nonogram{
		rows:       make([]*Line, len(rowClues)),
		cols:       make([]*Line, len(columnClues)),
		solvedRows: make(map[int]Pattern),
	}
	for i, clues := range rowClues {
		nonogram.rows[i] = NewLine(Row, i, len(columnClues), clues)
	}
	for i, clues := range columnClues {
		nonogram.cols[i] = NewLine(Column, i, len(rowClues), clues)
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

func (n *Nonogram) incStep() {
	if n.MaxStep > 0 && n.Step > n.MaxStep {
		panic("exceed max step " + strconv.Itoa(n.MaxStep))
	}
	n.Step++
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
			if pattern, ok := n.solvedRows[row]; ok {
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
	// propagate forced values until no more changes
	colsChanges, rowsChanges := map[int]tailChange{}, map[int]tailChange{}
	for {
		curRowsChanges := map[int]tailChange{}
		for _, row := range n.rows {
			n.incStep()
			ok, changes := n.propagate(row)
			if !ok {
				n.rollbackChanges(rowsChanges, colsChanges)
				return false
			}
			mergeChanges(curRowsChanges, changes)
			if row.tail == 1 {
				n.solvedRows[row.index] = row.candidates[0]
			}
		}
		mergeChanges(rowsChanges, curRowsChanges)

		curColsChanges := map[int]tailChange{}
		for _, col := range n.cols {
			n.incStep()
			ok, changes := n.propagate(col)
			if !ok {
				n.rollbackChanges(rowsChanges, colsChanges)
				return false
			}
			mergeChanges(curColsChanges, changes)
		}
		mergeChanges(colsChanges, curColsChanges)
		if len(curRowsChanges) == 0 && len(curColsChanges) == 0 {
			break
		}
	}

	minCandidates := math.MaxUint32
	mrvRow := -1
	for i, row := range n.rows {
		if _, ok := n.solvedRows[i]; ok {
			continue
		}
		if row.tail < minCandidates {
			minCandidates = row.tail
			mrvRow = i
		}
	}
	if mrvRow == -1 { // all rows are solved
		return true
	}

	n.println("Try to fill row", mrvRow, "with", minCandidates, "candidates")
	for _, pattern := range n.rows[mrvRow].candidates[:minCandidates] {
		n.incStep()
		fillOk, branchColsChanges := n.applyRow(mrvRow, pattern)
		n.println("  Try pattern", fmt.Sprintf("%0*b", len(n.cols), pattern), "fillOk:", fillOk)
		if fillOk {
			n.solvedRows[mrvRow] = pattern
			n.println(n.GridString())
			solved := n.Solve()
			if solved {
				return true
			} else {
				delete(n.solvedRows, mrvRow)
			}
		}
		n.rollbackChanges(map[int]tailChange{}, branchColsChanges)
	}
	n.rollbackChanges(rowsChanges, colsChanges)
	return false
}

type tailChange struct {
	from int
	to   int
}

func mergeChanges(target, delta map[int]tailChange) {
	for i, dChange := range delta {
		if tChange, ok := target[i]; ok {
			tChange.to = dChange.to
			target[i] = tChange
		} else {
			target[i] = dChange
		}
	}
}

// merge two propagate functions into one propagate function
func (n *Nonogram) propagate(source *Line) (bool, map[int]tailChange) {
	ok, changes := true, make(map[int]tailChange)
	forcedFilled, forcedEmpty := source.forcedValues()
	if forcedFilled == 0 && forcedEmpty == 0 {
		return ok, changes
	}
	n.printf("  Propagate %v %d, forced filled: %0*b, forced empty: %0*b\n", source.direction, source.index, source.size, forcedFilled, source.size, forcedEmpty)
	targetLines := n.cols // source.direction == Row
	if source.direction == Column {
		targetLines = n.rows
	}
	for i, target := range targetLines {
		fill := (Pattern(1)<<(source.size-1-i))&forcedFilled != 0
		empty := (Pattern(1)<<(source.size-1-i))&forcedEmpty != 0
		if fill {
			change := target.filterAt(source.index, true)
			n.printf("    Forced fill %v %d at %v %d, candidates from %d to %d\n", target.direction, target.index, source.direction, source.index, change.from, change.to)
			if change.from != change.to {
				changes[i] = change
			}
			if change.to == 0 {
				ok = false
				break
			}
		} else if empty {
			change := target.filterAt(source.index, false)
			n.printf("    Forced empty %v %d at %v %d, candidates from %d to %d\n", target.direction, target.index, source.direction, source.index, change.from, change.to)
			if change.from != change.to {
				changes[i] = change
			}
			if change.to == 0 {
				ok = false
				break
			}
		}
	}
	return ok, changes
}

func (n *Nonogram) applyRow(row int, pattern Pattern) (bool, map[int]tailChange) {
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

func (n *Nonogram) rollbackChanges(rowsChanges, colsChanges map[int]tailChange) {
	for i, change := range rowsChanges {
		row := n.rows[i]
		row.tail = change.from
	}
	for i, change := range colsChanges {
		col := n.cols[i]
		col.tail = change.from
	}
}
