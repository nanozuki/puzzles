package sudoku

import (
	"fmt"
	"math/bits"
	"strings"
)

type Sudoku struct {
	board       [81]int // from top-left to bottom-right
	rowMask     [9]int  // from top to bottom
	colMask     [9]int  // from left to right
	regionMask  [9]int  // from top-left to bottom-right
	searchStack []action

	Verbose bool
}

type action struct {
	index  int
	number int
}

func NewFromString(s string) *Sudoku {
	if len(s) != 81 {
		panic("Input string must be exactly 81 characters long")
	}
	sudoku := &Sudoku{}
	for i, ch := range s {
		if ch >= '1' && ch <= '9' {
			number := int(ch - '0')
			row := i / 9
			col := i % 9
			sudoku.fill(row, col, number)
		}
	}
	return sudoku
}

func (s *Sudoku) String() string {
	builder := strings.Builder{}
	for _, number := range s.board {
		builder.WriteByte(byte(number + '0'))
	}
	return builder.String()
}

func (s *Sudoku) PrettifyString() string {
	builder := strings.Builder{}
	for i, number := range s.board {
		row, col := i/9, i%9
		if number == 0 {
			builder.WriteByte(' ')
		} else {
			builder.WriteByte(byte(s.board[i] + '0'))
		}
		switch col {
		case 2, 5:
			builder.WriteByte('|') // region's vertical separator
		case 8:
			switch row {
			case 2, 5:
				builder.WriteString("\n---+---+---\n") // region's horizontal separator
			case 0, 1, 3, 4, 6, 7:
				builder.WriteByte('\n')
			case 8:
			}
		case 0, 1, 3, 4, 6, 7:
		}
	}
	return builder.String()
}

func (s *Sudoku) Solve() bool {
	// 1. Iterator all cells to find the minimum candidate cell
	minCandidates := 10
	minIndex := -1
	for i, number := range s.board {
		if number == 0 {
			row, col, reg := i/9, i%9, (i/9)/3*3+(i%9)/3
			n := bits.OnesCount(uint(^(s.rowMask[row] | s.colMask[col] | s.regionMask[reg]) & 0x1FF))
			if n == 0 {
				return false // no candidate, backtrack
			}
			if n < minCandidates {
				minCandidates = n
				minIndex = i
			}
		}
	}
	if minIndex == -1 {
		return true // solved
	}

	// 2. Get Candidates for the selected cell
	row, col, reg := minIndex/9, minIndex%9, (minIndex/9)/3*3+(minIndex%9)/3
	candidatesMask := ^(s.rowMask[row] | s.colMask[col] | s.regionMask[reg]) & 0x1FF

	// 3. Try each candidate
	for c := range 9 {
		if candidatesMask&(1<<c) != 0 {
			number := c + 1
			s.fill(row, col, number)
			s.searchStack = append(s.searchStack, action{index: minIndex, number: number})
			if s.Verbose {
				fmt.Printf("Fill %d at (%d, %d)\n", number, row, col)
				fmt.Println(s.PrettifyString())
			}
			if s.Solve() {
				return true
			} else {
				s.erase(row, col)
				s.searchStack = s.searchStack[:len(s.searchStack)-1]
				if s.Verbose {
					fmt.Printf("Erase %d at (%d, %d)\n", number, row, col)
					fmt.Println(s.PrettifyString())
				}
			}
		}
	}
	return false
}

func (s *Sudoku) fill(row, col, number int) {
	index := row*9 + col
	if s.board[index] != 0 {
		panic("Cell is already filled")
	}
	s.board[index] = number
	s.rowMask[row] |= 1 << (number - 1)
	s.colMask[col] |= 1 << (number - 1)
	s.regionMask[(row/3)*3+(col/3)] |= 1 << (number - 1)
}

func (s *Sudoku) erase(row, col int) {
	index := row*9 + col
	if s.board[index] == 0 {
		panic("Cell is already empty")
	}
	number := s.board[index]
	s.board[index] = 0
	s.rowMask[row] &^= 1 << (number - 1)
	s.colMask[col] &^= 1 << (number - 1)
	s.regionMask[(row/3)*3+(col/3)] &^= 1 << (number - 1)
}
