package dayNine

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func partOne(puzzleFile string) {
	slog.Info("Day Nine part one", "puzzle file", puzzleFile)
}

func partTwo(puzzleFile string) {
	slog.Info("Day Nine part two", "puzzle file", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "dayNine",
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
