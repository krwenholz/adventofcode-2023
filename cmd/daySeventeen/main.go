package daySeventeen

import (
	"adventofcode/cmd/fileReader"
	"container/heap"
	"fmt"
	"log/slog"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type Coordinate struct {
	Row int
	Col int
}

func (c *Coordinate) String() string {
	return fmt.Sprintf("(%d, %d)", c.Row, c.Col)
}

type Cell struct {
	parent      *Coordinate
	f           float64 // Total cost of the cell (g + h)
	g           float64 // Cost from start to this cell
	h           float64 // Heuristic cost from this cell to destination
	prevParents []*Coordinate
}

func (c *Cell) String() string {
	return fmt.Sprintf("Cell{%s, f: %f, g: %f, h: %f, prevParents: %v}", c.parent.String(), c.f, c.g, c.h, c.prevParents)
}

func (c *Cell) Next(row, col int) *Cell {
	newParents := append(c.prevParents, c.parent)
	if len(newParents) > 3 {
		newParents = newParents[1:]
	}
	return &Cell{&Coordinate{row, col}, 0, 0, 0, newParents}
}

func (c *Cell) F() float64 {
	return c.g + c.h
}

func (c *Cell) IsValid(rows []string) bool {
	countRow := map[int]int{}
	countCol := map[int]int{}
	for _, parent := range c.prevParents[1:] {
		if _, ok := countRow[parent.Row]; ok {
			countRow[parent.Row]++
		} else {
			countRow[parent.Row] = 1
		}
		if _, ok := countCol[parent.Col]; ok {
			countCol[parent.Col]++
		} else {
			countCol[parent.Col] = 1
		}
	}
	if len(c.prevParents) == 3 && (len(countRow) == 1 || len(countCol) == 1) {
		return false
	}
	return c.parent.Row >= 0 && c.parent.Row < len(rows) && c.parent.Col >= 0 && c.parent.Col < len(rows[0])
}

func IsDestination(row, col int, dest []int) bool {
	return row == dest[0] && col == dest[1]
}

// Calculate the heuristic value of a cell: Manhattan distance
func HValue(row, col int, dest []int) float64 {
	return float64((dest[0] - row) + (dest[1] - col))
	//return math.Pow((math.Pow(float64(row-dest[0]), 2) + math.Pow(float64(col-dest[1]), 2)), 0.5)
}

func TracePath(cellDetails [][]*Cell, dest []int) [][]int {
	path := [][]int{}
	row := dest[0]
	col := dest[1]

	// Trace the path from destination to source using parent cells
	for !(cellDetails[row][col].parent.Row == row && cellDetails[row][col].parent.Col == col) {
		path = append(path, []int{row, col})
		tempRow := cellDetails[row][col].parent.Row
		tempCol := cellDetails[row][col].parent.Col
		row = tempRow
		col = tempCol
		slog.Debug("tracing path", "row", row, "col", col)
	}

	// Add the source cell to the path
	path = append(path, []int{row, col})
	// Reverse the path to get the path from source to destination
	slices.Reverse(path)

	return path
}

func HeatLoss(row, col int, grid []string) float64 {
	h, _ := strconv.Atoi(string(grid[row][col]))
	return float64(h)
}

func directions() [][]int {
	return [][]int{
		{0, 1},  // right
		{0, -1}, // left
		{1, 0},  // "down"
		{-1, 0}, // "up"
	}
}

// An A* implementation!
func AStarSearch(grid []string, src, dest []int) [][]int {
	// Initialize the closed list (visited cells)
	closedList := make([][]bool, len(grid))
	for i := 0; i < len(grid); i++ {
		closedList[i] = make([]bool, len(grid[0]))
	}

	// Initialize the details of each cell
	cellDetails := make([][]*Cell, len(grid))
	for i := 0; i < len(grid); i++ {
		cellDetails[i] = make([]*Cell, len(grid[0]))
		for j := 0; j < len(grid[0]); j++ {
			cellDetails[i][j] = &Cell{
				&Coordinate{-1, -1},
				float64(math.Inf(1)),
				float64(math.Inf(1)),
				float64(math.Inf(1)),
				[]*Coordinate{},
			}
		}
	}

	// Initialize the start cell details
	row := src[0]
	col := src[1]
	start := &Cell{&Coordinate{row, col}, 0, 0, 0, []*Coordinate{}}
	cellDetails[row][col] = start

	// Initialize the open list (cells to be visited) with the start cell
	openList := CellHeap{start}

	// Main loop of A* search algorithm
	for len(openList) > 0 {
		p := heap.Pop(&openList).(*Cell)
		row := p.parent.Row
		col := p.parent.Col
		closedList[row][col] = true

		// For each direction, check the successors
		for _, dir := range directions() {
			newRow := row + dir[0]
			newCol := col + dir[1]

			// If the successor is valid and not visited
			newCell := p.Next(newRow, newCol)
			if newCell.IsValid(grid) && !closedList[newRow][newCol] {
				// If the successor is the destination
				if IsDestination(newRow, newCol, dest) {
					// Set the parent of the destination cell
					cellDetails[newRow][newCol].parent.Row = row
					cellDetails[newRow][newCol].parent.Col = col
					slog.Info("The destination cell is found")
					// Trace and print the path from source to destination
					return TracePath(cellDetails, dest)
				} else {
					// Calculate the new f, g, and h values
					newCell.g = cellDetails[row][col].g + HeatLoss(newRow, newCol, grid)
					newCell.h = HValue(newRow, newCol, dest)
					slog.Debug("Checking successor", "parent", p, "cell", newCell, "direction", dir)

					// If the cell is not in the open list or the new f value is smaller
					if cellDetails[newRow][newCol].f == float64(math.Inf(1)) || cellDetails[newRow][newCol].f > newCell.F() {
						// Add the cell to the open list
						heap.Push(&openList, newCell)
						// Update the cell details
						cellDetails[newRow][newCol].f = newCell.F()
						cellDetails[newRow][newCol].g = newCell.g
						cellDetails[newRow][newCol].h = newCell.h
						cellDetails[newRow][newCol].parent.Row = row
						cellDetails[newRow][newCol].parent.Col = col
					}
				}
			}
		}
	}

	panic("Did not find the destination cell")
}

func partOne(puzzleFile string) {
	slog.Info("Day Seventeen part one", "puzzle file", puzzleFile)
	rows := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")

	src := []int{0, 0}
	dest := []int{len(rows) - 1, len(rows[0]) - 1}

	path := AStarSearch(rows, src, dest)
	heatLoss := 0
	for _, coord := range path {
		// TODO(kyle): need to account for the actual values
		heatLoss += int(HeatLoss(coord[0], coord[1], rows))
	}

	slog.Info("The path from source to destination found", "path", path, "heat loss", heatLoss)
}

func partTwo(puzzleFile string) {
	slog.Info("Day Seventeen part two", "puzzle file", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "daySeventeen",
	Run: func(cmd *cobra.Command, args []string) {
		puzzleInput, _ := cmd.Flags().GetString("puzzle-input")
		if !cmd.Flag("part-two").Changed {
			partOne(puzzleInput)
		} else {
			partTwo(puzzleInput)
		}
	},
}

func init() {
	Cmd.Flags().Bool("part-two", false, "Whether to run part two of the day's challenge")
}
