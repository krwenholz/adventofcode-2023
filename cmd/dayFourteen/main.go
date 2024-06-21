package dayFourteen

import (
	"adventofcode/cmd/fileReader"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
)

/**
1. roll rocks north, they stop at either the top or the first cube or round rock they hit
2. the load is equal to the row they land in +1 (one-indexed)
3. sum it for all rounded rocks
**/

func partOne(puzzleFile string) {
	slog.Info("Day Fourteen part one", "puzzle file", puzzleFile)

	rows := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")

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

	slog.Info("Day fourteen part one total load", "load", loadTotal)
}

func partTwo(puzzleFile string) {
	slog.Info("Day Fourteen part two", "puzzle file", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "dayFourteen",
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
