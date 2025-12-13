package main

import (
	"fmt"
	"os"

	"github.com/nanozuki/puzzles/go/sudoku"
)

// run soduku <puzzle> to solve a sudoku puzzle
// where <puzzle> is a string of 81 characters, each being '1'-'9' or '0' for empty cells

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: sudoku <puzzle>")
		return
	}
	puzzleStr := os.Args[1]
	puzzle := sudoku.NewFromString(puzzleStr)
	puzzle.Verbose = true
	puzzle.Solve()
}
