// Code generated by new-day DO NOT EDIT
package cmd

import (
	"adventofcode/cmd/dayEight"
	"adventofcode/cmd/dayEighteen"
	"adventofcode/cmd/dayEleven"
	"adventofcode/cmd/dayFifteen"
	"adventofcode/cmd/dayFive"
	"adventofcode/cmd/dayFour"
	"adventofcode/cmd/dayFourteen"
	"adventofcode/cmd/dayNine"
	"adventofcode/cmd/dayNineteen"
	"adventofcode/cmd/daySeven"
	"adventofcode/cmd/daySeventeen"
	"adventofcode/cmd/daySix"
	"adventofcode/cmd/daySixteen"
	"adventofcode/cmd/dayTen"
	"adventofcode/cmd/dayThirteen"
	"adventofcode/cmd/dayThree"
	"adventofcode/cmd/dayTwelve"
	"adventofcode/cmd/dayTwenty"
	"adventofcode/cmd/dayTwentyFour"
	"adventofcode/cmd/dayTwentyOne"
	"adventofcode/cmd/dayTwentyThree"
	"adventofcode/cmd/dayTwentyTwo"
	"adventofcode/cmd/dayTwo"
	

	"github.com/spf13/cobra"
)

func init() {
	for _, c := range []*cobra.Command{
		dayEight.Cmd,
		dayEighteen.Cmd,
		dayEleven.Cmd,
		dayFifteen.Cmd,
		dayFive.Cmd,
		dayFour.Cmd,
		dayFourteen.Cmd,
		dayNine.Cmd,
		dayNineteen.Cmd,
		daySeven.Cmd,
		daySeventeen.Cmd,
		daySix.Cmd,
		daySixteen.Cmd,
		dayTen.Cmd,
		dayThirteen.Cmd,
		dayThree.Cmd,
		dayTwelve.Cmd,
		dayTwenty.Cmd,
		dayTwentyFour.Cmd,
		dayTwentyOne.Cmd,
		dayTwentyThree.Cmd,
		dayTwentyTwo.Cmd,
		dayTwo.Cmd,
		
	} {
		rootCmd.AddCommand(c)
	}
}
	