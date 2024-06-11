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
	Conditions     []string `@("#" | "." | "?")+`
	GroupSizes     []int    `(@Int ","?)+`
	TrueConditions []Condition
	conditionsI    int
	sizesI         int
	truesI         int
	damagedCount   int
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

func (c *ConditionRecord) GenerateReplacements() []*ConditionRecord {
	ret := []*ConditionRecord{c}
	for i := 0; i < c.UnknownCount(); i++ {
		newRs := make([]*ConditionRecord, 0)
		for _, cr := range ret {
			for _, cond := range []Condition{Operational, Damaged} {
				newR := &ConditionRecord{
					Conditions:     cr.Conditions,
					GroupSizes:     cr.GroupSizes,
					TrueConditions: cr.TrueConditions,
					conditionsI:    cr.conditionsI,
					sizesI:         cr.sizesI,
					truesI:         cr.truesI,
					damagedCount:   cr.damagedCount,
				}
				newR.TrueConditions = append(newR.TrueConditions, cond)
				if newR.IsValid() {
					slog.Debug("found replacement", "replacement", newR, "sizes", newR.StringSizes())
					newRs = append(newRs, newR)
				}
			}
		}
		ret = newRs
	}
	return ret
}

func (c *ConditionRecord) IsValid() bool {
	for c.conditionsI < len(c.Conditions) {
		cond := NewCondition(rune(c.Conditions[c.conditionsI][0]))
		if cond == Unknown {
			if c.truesI >= len(c.TrueConditions) {
				// we need more replacement options
				return c.UnknownCount() != len(c.TrueConditions)
			}

			cond = c.TrueConditions[c.truesI]
			c.truesI++
		}

		switch cond {
		case Operational:
			if c.damagedCount > 0 {
				if c.sizesI >= len(c.GroupSizes) {
					return false
				}
				if c.damagedCount != c.GroupSizes[c.sizesI] {
					return false
				}
				c.damagedCount = 0
				c.sizesI++
			}
		case Damaged:
			c.damagedCount++
		default:
			panic("Invalid condition")
		}
		c.conditionsI++
	}

	if c.damagedCount > 0 {
		if c.sizesI >= len(c.GroupSizes) {
			return false
		}
		if c.damagedCount != c.GroupSizes[c.sizesI] {
			return false
		}
		if c.conditionsI == len(c.Conditions) {
			c.sizesI++
		}
	}

	if c.conditionsI == len(c.Conditions) && c.sizesI < len(c.GroupSizes) {
		return false
	}

	return true
}

func (c ConditionRecord) String() string {
	return fmt.Sprintf(
		"Conditions: %v, GroupSizes: %v, TrueConditions: %v",
		c.Conditions,
		c.GroupSizes,
		c.TrueConditions,
	)
}

func (c ConditionRecord) StringSizes() string {
	return fmt.Sprintf(
		"conditionsI: %d, sizesI: %d, truesI: %d, damagedCount: %d",
		c.conditionsI,
		c.sizesI,
		c.truesI,
		c.damagedCount,
	)
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
