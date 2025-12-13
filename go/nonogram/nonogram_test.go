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
			name:       "ID P001",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewFromString(tt.puzzle)
			gotSolved := n.Solve()
			gotGrid := n.GridString()
			if got := n.Solve(); gotSolved != tt.wantSolved {
				t.Errorf("Nonogram.Solve() = %v, want %v", got, tt.wantSolved)
			}
			if gotGrid != tt.wantGrid {
				t.Errorf("Nonogram.GridString() =\n%v, want\n%v", gotGrid, tt.wantGrid)
			}
		})
	}
}
