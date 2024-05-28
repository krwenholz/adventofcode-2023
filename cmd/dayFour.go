package cmd

import (
	"adventofcode/cmd/dayFour"
	"fmt"

	"github.com/spf13/cobra"
)

// dayFourCmd represents the dayTemplate command
var dayFourCmd = &cobra.Command{
	Use: "dayFour",
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flag("part-two").Changed {
			dayFour.PartOne(puzzleInput)
		} else {
			dayFour.PartTwo(puzzleInput)
		}
	},
}

func init() {
	rootCmd.AddCommand(dayFourCmd)
	dayFourCmd.Flags().BoolP("part-two", "p", false, "Whether to run part two of the day's challenge")
}

func partOneDayFour(puzzleFile string) {
	fmt.Println("Day template part one", puzzleFile)
}

func partTwoDayFour(puzzleFile string) {
	fmt.Println("Day template part two", puzzleFile)
}
