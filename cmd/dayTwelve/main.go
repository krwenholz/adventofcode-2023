package dayTwelve

import (
	"adventofcode/cmd/scanner"
	"fmt"
	"log"
	"log/slog"

	"github.com/alecthomas/participle/v2"
	"github.com/spf13/cobra"
)

type ConditionRecord struct {
	Conditions []string `@("#" | "." | "?")+`
	GroupSizes []int    `@Int+`
}

func (c ConditionRecord) String() string {
	return fmt.Sprintf("Conditions: %v, GroupSizes: %v", c.Conditions, c.GroupSizes)
}

func newScanner(puzzleFile string) *scanner.PuzzleScanner[ConditionRecord] {
	parser, err := participle.Build[ConditionRecord]()
	if err != nil {
		log.Fatal(err)
	}

	return scanner.NewScanner[ConditionRecord](parser, puzzleFile)
}
func partOne(puzzleFile string) {
	slog.Info("Day Twelve part one", "puzzle file", puzzleFile)
	sc := newScanner(puzzleFile)

	for sc.Scan() {
		r := sc.Struct()
		slog.Debug("parsed record", "record", r)
	}
}

func partTwo(puzzleFile string) {
	slog.Info("Day Twelve part two", "puzzle file", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "dayTwelve",
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
