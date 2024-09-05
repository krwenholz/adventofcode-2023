package dayTwentyThree

import (
	"adventofcode/cmd/fileReader"
	"adventofcode/cmd/util"
	"container/heap"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var counter = 0

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
	prevs  map[string]bool
	f      int // Total cost of the cell (g + h)
	g      int // Cost from start to this cell
	h      int // Heuristic cost from this cell to destination
}

func (c *Cell) String() string {
	return fmt.Sprintf("Cell{%s, f: %d, g: %d, h: %d, prevs: %v}", c.coords.String(), c.f, c.g, c.h, c.prevs)
}

func (c *Cell) Next(dir *Direction, dest *Coordinate, grid [][]string, validPositions int) (*Cell, error) {
	newCoords := c.coords.Move(dir)
	if !(newCoords.Row >= 0 && newCoords.Row < len(grid) && newCoords.Col >= 0 && newCoords.Col < len(grid[0])) ||
		grid[newCoords.Row][newCoords.Col] == "#" ||
		c.prevs[newCoords.String()] {
		return nil, fmt.Errorf("invalid cell")
	}
	switch grid[c.coords.Row][c.coords.Col] {
	case "^":
		if !dir.Equals(&Direction{-1, 0}) {
			return nil, fmt.Errorf("wrong way on the slope")
		}
	case ">":
		if !dir.Equals(&Direction{0, 1}) {
			return nil, fmt.Errorf("wrong way on the slope")
		}
	case "v":
		if !dir.Equals(&Direction{1, 0}) {
			return nil, fmt.Errorf("wrong way on the slope")
		}
	case "<":
		if !dir.Equals(&Direction{0, -1}) {
			return nil, fmt.Errorf("wrong way on the slope")
		}
	}

	// NOTE: the values here are _negative_, because that's awesome and gets us a longest path
	h := len(c.prevs) - validPositions - (util.Abs(dest.Row-newCoords.Row) + util.Abs(dest.Col-newCoords.Col))
	// g is just distance traveled
	g := c.g - 1

	newPrevs := map[string]bool{}
	for k, v := range c.prevs {
		newPrevs[k] = v
	}
	newPrevs[newCoords.String()] = true

	// Invert all values to look for the _longest_ distance
	return &Cell{
		newCoords,
		newPrevs,
		g + h,
		g,
		h,
	}, nil
}

func (c *Cell) CellState() string {
	return fmt.Sprintf("Cell{%s, %d}", c.coords.String(), c.f)
}

func Directions() []*Direction {
	dirs := []*Direction{
		{1, 0},
		{-1, 0},
		{0, 1},
		{0, -1},
	}
	return dirs
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

func CountValidSpaces(grid [][]string) int {
	count := 0
	for _, row := range grid {
		for _, pos := range row {
			if pos != "#" {
				count++
			}
		}
	}
	return count
}

// An A* implementation!
func AStarSearch(grid [][]string,
	src, dest *Coordinate,
	finished func(c *Cell, d *Coordinate) bool,
) ([]*Cell, *Cell) {
	validPositions := CountValidSpaces(grid)
	// Initialize the closed list (visited cells)
	seen := map[string]*Cell{}
	// Track the best paths
	cameFrom := map[string]*Cell{}

	// Initialize the start cell details
	start := &Cell{src, map[string]bool{}, 0, 0, 0}
	// Initialize the open list (cells to be visited) with the start cell
	openSet := &CellHeap{start}

	// A helpful debug cell set
	dbgCells := map[string][]*Cell{}

	// Main loop of A* search algorithm
	for len(*openSet) > 0 {
		current := heap.Pop(openSet).(*Cell)
		if current.coords.Equals(&Coordinate{6, 3}) || current.coords.Equals(&Coordinate{5, 4}) {
			if dbgCells[current.coords.String()] == nil {
				dbgCells[current.coords.String()] = []*Cell{}
			}
			dbgCells[current.coords.String()] = append(dbgCells[current.coords.String()], current)
		}

		if finished(current, dest) {
			slog.Debug("found the destination!", "cell", current, "open list", *openSet)
			dbgCellsOut := []byte{}
			for k, v := range dbgCells {
				dbgCellsOut = append(dbgCellsOut, []byte(fmt.Sprintf("%s\n", k))...)
				for _, c := range v {
					dbgCellsOut = append(dbgCellsOut, []byte(fmt.Sprintf("\t%s\n", c))...)
				}
				dbgCellsOut = append(dbgCellsOut, []byte("\n")...)
			}
			os.WriteFile("/tmp/dbgCells.txt", dbgCellsOut, 0644)
			return ReconstructPath(cameFrom, current), current
		}

		slog.Debug("popped!", "cell", current, "open list len", len(*openSet))

		// For each direction, check the successors
		for _, dir := range Directions() {
			neighbor, err := current.Next(dir, dest, grid, validPositions)
			if err != nil {
				slog.Debug("invalid state", "cell", current, "dir", dir, "error", err)
				continue
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

func findOnlySlot(grid [][]string, row int) *Coordinate {
	for i, pos := range grid[row] {
		if pos == "." {
			return &Coordinate{row, i}
		}
	}
	return nil
}

func PrintPath(path []*Cell, grid []string) {
	for _, c := range path {
		row := c.coords.Row
		col := c.coords.Col
		dir := "O"
		grid[row] = grid[row][:col] + dir + grid[row][col+1:]
	}
	os.WriteFile("/tmp/grid.txt", []byte(strings.Join(grid, "\n")), 0644)
	niceCellPath := []string{}
	for _, c := range path {
		niceCellPath = append(niceCellPath, c.CellState())
	}
	os.WriteFile("/tmp/path.txt", []byte(strings.Join(niceCellPath, "\n")), 0644)
}

func partOne(puzzleFile string) {
	slog.Info("Day TwentyThree part one", "puzzle file", puzzleFile)

	rows := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")
	expected := rows[0]
	rows = rows[2:]
	grid := make([][]string, len(rows))
	for i, row := range rows {
		grid[i] = make([]string, len(row))
		for j, pos := range row {
			grid[i][j] = string(pos)
		}
	}

	start := findOnlySlot(grid, 0)
	end := findOnlySlot(grid, len(grid)-1)

	path, finalCell := AStarSearch(grid, start, end, func(c *Cell, d *Coordinate) bool {
		return c.coords.Equals(d)
	})

	PrintPath(path, rows)

	slog.Info("Day TwentyThree part one", "expected", expected, "distance", -finalCell.g)
}

func partTwo(puzzleFile string) {
	slog.Info("Day TwentyThree part two", "puzzle file", puzzleFile)
	// TODO: I think I just need to rewrite this to use a DFS with a couple optimizations. See the Reddit thread:
	// https://www.reddit.com/r/adventofcode/comments/18rak3k/2023_rust_solving_everything_under_1_second/

	rows := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")
	expected := rows[1]
	rows = rows[2:]

	grid := make([][]string, len(rows))
	for i, row := range rows {
		grid[i] = make([]string, len(row))
		for j, pos := range row {
			if pos == '.' || pos == '#' {
				grid[i][j] = string(pos)
			} else {
				grid[i][j] = "."
			}
		}
	}

	start := findOnlySlot(grid, 0)
	end := findOnlySlot(grid, len(grid)-1)

	path, finalCell := AStarSearch(grid, start, end, func(c *Cell, d *Coordinate) bool {
		return c.coords.Equals(d)
	})

	PrintPath(path, rows)

	slog.Info("Day TwentyThree part one", "expected", expected, "distance", -finalCell.g)
}

var Cmd = &cobra.Command{
	Use: "dayTwentyThree",
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
