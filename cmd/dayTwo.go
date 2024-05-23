package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dayTwoCmd represents the dayTwo command
var dayTwoCmd = &cobra.Command{
	Use: "dayTwo",
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flag("part-two").Changed {
			partOneDayTwo(puzzleInput)
		} else {
			partTwoDayTwo(puzzleInput)
		}
	},
}

func init() {
	rootCmd.AddCommand(dayTwoCmd)
	dayTwoCmd.Flags().BoolP("part-two", "p", false, "Whether to run part two of the day's challenge")
}

func partOneDayTwo(puzzleFile string) {
	fmt.Println("Day template part one", puzzleFile)
}

func partTwoDayTwo(puzzleFile string) {
	fmt.Println("Day template part two", puzzleFile)
}
