package dayFive

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
	Use: "dayFive",
	Run: func(cmd *cobra.Command, args []string) {
		puzzleInput, _ := cmd.Flags().GetString("puzzle-input")
		if !cmd.Flag("part-two").Changed {
			partOne(puzzleInput)
		} else {
			partTwo(puzzleInput)
		}
	},
}
