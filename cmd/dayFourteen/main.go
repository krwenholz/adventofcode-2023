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
	colEnds := map[int]int{} // map the last row index for each column
	for i := range rows[0] {
		colEnds[i] = len(rows)
	}

	//colRoundRocks := make(map[int][]int, len(rows[0]))
	loadTotal := 0

	// for every row we're either updating a column end, not, or adding a new rounded rock
	for i, row := range rows {
		for col, rock := range row {
			switch rock {
			case '.':
				// nothing for an empty spot!
			case '#':
				// update the end of the column!
				colEnds[col] = len(rows) - i - 1
			case 'O':
				// roll a rounded rock, calculate the load, and move that column end!
				loadTotal += colEnds[col]
				colEnds[col] = colEnds[col] - 1
			default:
				panic("unexpected rock type, WTF")
			}
		}
	}

	return loadTotal
}

func partOne(puzzleFile string) {
	slog.Info("Day Fourteen part one", "puzzle file", puzzleFile)

	rows := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")

	slog.Info("Day fourteen part one total load", "load", load(rows))
}

func replace(row string, col int, char byte) string {
	return row[:col] + string(char) + row[col+1:]
}

func tiltNorth(col int, rows []string) []string {
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
			rows[rowI] = replace(row, col, '.')
			rows[stopRow] = replace(rows[stopRow], col, 'O')
			stopRow = stopRow + 1
		}
	}
	return rows
}

func tiltSouth(col int, rows []string) []string {
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
			rows[rowI] = replace(row, col, '.')
			rows[stopRow] = replace(rows[stopRow], col, 'O')
			stopRow = stopRow - 1
		}
	}

	return rows
}

func tiltWest(row int, rows []string) []string {
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
			rows[row] = replace(rows[row], col, '.')
			rows[row] = replace(rows[row], stopCol, 'O')
			stopCol = stopCol + 1
		}
	}
	return rows
}

func tiltEast(row int, rows []string) []string {
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
			rows[row] = replace(rows[row], col, '.')
			rows[row] = replace(rows[row], stopCol, 'O')
			stopCol = stopCol - 1
		}
	}
	return rows
}

func spinCycle(rows []string) []string {
	for i := range rows[0] {
		rows = tiltNorth(i, rows)
	}
	if os.Getenv("LOG_TILTS") == "YES" {
		fmt.Print("tilted North\n", strings.Join(rows, "\n"), "\n")
	}
	for i := range rows {
		rows = tiltWest(i, rows)
	}
	if os.Getenv("LOG_TILTS") == "YES" {
		fmt.Print("tilted West\n", strings.Join(rows, "\n"), "\n")
	}
	for i := range rows[0] {
		rows = tiltSouth(i, rows)
	}
	if os.Getenv("LOG_TILTS") == "YES" {
		fmt.Print("tilted South\n", strings.Join(rows, "\n"), "\n")
	}
	for i := range rows {
		rows = tiltEast(i, rows)
	}
	if os.Getenv("LOG_TILTS") == "YES" {
		fmt.Print("tilted East\n", strings.Join(rows, "\n"), "\n")
	}

	if os.Getenv("LOG_CYCLES") == "YES" {
		fmt.Print("cycle\n", strings.Join(rows, "\n"), "\n")
	}

	return rows
}

func partTwo(puzzleFile string, cycles int) {
	slog.Info("Day Fourteen part two", "puzzle file", puzzleFile)

	rows := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")
	for i := 0; i < cycles; i++ {
		spinCycle(rows)
	}

	slog.Info("Day fourteen part two total load", "load", load(rows))
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
