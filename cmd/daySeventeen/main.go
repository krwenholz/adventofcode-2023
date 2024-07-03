package daySeventeen

import (
	"adventofcode/cmd/fileReader"
	"container/heap"
	"fmt"
	"log/slog"
	"math"
	"os"
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
	return fmt.Sprintf("(%d,%d)", c.Row, c.Col)
}

func (c *Coordinate) Equals(other *Coordinate) bool {
	return c.Row == other.Row && c.Col == other.Col
}

func (c *Coordinate) Move(dir *Direction) *Coordinate {
	return &Coordinate{c.Row + dir.Row, c.Col + dir.Col}
}

type Direction struct {
	Row int
	Col int
}

func (d *Direction) String() string {
	if d.Row == 0 && d.Col == 1 {
		return ">"
	} else if d.Row == 0 && d.Col == -1 {
		return "<"
	} else if d.Row == 1 && d.Col == 0 {
		return "v"
	} else if d.Row == -1 && d.Col == 0 {
		return "^"
	}
	panic("shit, invalid direction")
}

func (d *Direction) Equals(other *Direction) bool {
	return d.Row == other.Row && d.Col == other.Col
}

func (d *Direction) IsPerpendicular(other *Direction) bool {
	if (d.Row == 0 && d.Col == 1) || (d.Row == 0 && d.Col == -1) {
		if (other.Row == 1 && other.Col == 0) || (other.Row == -1 && other.Col == 0) {
			return true
		}
	} else if (d.Row == 1 && d.Col == 0) || (d.Row == -1 && d.Col == 0) {
		if (other.Row == 0 && other.Col == 1) || (other.Row == 0 && other.Col == -1) {
			return true
		}
	}
	return false
}

type Cell struct {
	coords *Coordinate
	dir    *Direction
	steps  int
	f      int // Total cost of the cell (g + h)
	g      int // Cost from start to this cell
	h      int // Heuristic cost from this cell to destination
}

func (c *Cell) String() string {
	return fmt.Sprintf("Cell{%s, dir: %s, steps: %d, f: %d, g: %d, h: %d}", c.coords.String(), c.dir, c.steps, c.f, c.g, c.h)
}

func (c *Cell) Next(dir *Direction, dest *Coordinate, grid []string) (*Cell, error) {
	var steps int
	if c.dir != nil && c.steps == 3 && !c.dir.IsPerpendicular(dir) {
		return nil, fmt.Errorf("too many steps in the same direction without a 90 degree turn")
	}

	if c.dir != nil && c.dir.Equals(dir) {
		steps = c.steps + 1
	} else {
		steps = 1
	}

	newCell := &Cell{
		c.coords.Move(dir),
		dir,
		steps,
		0,
		0,
		0,
	}

	if !(newCell.coords.Row >= 0 && newCell.coords.Row < len(grid) && newCell.coords.Col >= 0 && newCell.coords.Col < len(grid[0])) {
		return nil, fmt.Errorf("invalid cell")
	}

	// Manhattan distance is our h estimate
	newCell.h = (dest.Row - newCell.coords.Row) + (dest.Col - newCell.coords.Col)
	newCell.g = c.g + HeatLoss(newCell.coords.Row, newCell.coords.Col, grid)
	newCell.f = newCell.g + newCell.h

	return newCell, nil
}

func (c *Cell) CellState() string {
	return fmt.Sprintf("Cell{%s, %s, %d}", c.coords.String(), c.dir, c.steps)
}

func ReconstructPath(cameFrom map[string]*Cell, current *Cell) [][]int {
	totalPath := [][]int{{current.coords.Row, current.coords.Col}}
	for {
		if prev, ok := cameFrom[current.CellState()]; ok {
			totalPath = append(totalPath, []int{current.coords.Row, current.coords.Col})
			current = prev
		} else {
			break
		}
	}

	// Reverse the path to get the path from source to destination
	slices.Reverse(totalPath)

	return totalPath
}

func HeatLoss(row, col int, grid []string) int {
	h, _ := strconv.Atoi(string(grid[row][col]))
	return h
}

func directions() []*Direction {
	return []*Direction{
		{0, 1},  // right
		{0, -1}, // left
		{1, 0},  // "down"
		{-1, 0}, // "up"
	}
}

// An A* implementation!
func AStarSearch(grid []string, src, dest *Coordinate) ([][]int, [][]int) {
	// Initialize the closed list (visited cells)
	seen := map[string]*Cell{}
	// Track the best paths
	cameFrom := map[string]*Cell{}

	gScore := make([][]int, len(grid))
	for i := 0; i < len(grid); i++ {
		gScore[i] = make([]int, len(grid[0]))
		for j := 0; j < len(grid[0]); j++ {
			gScore[i][j] = int(math.Inf(1))
		}
	}
	gScore[src.Row][src.Col] = 0

	// Initialize the start cell details
	start := &Cell{src, nil, 0, 0, 0, 0}
	// Initialize the open list (cells to be visited) with the start cell
	openSet := &CellHeap{start}

	// Main loop of A* search algorithm
	for len(*openSet) > 0 {
		current := heap.Pop(openSet).(*Cell)

		slog.Debug("popped!", "cell", current, "open list len", len(*openSet))

		if current.coords.Equals(dest) {
			return ReconstructPath(cameFrom, current), gScore
		}
		// For each direction, check the successors
		for _, dir := range directions() {
			// Cell{(0,11), f: 58, g: 45, h: 13, parents: [(0,10) (1,10) (2,10) (2,11) (0,12)]
			neighbor, err := current.Next(dir, dest, grid)
			if err != nil {
				slog.Debug("invalid state", "cell", current, "dir", dir, "error", err)
				continue
			}

			if neighbor.g < gScore[neighbor.coords.Row][neighbor.coords.Col] {
				// record our visual g scores
				gScore[neighbor.coords.Row][neighbor.coords.Col] = neighbor.g
			}

			if prev, ok := seen[neighbor.CellState()]; !ok || neighbor.g < prev.g {
				// This path to neighbor is better than any previous one. Record it!
				slog.Debug("found new best path!", "cell", neighbor, "open list", len(*openSet))
				cameFrom[neighbor.CellState()] = current
				seen[neighbor.CellState()] = neighbor
				heap.Push(openSet, neighbor)
			}
		}
	}

	PrintCellDetails(gScore)
	panic("Did not find the destination cell")
}

func PrintPath(path [][]int, grid []string) {
	for _, coord := range path {
		row := coord[0]
		col := coord[1]
		grid[row] = grid[row][:col] + "X" + grid[row][col+1:]
	}
	os.WriteFile("/tmp/grid.txt", []byte(strings.Join(grid, "\n")), 0644)
}

func PrintCellDetails(cellDetails [][]int) {
	out := ""
	for _, row := range cellDetails {
		for _, cell := range row {
			if cell == int(math.Inf(1)) {
				out += "---- "
				continue
			}

			out += fmt.Sprintf("%04d ", cell)
		}
		out += "\n"
	}
	os.WriteFile("/tmp/cells.txt", []byte(out), 0644)
}

func partOne(puzzleFile string) {
	slog.Info("Day Seventeen part one", "puzzle file", puzzleFile)
	rows := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")

	src := &Coordinate{0, 0}
	dest := &Coordinate{len(rows) - 1, len(rows[0]) - 1}

	path, gScore := AStarSearch(rows, src, dest)

	PrintPath(path, rows)
	PrintCellDetails(gScore)

	slog.Info("The path from source to destination found", "path", path, "heat loss", gScore[dest.Row][dest.Col])
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
