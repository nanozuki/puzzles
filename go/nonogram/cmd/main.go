package main

import (
	"os"

	"github.com/nanozuki/puzzles/go/nonogram"
)

func main() {
	if len(os.Args) != 2 {
		println("usage: nonogram <puzzle>")
		return
	}
	puzzle := nonogram.NewFromString(os.Args[1])
	puzzle.Debug = true
	// puzzle.MaxStep = 1
	ok := puzzle.Solve()
	if !ok {
		println("no solution")
		return
	}
	println(puzzle.GridString())
	println("used", puzzle.Step, "steps")
}
