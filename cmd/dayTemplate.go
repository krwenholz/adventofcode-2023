package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dayTemplateCmd represents the dayTemplate command
var dayTemplateCmd = &cobra.Command{
	Use: "dayTemplate",
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flag("part-two").Changed {
			partOneDayTemplate(puzzleInput)
		} else {
			partTwoDayTemplate(puzzleInput)
		}
	},
}

func init() {
	rootCmd.AddCommand(dayTemplateCmd)
	dayTemplateCmd.Flags().BoolP("part-two", "p", false, "Whether to run part two of the day's challenge")
}

func partOneDayTemplate(puzzleFile string) {
	fmt.Println("Day template part one", puzzleFile)
}

func partTwoDayTemplate(puzzleFile string) {
	fmt.Println("Day template part two", puzzleFile)
}
