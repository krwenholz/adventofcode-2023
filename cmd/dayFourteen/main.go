package dayFourteen

import (
	"adventofcode/cmd/fileReader"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

/**
1. roll rocks north, they stop at either the top or the first cube or round rock they hit
2. the load is equal to the row they land in +1 (one-indexed)
3. sum it for all rounded rocks
**/

func load(rows []string) int {
	//colRoundRocks := make(map[int][]int, len(rows[0]))
	loadTotal := 0

	// for every row we're either updating a column end, not, or adding a new rounded rock
	for i, row := range rows {
		for _, rock := range row {
			if rock == 'O' {
				// roll a rounded rock, calculate the load, and move that column end!
				loadTotal += len(rows) - i
			}
		}
	}

	return loadTotal
}

func partOne(puzzleFile string) {
	slog.Info("Day Fourteen part one", "puzzle file", puzzleFile)

	rows := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")
	g := toGrid(rows)
	for i := range rows[0] {
		g = tiltNorth(i, g)
	}

	slog.Info("Day fourteen part one total load", "load", load(rows))
}

func tiltNorth(col int, rows [][]rune) [][]rune {
	stopRow := 0
	for rowI := range rows {
		row := rows[rowI]
		switch row[col] {
		case '.':
			// keep going
		case '#':
			// new stop
			stopRow = rowI + 1
		case 'O':
			// roll!
			rows[rowI][col] = '.'
			rows[stopRow][col] = 'O'
			stopRow = stopRow + 1
		}
	}
	return rows
}

func tiltSouth(col int, rows [][]rune) [][]rune {
	stopRow := len(rows) - 1
	for rowI := len(rows) - 1; rowI >= 0; rowI-- {
		row := rows[rowI]
		switch row[col] {
		case '.':
			// keep going
		case '#':
			// new stop
			stopRow = rowI - 1
		case 'O':
			// roll!
			rows[rowI][col] = '.'
			rows[stopRow][col] = 'O'
			stopRow = stopRow - 1
		}
	}

	return rows
}

func tiltWest(row int, rows [][]rune) [][]rune {
	stopCol := 0
	for col := range rows[0] {
		switch rows[row][col] {
		case '.':
			// keep going
		case '#':
			// new stop
			stopCol = col + 1
		case 'O':
			// roll!
			rows[row][col] = '.'
			rows[row][stopCol] = 'O'
			stopCol = stopCol + 1
		}
	}
	return rows
}

func tiltEast(row int, rows [][]rune) [][]rune {
	stopCol := len(rows[0]) - 1
	for col := len(rows[0]) - 1; col >= 0; col-- {
		switch rows[row][col] {
		case '.':
			// keep going
		case '#':
			// new stop
			stopCol = col - 1
		case 'O':
			// roll!
			rows[row][col] = '.'
			rows[row][stopCol] = 'O'
			stopCol = stopCol - 1
		}
	}
	return rows
}

func printGrid(dir string, rows [][]rune) {
	fmt.Println("grid after", dir)
	for _, row := range rows {
		fmt.Println(string(row))
	}
}

func spinCycle(rows [][]rune) [][]rune {
	for i := range rows[0] {
		rows = tiltNorth(i, rows)
	}
	if os.Getenv("LOG_TILTS") == "YES" {
		printGrid("North", rows)
	}
	for i := range rows {
		rows = tiltWest(i, rows)
	}
	if os.Getenv("LOG_TILTS") == "YES" {
		printGrid("West", rows)
	}
	for i := range rows[0] {
		rows = tiltSouth(i, rows)
	}
	if os.Getenv("LOG_TILTS") == "YES" {
		printGrid("South", rows)
	}
	for i := range rows {
		rows = tiltEast(i, rows)
	}
	if os.Getenv("LOG_TILTS") == "YES" {
		printGrid("East", rows)
	}

	if os.Getenv("LOG_CYCLES") == "YES" {
		printGrid("Cycle", rows)
	}

	return rows
}

func gridChecksum(rows [][]rune) string {
	checksum := ""
	for _, row := range rows {
		checksum += string(row)
	}
	return checksum
}

func toGrid(rows []string) [][]rune {
	grid := make([][]rune, len(rows))
	for i, row := range rows {
		grid[i] = []rune(row)
	}
	return grid
}

func toRows(grid [][]rune) []string {
	rows := make([]string, len(grid))
	for i, row := range grid {
		rows[i] = string(row)
	}
	return rows
}

func partTwo(puzzleFile string, cycles int) {
	slog.Info("Day Fourteen part two", "puzzle file", puzzleFile)

	rows := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")
	g := toGrid(rows)

	seenGrids := map[string]int{}

	for i := 0; i < cycles; i++ {
		g = spinCycle(g)
		if seen, ok := seenGrids[gridChecksum(g)]; ok {
			slog.Debug("Day fourteen part two repeat found", "cycle", i, "seen", seen)
			maxSkips := (cycles - i) / (i - seen)
			i = i + maxSkips*(i-seen)
			slog.Debug("skipping", "new cycle", i, "maxSkips", maxSkips)
			continue
		}
		seenGrids[gridChecksum(g)] = i
		slog.Debug("Day fourteen part two cycle", "cycle", i, "load", load(toRows(g)))
	}

	printGrid("Final", g)
	slog.Info("Day fourteen part two total load", "load", load(toRows(g)))
}

var Cmd = &cobra.Command{
	Use: "dayFourteen",
	Run: func(cmd *cobra.Command, args []string) {
		puzzleInput, _ := cmd.Flags().GetString("puzzle-input")
		if !cmd.Flag("part-two").Changed {
			partOne(puzzleInput)
		} else {
			cycles, _ := cmd.Flags().GetInt("cycles")
			partTwo(puzzleInput, cycles)
		}
	},
}

func init() {
	Cmd.Flags().Bool("part-two", false, "Whether to run part two of the day's challenge")
	Cmd.Flags().Int("cycles", 1, "Cycles to run, only applicable for part two")
}
