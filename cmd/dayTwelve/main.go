package dayTwelve

import (
	"adventofcode/cmd/scanner"
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"

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

var VALID_CACHE = make(map[string]bool)

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

				var isValid bool
				if cacheVal, ok := VALID_CACHE[newTC.String()]; ok {
					isValid = cacheVal
				} else {
					isValid = c.IsValid(newTC)
					VALID_CACHE[newTC.String()] = isValid
				}

				if isValid {
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
		slog.Info("parsed record")
	}

	slog.Info("finished day twelve part one", "sum options", sumOptions)
}

func partTwo(puzzleFile string) {
	slog.Info("Day Twelve part two", "puzzle file", puzzleFile)
	sc := newScanner(puzzleFile)

	sumOptions := 0
	for sc.Scan() {
		r := sc.Struct()
		newCons := r.Conditions
		newGs := r.GroupSizes
		for i := 0; i < 4; i++ {
			newCons = append(newCons, "?")
			newCons = append(newCons, r.Conditions...)
			newGs = append(newGs, r.GroupSizes...)
		}
		r.Conditions = newCons
		r.GroupSizes = newGs
		slog.Debug("Checking Record", "record", r, "unknown count", r.UnknownCount())
		replacementOptions := r.GenerateReplacements()
		sumOptions += len(replacementOptions)
	}

	slog.Info("finished day twelve part two", "sum options", sumOptions)
}

// Borrowed! https://github.com/Rtchaik/AoC-2023/blob/main/Day12/solution.py
type Data struct {
	Row        string
	GroupSizes []int
}

func withReddit(puzzleFile string) {
	slog.Info("Day twelve part one with Reddit", "f", puzzleFile)
	f, err := os.Open(puzzleFile)
	if err != nil {
		log.Fatal(err)
	}

	sc := bufio.NewScanner(f)

	data := []*Data{}
	for sc.Scan() {
		d := &Data{
			Row:        "",
			GroupSizes: []int{},
		}

		t := sc.Text()
		d.Row = strings.Split(t, " ")[0]
		for _, c := range strings.Split(strings.Split(t, " ")[1], ",") {
			i, _ := strconv.Atoi(c)
			d.GroupSizes = append(d.GroupSizes, i)
		}

		data = append(data, d)
	}

	sum := 0
	for _, d := range data {
		sum += springsFinder(d.Row+".", d.GroupSizes)
		slog.Debug("processed row", "data", d, "sum", sum)
	}
	slog.Info("finished day twelve part one with Reddit", "sum", sum)

	sum = 0
	for _, d := range data {
		expandedD := &Data{
			Row:        d.Row,
			GroupSizes: d.GroupSizes,
		}
		for i := 0; i < 4; i++ {
			expandedD.Row += "?" + d.Row
			expandedD.GroupSizes = append(expandedD.GroupSizes, d.GroupSizes...)
		}
		sum += springsFinder(expandedD.Row+".", expandedD.GroupSizes)
	}
	slog.Info("finished day twelve part two with Reddit", "sum", sum)
}

var SPRINGS_CACHE = make(map[string]int)

func springsCacheKey(row string, nums []int) string {
	return fmt.Sprintf("%s-%v", row, nums)
}

func sum(nums []int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

func springsFinder(row string, nums []int) int {
	slog.Debug("finding springs", "row", row, "nums", nums)
	if v, ok := SPRINGS_CACHE[springsCacheKey(row, nums)]; ok {
		return v
	}

	next := nums[1:]
	springs := []string{}
	for spr := 0; spr < len(row)-sum(nums)-len(next); spr++ {
		expandedSprings := ""
		for i := 0; i < spr; i++ {
			expandedSprings += "."
		}
		expandedBreaks := ""
		for i := 0; i < nums[0]; i++ {
			expandedBreaks += "#"
		}
		springs = append(springs, fmt.Sprintf("%s%s.", expandedSprings, expandedBreaks))
	}

	valid := []int{}
	for _, s := range springs {
		if !isValid(row, s) {
			continue
		}

		valid = append(valid, len(s))
	}

	slog.Debug("found valid springs", "springs", springs, "valid", valid)
	sum := 0
	if len(next) == 0 {
		for _, v := range valid {
			foundPound := false
			for i := v; i < len(row); i++ {
				if row[i] == '#' {
					foundPound = true
					break
				}
			}
			if !foundPound {
				sum++
			}
		}
	} else {
		for _, v := range valid {
			sum += springsFinder(row[v:], next)
		}
	}

	SPRINGS_CACHE[springsCacheKey(row, nums)] = sum
	return sum
}

func isValid(row string, spr string) bool {
	for i := 0; i < len(row) && i < len(spr); i++ {
		if row[i] == '?' {
			continue
		}
		if row[i] != spr[i] {
			return false
		}
	}
	return true
}

var Cmd = &cobra.Command{
	Use: "dayTwelve",
	Run: func(cmd *cobra.Command, args []string) {
		puzzleInput, _ := cmd.Flags().GetString("puzzle-input")
		if cmd.Flag("with-reddit").Changed {
			withReddit(puzzleInput)
			return
		}
		if !cmd.Flag("part-two").Changed {
			partOne(puzzleInput)
		} else {
			partTwo(puzzleInput)
		}
	},
}

func init() {
	Cmd.Flags().Bool("part-two", false, "Whether to run part two of the day's challenge")
	Cmd.Flags().Bool("with-reddit", false, "Use the suggested solution from Reddit")
}
