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

type TrueCondition struct {
	TrueConditions []Condition
	conditionsI    int
	sizesI         int
	truesI         int
	damagedCount   int
}

func (c ConditionRecord) String() string {
	return fmt.Sprintf(
		"Conditions: %v, GroupSizes: %v",
		c.Conditions,
		c.GroupSizes,
	)
}

func (t TrueCondition) String() string {
	return fmt.Sprintf(
		"trues: %v, conditionsI: %d, sizesI: %d, truesI: %d, damagedCount: %d",
		t.TrueConditions,
		t.conditionsI,
		t.sizesI,
		t.truesI,
		t.damagedCount,
	)
}

func (c *ConditionRecord) UnknownCount() int {
	count := 0
	for _, cond := range c.Conditions {
		if NewCondition(rune(cond[0])) == Unknown {
			count++
		}
	}
	return count

}

func (c *ConditionRecord) GenerateReplacements() []*TrueCondition {
	ret := []*TrueCondition{{TrueConditions: make([]Condition, 0)}}
	for i := 0; i < c.UnknownCount(); i++ {
		newTCs := make([]*TrueCondition, 0)
		for _, tc := range ret {
			for _, cond := range []Condition{Operational, Damaged} {
				newTC := &TrueCondition{
					TrueConditions: tc.TrueConditions,
					conditionsI:    tc.conditionsI,
					sizesI:         tc.sizesI,
					truesI:         tc.truesI,
					damagedCount:   tc.damagedCount,
				}
				newTC.TrueConditions = append(newTC.TrueConditions, cond)
				if c.IsValid(newTC) {
					slog.Debug("found replacement", "cond", c, "tc", newTC)
					newTCs = append(newTCs, newTC)
				}
			}
		}
		ret = newTCs
	}
	return ret
}

func (c *ConditionRecord) IsValid(tc *TrueCondition) bool {
	for tc.conditionsI < len(c.Conditions) {
		cond := NewCondition(rune(c.Conditions[tc.conditionsI][0]))
		if cond == Unknown {
			if tc.truesI >= len(tc.TrueConditions) {
				// we need more replacement options
				return c.UnknownCount() != len(tc.TrueConditions)
			}

			cond = tc.TrueConditions[tc.truesI]
			tc.truesI++
		}

		switch cond {
		case Operational:
			if tc.damagedCount > 0 {
				if tc.sizesI >= len(c.GroupSizes) {
					return false
				}
				if tc.damagedCount != c.GroupSizes[tc.sizesI] {
					return false
				}
				tc.damagedCount = 0
				tc.sizesI++
			}
		case Damaged:
			tc.damagedCount++
		default:
			panic("Invalid condition")
		}
		tc.conditionsI++
	}

	if tc.damagedCount > 0 {
		if tc.sizesI >= len(c.GroupSizes) {
			return false
		}
		if tc.damagedCount != c.GroupSizes[tc.sizesI] {
			return false
		}
		if tc.conditionsI == len(c.Conditions) {
			tc.sizesI++
		}
	}

	if tc.conditionsI == len(c.Conditions) && tc.sizesI < len(c.GroupSizes) {
		return false
	}

	return true
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
		slog.Debug("Checking Record", "record", r, "unknown count", r.UnknownCount())
		replacementOptions := r.GenerateReplacements()
		sumOptions += len(replacementOptions)

		slog.Debug("parsed record", "record", r, "total valid", len(replacementOptions), "valid options", replacementOptions)
	}

	slog.Info("finished day twelve part one", "sum options", sumOptions)
}

func partTwo(puzzleFile string) {
	slog.Info("Day Twelve part two", "puzzle file", puzzleFile)
	sc := newScanner(puzzleFile)

	sumOptions := 0
	for sc.Scan() {
		r := sc.Struct()
		for i := 0; i < 4; i++ {
			r.Conditions = append(r.Conditions, "?")
			r.Conditions = append(r.Conditions, r.Conditions...)
			r.GroupSizes = append(r.GroupSizes, r.GroupSizes...)
		}
		slog.Debug("Checking Record", "record", r, "unknown count", r.UnknownCount())
		replacementOptions := r.GenerateReplacements()
		sumOptions += len(replacementOptions)
	}

	slog.Info("finished day twelve part two", "sum options", sumOptions)
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
