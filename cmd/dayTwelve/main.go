package dayTwelve

import (
	"adventofcode/cmd/scanner"
	"fmt"
	"log"
	"log/slog"

	"github.com/alecthomas/participle/v2"
	"github.com/spf13/cobra"
)

type Condition int

const (
	Operational Condition = iota
	Damaged
	Unknown
)

func (c Condition) String() string {
	switch c {
	case Operational:
		return "."
	case Damaged:
		return "#"
	case Unknown:
		return "?"
	default:
		return ""
	}
}

func NewCondition(r rune) Condition {
	switch r {
	case '.':
		return Operational
	case '#':
		return Damaged
	case '?':
		return Unknown
	default:
		return Unknown
	}
}

type ConditionRecord struct {
	Conditions []string `@("#" | "." | "?")+`
	GroupSizes []int    `(@Int ","?)+`
}

func (c ConditionRecord) UnknownCount() int {
	count := 0
	for _, cond := range c.Conditions {
		if NewCondition(rune(cond[0])) == Unknown {
			count++
		}
	}
	return count

}

func GenerateOptions(unknownCount int) [][]Condition {
	ret := make([][]Condition, 1)
	for i := 0; i < unknownCount; i++ {
		n := make([][]Condition, 0)
		for _, o := range ret {
			slog.Debug("processing option", "option", o)
			newOperational := make([]Condition, len(o)+1)
			newOperational = append(o, Operational)
			newDamaged := append(o, Damaged)
			slog.Debug("will add", "operational", newOperational, "damaged", newDamaged)
			n = append(n, newOperational, newDamaged)
		}
		slog.Debug("generated options", "options", n)
		ret = n
	}
	return ret
}

func IsValid(cs []string, replacements []Condition, sizes []int) bool {
	damagedCount := 0
	replacementsIndex := 0
	sizesIndex := 0
	for _, cond := range cs {
		c := NewCondition(rune(cond[0]))
		if c == Unknown {
			c = replacements[replacementsIndex]
			replacementsIndex++
		}
		switch c {
		case Operational:
			if damagedCount > 0 {
				if damagedCount != sizes[sizesIndex] {
					return false
				}
				damagedCount = 0
				sizesIndex++
			}
		case Damaged:
			damagedCount++
		default:
			panic("Invalid condition")
		}
	}

	if damagedCount > 0 {
		if damagedCount != sizes[sizesIndex] {
			return false
		}
	}

	return true
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

	sumOptions := 0
	for sc.Scan() {
		r := sc.Struct()
		options := GenerateOptions(r.UnknownCount())
		validOptions := make([][]Condition, 0)
		for _, option := range options {
			slog.Debug("checking option", "option", option)
			if IsValid(r.Conditions, option, r.GroupSizes) {
				validOptions = append(validOptions, option)
			}
		}
		sumOptions += len(validOptions)

		slog.Debug("parsed record", "record", r, "total valid", len(validOptions), "valid options", validOptions)
	}

	slog.Info("finished day twelve part one", "sum options", sumOptions)
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
