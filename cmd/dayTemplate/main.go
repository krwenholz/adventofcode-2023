package dayTemplate

import (
	"fmt"

	"github.com/spf13/cobra"
)

func partOne(puzzleFile string) {
	fmt.Println("Day part one", puzzleFile)
}

func partTwo(puzzleFile string) {
	fmt.Println("Day part two", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "dayTemplate",
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
	Cmd.Flags().BoolP("part-two", "p", false, "Whether to run part two of the day's challenge")
}
