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

func (c *Cell) Next(dir *Direction, dest *Coordinate, grid [][]int) (*Cell, error) {
	var steps int

	if c.dir != nil && c.dir.Equals(dir) {
		steps = c.steps + 1
	} else {
		steps = 1
	}

	newCoords := c.coords.Move(dir)
	if !(newCoords.Row >= 0 && newCoords.Row < len(grid) && newCoords.Col >= 0 && newCoords.Col < len(grid[0])) {
		return nil, fmt.Errorf("invalid cell")
	}

	// Manhattan distance is our h estimate
	h := (dest.Row - newCoords.Row) + (dest.Col - newCoords.Col)
	g := c.g + grid[newCoords.Row][newCoords.Col]

	return &Cell{
		newCoords,
		dir,
		steps,
		g + h,
		g,
		h,
	}, nil
}

func (c *Cell) CellState() string {
	return fmt.Sprintf("Cell{%s, %s, %d}", c.coords.String(), c.dir, c.steps)
}

func ReconstructPath(cameFrom map[string]*Cell, current *Cell) []*Cell {
	path := []*Cell{current}
	for {
		if prev, ok := cameFrom[current.CellState()]; ok {
			path = append(path, prev)
			current = prev
		} else {
			break
		}
	}

	// Reverse the path to get the path from source to destination
	slices.Reverse(path)

	return path
}

func HeatLoss(row, col int, grid []string) int {
	h, _ := strconv.Atoi(string(grid[row][col]))
	return h
}

// An A* implementation!
func AStarSearch(grid [][]int,
	src, dest *Coordinate,
	dirsFunc func(c *Cell) []*Direction,
	finished func(c *Cell, d *Coordinate) bool,
) ([]*Cell, int, [][]int) {
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

		if finished(current, dest) {
			return ReconstructPath(cameFrom, current), current.f, gScore
		}

		slog.Debug("popped!", "cell", current, "open list len", len(*openSet))

		// For each direction, check the successors
		for _, dir := range dirsFunc(current) {
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

	panic("Did not find the destination cell")
}

func PrintPath(path []*Cell, grid []string) {
	for _, c := range path {
		row := c.coords.Row
		col := c.coords.Col
		dir := "X"
		if c.dir != nil {
			dir = c.dir.String()
		}
		grid[row] = grid[row][:col] + dir + grid[row][col+1:]
	}
	os.WriteFile("/tmp/grid.txt", []byte(strings.Join(grid, "\n")), 0644)
	niceCellPath := []string{}
	for _, c := range path {
		niceCellPath = append(niceCellPath, c.CellState())
	}
	os.WriteFile("/tmp/path.txt", []byte(strings.Join(niceCellPath, "\n")), 0644)
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

func Directions(c *Cell) []*Direction {
	if c.dir == nil {
		// Start point, only left and "down" are valid
		return []*Direction{
			{0, 1},
			{1, 0},
		}
	}
	dirs := []*Direction{
		// Nifty transform for turns
		{-c.dir.Col, c.dir.Row},
		{c.dir.Col, -c.dir.Row},
	}
	if c.steps < 3 {
		dirs = append(dirs, c.dir)
	}
	return dirs
}

func partOne(puzzleFile string) {
	slog.Info("Day Seventeen part one", "puzzle file", puzzleFile)
	rows := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")
	grid := make([][]int, len(rows))
	for i, row := range rows {
		grid[i] = make([]int, len(row))
		for j, cell := range row {
			h, _ := strconv.Atoi(string(cell))
			grid[i][j] = h
		}
	}

	src := &Coordinate{0, 0}
	dest := &Coordinate{len(rows) - 1, len(rows[0]) - 1}

	path, heatLoss, gScore := AStarSearch(grid, src, dest, Directions, func(c *Cell, d *Coordinate) bool {
		return c.coords.Equals(d)
	})

	PrintPath(path, rows)
	PrintCellDetails(gScore)

	slog.Info("The path from source to destination found", "path", path, "heat loss", heatLoss)
}

func UltraDirections(c *Cell) []*Direction {
	if c.dir == nil {
		// Start point, only left and "down" are valid
		return []*Direction{
			{0, 1},
			{1, 0},
		}
	}

	/**
	Once an ultra crucible starts moving in a direction, it needs to move a minimum of four
	blocks in that direction before it can turn (or even before it can stop at the end).
	However, it will eventually start to get wobbly: an ultra crucible can move a maximum of
	ten consecutive blocks without turning.
	**/
	if c.steps < 4 {
		return []*Direction{c.dir}
	}

	dirs := []*Direction{
		// Nifty transform for turns
		{-c.dir.Col, c.dir.Row},
		{c.dir.Col, -c.dir.Row},
	}
	if c.steps < 10 {
		dirs = append(dirs, c.dir)
	}
	return dirs
}

func UltraFinished(c *Cell, d *Coordinate) bool {
	return c.coords.Equals(d) && c.steps >= 4
}

func partTwo(puzzleFile string) {
	slog.Info("Day Seventeen part two", "puzzle file", puzzleFile)
	rows := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")
	grid := make([][]int, len(rows))
	for i, row := range rows {
		grid[i] = make([]int, len(row))
		for j, cell := range row {
			h, _ := strconv.Atoi(string(cell))
			grid[i][j] = h
		}
	}

	src := &Coordinate{0, 0}
	dest := &Coordinate{len(rows) - 1, len(rows[0]) - 1}

	path, heatLoss, gScore := AStarSearch(grid, src, dest, UltraDirections, UltraFinished)

	PrintPath(path, rows)
	PrintCellDetails(gScore)

	slog.Info("The path from source to destination found", "path", path, "heat loss", heatLoss)
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
