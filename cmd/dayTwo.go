package cmd

import (
	"adventofcode/cmd/dayTwo"

	"github.com/spf13/cobra"
)

// dayTwoCmd represents the dayTwo command
var dayTwoCmd = &cobra.Command{
	Use: "dayTwo",
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flag("part-two").Changed {
			dayTwo.PartOne(puzzleInput)
		} else {
			dayTwo.PartTwo(puzzleInput)
		}
	},
}

func init() {
	rootCmd.AddCommand(dayTwoCmd)
	dayTwoCmd.Flags().BoolP("part-two", "p", false, "Whether to run part two of the day's challenge")
}
