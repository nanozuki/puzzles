package nonogram

import (
	"strings"
	"testing"
)

func TestNonogram_Solve(t *testing.T) {
	tests := []struct {
		name       string
		puzzle     string
		wantSolved bool
		wantGrid   string
	}{
		{
			name:       "P001", // 5x5, easy",
			puzzle:     "3;1,1;3;1;2=0;3;1,3;3,1;0",
			wantSolved: true,
			wantGrid: strings.TrimSpace(`
+-+-+-+-+-+
|• o o o •|
|• o • o •|
|• o o o •|
|• • o • •|
|• • o o •|
+-+-+-+-+-+`),
		},
		{
			name:       "P060", // 10x10, hard
			puzzle:     "3,2;1,1,1,2;2,2,1;1,2,1;2,3;1,1,1,1;4,1,1;1,2;1,1;2,1,2=1,2;1,1,1;3,2,1;1,1,1,1;3,2,1,1;1,2,1;1,1,1;1,1,3;1,2,1,1;1,1,2",
			wantSolved: true,
			wantGrid: strings.TrimSpace(`
+-+-+-+-+-+-+-+-+-+-+
|• • • o o o • • o o|
|o • o • o • o o • •|
|• o o • o o • • • o|
|• • o • • o o • o •|
|• • • o o • • o o o|
|o • o • o • • • • o|
|o o o o • o • • o •|
|• • • • o • o o • •|
|• • • o • • • o • •|
|• o o • o • • o o •|
+-+-+-+-+-+-+-+-+-+-+`),
		},
		{
			name:       "P097", // 15x15, easy
			puzzle:     "1,1;4,2;6,1;3,1,1,1;2,4;2,3;6,1,1;2,2,3;2,1,3;2,1,2,1;2,2,1;1,1,1,1;1,1,2,1;2,1,1;7=2;1,1;5;2;2;2;5,3;7,3;3,2,3;2,1,1;6,1,1;1,1,2,2,1;1,6,1;3,1,6,1;1,2,2",
			wantSolved: true,
			wantGrid: strings.TrimSpace(`
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|• • • • • • • • • • • o • o •|
|• • • • • • • o o o o • o o •|
|• • • • • • o o o o o o • o •|
|• • • • • • o o o • o • o • o|
|• • • • • • o o • • o o o o •|
|• • • • • • o o • • o o o • •|
|• • • • • o o o o o o • o • o|
|• • • • o o • o o • • • o o o|
|• • • o o • o • • • • o o o •|
|• • o o • • o • • • o o • o •|
|• o o • • • o o • • • • • o •|
|o • o • • • • o • • • • • o •|
|o • o • • • • o o • • • • o •|
|• o o • • • • • o • • • • • o|
|• • • • • • • • o o o o o o o|
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+`),
		},
		{
			name:       "P136", // 15*20, easy
			puzzle:     "4,3,3;4,11;2,3,1,2,4;5,4,3,3;2,1,2,2,2;5,1,2,2,1;5,2,6,3;1,1,2,3,8;3,1,1;3,2;2,2;2,2;4,2;4,1,3;2,1,1,2,1=2;1;2,3;1,2,2;2,5,2;6,8;4,2,6;2,1,4,3;1,1,1,1,1;2,2,1,2;5,1;1,1,3;2,4;5,3,1;2,1,2,1,2;4,9;2,1,1,5;4,2;4,2;1,1,2",
			wantSolved: true,
			wantGrid: strings.TrimSpace(`
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|• • • o o o o • o o o • • o o o • • • •|
|• • • • o o o o • o o o o o o o o o o o|
|• • o o • o o o • • o • o o • o o o o •|
|• • o o o o o • o o o o • o o o • o o o|
|• • • • o o • o • o o • o o • • • o o •|
|• • o o o o o • o • • o o • • o o • • o|
|o o o o o • o o • • o o o o o o • o o o|
|o • o • o o • o o o • o o o o o o o o •|
|• • • • • o o o • • • • • o • o • • • •|
|• • • • • o o o • • • • • • o o • • • •|
|• • • • • o o • • • • • • • • o o • • •|
|• • • • • o o • • • • • • • • o o • • •|
|• • • • • o o o o • • • • • • o o • • •|
|• • • • o o o o • o • • • • o o o • • •|
|• • • • o o • o • o • • • o o • o • • •|
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			`),
		},
		{
			name:       "P150", // 15*20,hard
			puzzle:     "1,1,1,1,1;5,1,3,3;1,5,2,2;1,1,1,4,1,1;2,5,3;3,2,1,1,1,3;1,1,1,4,1,2;1,1,1,1,1,2,1;2,5,3;5,1,1,1,1,2;1,1,1,1,3,4;1,1,2,2,3,2,1;2,1,11;1,1,1,6,2;3,1,1,4,1,1=2;2,1,1;1,2,3;2,1,3,1;1,2,2,2;1,2,1,1,2,1;1,1,2,2,1;3,1,2,1;4,3,2,1;1,2,1,3;5,1,3;3,2,1,3;1,4,5;2,1,1,6;2,2,1,3;1,1,4;2,1,7;1,2,1,1,3;12,1;1,3,2,1",
			wantSolved: true,
			wantGrid: strings.TrimSpace(`
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|• • • • • • • o • • o • • o • • o • o •|
|• • • • o o o o o • o • o o o • o o o •|
|• • • o • • • o o o o o • • o o • • o o|
|• • • o • o • • o • o o o o • • o • o •|
|• • • • o o • • o o o o o • • • • o o o|
|• • o o o • o o • o • • o • o • • o o o|
|• o • • • o • • o • • o o o o • o • o o|
|• o • o • • o • o • • o • • • o o • o •|
|• • o o • • o o o o o • • • • • o o o •|
|• o o o o o • o • • • o • o • • o • o o|
|o • • • o • o • • o • • o o o • o o o o|
|o • o • • o o • o o • o o o • o o • o •|
|• o o • • o • o o o o o o o o o o o • •|
|• • o • o • o • • • o o o o o o • o o •|
|• • • o o o • • o • o • o o o o • o • o|
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewFromString(tt.puzzle)
			n.Debug = true
			gotSolved := n.Solve()
			gotGrid := n.GridString()
			if got := n.Solve(); gotSolved != tt.wantSolved {
				t.Errorf("Nonogram.Solve() = %v, want %v", got, tt.wantSolved)
			}
			if gotSolved {
				t.Logf("Solved Nonogram use %d steps: \n%s\n", n.Step, gotGrid)
			}
			if gotGrid != tt.wantGrid {
				t.Errorf("Nonogram.GridString() =\n%v\nwant\n%v", gotGrid, tt.wantGrid)
			}
		})
	}
}
