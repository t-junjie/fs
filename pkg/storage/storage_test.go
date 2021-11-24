package storage

import (
	"testing"

	"github.com/google/uuid"
)

func TestSize(t *testing.T) {
	type TestResult struct {
		Row int
		Col int
	}

	tests := []struct {
		Name        string
		Tbl         [][]int
		ExpectedVal TestResult
	}{
		{
			Name: "2x2 input returns row = 2, col = 2",
			Tbl: [][]int{
				{0, 1},
				{4, 5},
			},
			ExpectedVal: TestResult{2, 2},
		},
		{
			Name:        "0x0 input returns row = 0, col = 0",
			Tbl:         [][]int{},
			ExpectedVal: TestResult{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			gotRow, gotCol := Size(tt.Tbl)
			actualVal := TestResult{gotRow, gotCol}
			if actualVal != tt.ExpectedVal {
				t.Fail()
			}
		})
	}
}

func TestItems(t *testing.T) {
	tests := []struct {
		Name        string
		Tbl         [][]int
		Cells       []Cell
		ExpectedVal []int
	}{
		{
			Name:        "Choosing (0,0) and (1,1) from [][]int{{0, 1},{4, 5}} returns []int{0,5}",
			Tbl:         [][]int{{0, 1}, {4, 5}},
			Cells:       []Cell{{uuid.New(), 0, 0}, {uuid.New(), 1, 1}}, // uuid is arbitrary
			ExpectedVal: []int{0, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			actualVal := Items(tt.Tbl, tt.Cells)
			actualSum := 0
			for _, val := range actualVal {
				actualSum += val
			}

			expectedSum := 0
			for _, val := range tt.ExpectedVal {
				expectedSum += val
			}

			if actualSum != expectedSum {
				t.Fail()
			}
		})
	}
}
