package nonogram

import (
	"fmt"
	"maps"
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

func (l *Line) filterAt(position int, fill bool) *LineChange {
	tailFrom := l.tail
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
	if tailFrom == l.tail {
		return nil
	}
	return &LineChange{
		line:     l,
		tailFrom: tailFrom,
		tailTo:   l.tail,
	}
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

type Cell int

const (
	Unknown Cell = iota
	Empty
	Filled
)

type Action struct {
	cell CellChange
	line *LineChange
}

type CellChange struct {
	value Cell
	row   int
	col   int
}

type LineChange struct {
	line     *Line
	tailFrom int
	tailTo   int
}

type Nonogram struct {
	rows        []*Line
	cols        []*Line
	grid        []Cell
	actions     []Action
	filledCount int

	Debug   bool
	Step    int
	MaxStep int
}

func New(rowClues [][]int, columnClues [][]int) *Nonogram {
	nonogram := &Nonogram{
		rows: make([]*Line, len(rowClues)),
		cols: make([]*Line, len(columnClues)),
		grid: make([]Cell, len(rowClues)*len(columnClues)),
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

func (n *Nonogram) FillCell(line *Line, index int, value Cell) (bool, *LineChange) {
	var row, col int
	var effected *Line
	if line.direction == Row {
		row = line.index
		col = index
		effected = n.cols[col]
	} else {
		row = index
		col = line.index
		effected = n.rows[row]
	}
	cellIndex := row*len(n.cols) + col
	if n.grid[cellIndex] != Unknown {
		if n.grid[cellIndex] != value {
			panic(fmt.Sprintf("Conflict fill at cell (%d,%d), current value: %v, new value: %v", row, col, n.grid[cellIndex], value))
		}
		return true, nil
	}
	n.grid[cellIndex] = value
	n.filledCount++
	n.println(n.GridString())
	action := Action{cell: CellChange{
		value: value,
		row:   row,
		col:   col,
	}}
	action.line = effected.filterAt(line.index, value == Filled)
	if action.line != nil && action.line.tailTo == 0 {
		return false, nil
	}
	n.actions = append(n.actions, action)
	return true, action.line
}

func (n *Nonogram) Savepoint() int {
	return len(n.actions)
}

func (n *Nonogram) Rollback(to int) {
	for i := len(n.actions) - 1; i >= to; i-- {
		action := n.actions[i]
		n.grid[action.cell.row*len(n.cols)+action.cell.col] = Unknown
		n.filledCount--
		line := action.line.line
		line.tail = action.line.tailFrom
	}
	n.actions = n.actions[:to]
}

func (n *Nonogram) IsSolved() bool {
	return n.filledCount == len(n.grid)
}

func (n *Nonogram) IsRowSolved(row int) bool {
	for col := 0; col < len(n.cols); col++ {
		cell := n.grid[row*len(n.cols)+col]
		if cell == Unknown {
			return false
		}
	}
	return true
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
	for i, cell := range n.grid {
		row, col := i/len(n.cols), i%len(n.cols)
		switch {
		case col == 0 && row == 0:
			builder.WriteString("\n|")
		case col == 0:
			builder.WriteString("|\n|")
		default:
			builder.WriteString(" ")
		}
		switch cell {
		case Filled:
			builder.WriteString("o")
		case Empty:
			builder.WriteString("·")
		case Unknown:
			builder.WriteString(" ")
		}
	}
	builder.WriteString("|\n")
	printHorizonBorder()
	return builder.String()
}

func (n *Nonogram) Solve() bool {
	// propagate forced values until no more changes
	beforePropagation := n.Savepoint()
	lineQueue := append([][]*Line{}, n.rows)
	lineQueue = append(lineQueue, n.cols)
	cursor := 0
	for cursor < len(lineQueue) {
		lines := lineQueue[cursor]
		effectedLineMap := map[int]*Line{}
		for _, line := range lines {
			n.incStep()
			ok, effected := n.propagate(line)
			if n.IsSolved() {
				return true
			}
			if !ok {
				n.Rollback(beforePropagation)
				return false
			}
			maps.Copy(effectedLineMap, effected)
		}
		var effectedLines []*Line
		for _, l := range effectedLineMap {
			effectedLines = append(effectedLines, l)
		}
		if len(effectedLines) > 0 {
			lineQueue = append(lineQueue, effectedLines)
		}
		cursor++
	}

	minCandidates := math.MaxUint32
	mrvRow := -1
	for i, row := range n.rows {
		if n.IsRowSolved(i) {
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
	beforeBranches := n.Savepoint()
	for _, pattern := range n.rows[mrvRow].candidates[:minCandidates] {
		n.incStep()
		applyOk := n.applyRow(mrvRow, pattern)
		n.printf("[%d] Try pattern %0*b for row %d, result: %v\n", n.Step, len(n.cols), pattern, mrvRow, applyOk)
		if applyOk {
			solved := n.Solve()
			if solved {
				return true
			} else {
			}
		}
		n.Rollback(beforeBranches)
	}
	n.Rollback(beforePropagation)
	return false
}

// merge two propagate functions into one propagate function
func (n *Nonogram) propagate(source *Line) (bool, map[int]*Line) {
	effected := map[int]*Line{}
	forcedFilled, forcedEmpty := source.forcedValues()
	if forcedFilled == 0 && forcedEmpty == 0 {
		return true, effected
	}
	n.printf("[%d] Propagate %v %d, forced filled: %0*b, forced empty: %0*b\n", n.Step, source.direction, source.index, source.size, forcedFilled, source.size, forcedEmpty)

	for i := 0; i < source.size; i++ {
		fill := (Pattern(1)<<(source.size-1-i))&forcedFilled != 0
		empty := (Pattern(1)<<(source.size-1-i))&forcedEmpty != 0
		var value Cell
		if fill {
			value = Filled
		} else if empty {
			value = Empty
		}
		if value == Unknown {
			continue
		}
		fillOk, lineChange := n.FillCell(source, i, value)
		if !fillOk {
			return false, nil
		}
		if lineChange != nil {
			effected[lineChange.line.index] = lineChange.line
		}
	}
	return true, effected
}

func (n *Nonogram) applyRow(row int, pattern Pattern) bool {
	for i, col := range n.cols {
		// pick i-th (from highest) bit in pattern, to check if need to fill or empty
		fill := (Pattern(1)<<(len(n.cols)-1-i))&pattern != 0
		// move to row-th (from-highest) position, and filter column candidates
		value := Filled
		if !fill {
			value = Empty
		}
		fillOk, _ := n.FillCell(col, row, value)
		if !fillOk {
			return false
		}
	}
	return true
}
