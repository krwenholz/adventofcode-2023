package cmd

import (
	"adventofcode/cmd/dayThree"

	"github.com/spf13/cobra"
)

// dayThreeCmd represents the dayTwo command
var dayThreeCmd = &cobra.Command{
	Use: "dayThree",
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flag("part-two").Changed {
			dayThree.PartOne(puzzleInput)
		} else {
			dayThree.PartTwo(puzzleInput)
		}
	},
}

func init() {
	rootCmd.AddCommand(dayThreeCmd)
	dayThreeCmd.Flags().BoolP("part-two", "p", false, "Whether to run part two of the day's challenge")
}
