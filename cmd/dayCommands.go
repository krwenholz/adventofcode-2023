// Code generated by new-day DO NOT EDIT
package cmd

import (
	"adventofcode/cmd/dayEight"
	"adventofcode/cmd/dayEleven"
	"adventofcode/cmd/dayFive"
	"adventofcode/cmd/dayFour"
	"adventofcode/cmd/dayNine"
	"adventofcode/cmd/daySeven"
	"adventofcode/cmd/daySix"
	"adventofcode/cmd/dayTen"
	"adventofcode/cmd/dayThirteen"
	"adventofcode/cmd/dayThree"
	"adventofcode/cmd/dayTwelve"
	"adventofcode/cmd/dayTwo"
	

	"github.com/spf13/cobra"
)

func init() {
	for _, c := range []*cobra.Command{
		dayEight.Cmd,
		dayEleven.Cmd,
		dayFive.Cmd,
		dayFour.Cmd,
		dayNine.Cmd,
		daySeven.Cmd,
		daySix.Cmd,
		dayTen.Cmd,
		dayThirteen.Cmd,
		dayThree.Cmd,
		dayTwelve.Cmd,
		dayTwo.Cmd,
		
	} {
		rootCmd.AddCommand(c)
	}
}
	