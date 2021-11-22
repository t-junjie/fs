// package storage provides methods to manipulate csv file content
// and perform calculations.
package storage

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/google/uuid"
)

type FileInfo struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Size     string    `json:"size,omitempty"`
	RowCount int       `json:"rowcount,omitempty"`
	ColCount int       `json:"colcount,omitempty"`
}

type Cell struct {
	row int
	col int
}

// ConvertToTable takes a csv file and converts it to a Table.
func ConvertToTable(c io.Reader) ([][]int, error) {
	csvReader := csv.NewReader(c)
	var tbl [][]int

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		var r []int
		for _, col := range record {
			val, err := strconv.Atoi(col)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			r = append(r, val)
		}

		tbl = append(tbl, r)
	}

	return tbl, nil
}

// Size returns the row and column count of the table.
// Assume no empty cells in the csv file i.e. cell(row, col) >= 0 for all cells.
func Size(tbl [][]int) (row, col int) {
	row = len(tbl)
	for _, c := range tbl {
		col = len(c)
		break
	}
	return
}

// // Items returns a list of integers from the cells of the table.
// func Items(cells []Cell) []int {
// 	list := make([]int, 0, len(cells))
// 	for _, c := range cells {
// 		list = append(list, t[c.row][c.col])
// 	}
// 	return list
// }
