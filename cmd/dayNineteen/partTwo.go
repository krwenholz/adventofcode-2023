package dayNineteen

import (
	"fmt"
	"log/slog"
	"strings"
)

/*
*
Okay, so how do I find all of the acceptable compbinations?
  - find all A's
  - if A is a default rule
    then all previous rules are constraints
    find parent(s) of this workflow, repeat
    end at "in"
  - if A is not a default rule
    then all previous rules are constraints
    find parent(s) of this workflow, repeat
    end at "in"

OR!
- for every workflow
	for every rule
		if it is an A state, add a valid range that leads to A
		if it leads directly to a state with a range, add that range with new constraints
- repeat until every "in" has a ValidRange attached to every rule

*/

type ValidRange struct {
	MinX int
	MaxX int
	MinM int
	MaxM int
	MinA int
	MaxA int
	MinS int
	MaxS int
}

func (v *ValidRange) String() string {
	return fmt.Sprintf("ValidRange: x(%d, %d), m(%d, %d), a(%d, %d), s(%d, %d)", v.MinX, v.MaxX, v.MinM, v.MaxM, v.MinA, v.MaxA, v.MinS, v.MaxS)
}

func DefaultValidRange() *ValidRange {
	return &ValidRange{
		MinX: -1,
		MaxX: -1,
		MinM: -1,
		MaxM: -1,
		MinA: -1,
		MaxA: -1,
		MinS: -1,
		MaxS: -1,
	}
}

func min(a, b int) int {
	if a == -1 {
		return b
	}
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a == -1 {
		return b
	}
	if a > b {
		return a
	}
	return b
}

func (v *ValidRange) ConstrainWithRule(r *Rule) *ValidRange {
	newV := *v
	switch r.Category {
	case "x":
		if r.Comparator == ">" {
			newV.MinX = max(newV.MinX, r.Value) + 1
		} else {
			newV.MaxX = min(newV.MaxX, r.Value) - 1
		}
	case "m":
		if r.Comparator == ">" {
			newV.MinM = max(newV.MinM, r.Value) + 1
		} else {
			newV.MaxM = min(newV.MaxM, r.Value) - 1
		}
	case "a":
		if r.Comparator == ">" {
			newV.MinA = max(newV.MinA, r.Value) + 1
		} else {
			newV.MaxA = min(newV.MaxA, r.Value) - 1
		}
	case "s":
		if r.Comparator == ">" {
			newV.MinS = max(newV.MinS, r.Value) + 1
		} else {
			newV.MaxS = min(newV.MaxS, r.Value) - 1
		}
	}

	return &newV
}

func (v *ValidRange) EvadeRule(r *Rule) *ValidRange {
	newV := *v
	switch r.Category {
	case "x":
		if r.Comparator == "<" {
			newV.MinX = max(newV.MinX, r.Value)
		} else {
			newV.MaxX = min(newV.MaxX, r.Value)
		}
	case "m":
		if r.Comparator == "<" {
			newV.MinM = max(newV.MinM, r.Value)
		} else {
			newV.MaxM = min(newV.MaxM, r.Value)
		}
	case "a":
		if r.Comparator == "<" {
			newV.MinA = max(newV.MinA, r.Value)
		} else {
			newV.MaxA = min(newV.MaxA, r.Value)
		}
	case "s":
		if r.Comparator == "<" {
			newV.MinS = max(newV.MinS, r.Value)
		} else {
			newV.MaxS = min(newV.MaxS, r.Value)
		}
	}

	return &newV
}

func (vr *ValidRange) AllCombinations() int {
	// The +1 ensures we _include_ our bounds
	return ((vr.MaxX - vr.MinX + 1) *
		(vr.MaxM - vr.MinM + 1) *
		(vr.MaxA - vr.MinA + 1) *
		(vr.MaxS - vr.MinS + 1))
}

func (vr *ValidRange) Finalize() {
	vr.MinX = max(vr.MinX, 1)
	vr.MaxX = min(vr.MaxX, 4000)
	vr.MinM = max(vr.MinM, 1)
	vr.MaxM = min(vr.MaxM, 4000)
	vr.MinA = max(vr.MinA, 1)
	vr.MaxA = min(vr.MaxA, 4000)
	vr.MinS = max(vr.MinS, 1)
	vr.MaxS = min(vr.MaxS, 4000)
}

func (f *Flower) FindCombinations() int {
	finalRanges := []*ValidRange{}
	validRanges := map[string][]*ValidRange{"in": {DefaultValidRange()}}
	workflows := []string{"in"}
	for len(workflows) > 0 {
		next := f.mappedWorkflows[workflows[0]]
		workflows = workflows[1:]
		curRanges := validRanges[next.Name]
		for _, r := range f.mappedWorkflows[next.Name].Rules {
			if _, ok := validRanges[r.Destination]; !ok {
				validRanges[r.Destination] = []*ValidRange{}
			}

			nextCurRanges := []*ValidRange{}
			for _, vr := range curRanges {
				constrained := vr.ConstrainWithRule(r)
				switch r.Destination {
				case "R":
					// do nothing
				case "A":
					constrained.Finalize()
					finalRanges = append(finalRanges, constrained)
					slog.Debug("rule based combination", "w", next.Name, "r", r, "valid range", constrained)
				default:
					validRanges[r.Destination] = append(validRanges[r.Destination], constrained)
				}

				evaded := vr.EvadeRule(r)
				nextCurRanges = append(nextCurRanges, evaded)
			}
			curRanges = nextCurRanges

			if r.Destination != "A" && r.Destination != "R" {
				workflows = append(workflows, r.Destination)
			}
		}

		// handle default rule, no bifurcating
		for _, vr := range curRanges {
			switch next.DefaultRule {
			case "R":
				// do nothing
			case "A":
				vr.Finalize()
				finalRanges = append(finalRanges, vr)
				slog.Debug("default rule combination", "w", next.Name, "valid range", vr)
			default:
				validRanges[next.DefaultRule] = append(validRanges[next.DefaultRule], vr)
			}
		}

		if next.DefaultRule != "A" && next.DefaultRule != "R" {
			workflows = append(workflows, next.DefaultRule)
		}
	}

	slog.Debug("Final valid ranges", "ranges", finalRanges)
	combinations := 0
	for _, vr := range finalRanges {
		combinations += vr.AllCombinations()
	}
	return combinations
}

func partTwo(puzzleFile string) {
	slog.Info("Day Nineteen part two", "puzzle file", puzzleFile)

	flower, _ := ParseCommands(puzzleFile)

	cs := flower.FindCombinations()

	slog.Info("Found combinations", "combinations", cs)
	if strings.Contains(puzzleFile, "sample") {
		slog.Debug("Expected combinations", "expected", 167409079868000, "diff", 167409079868000-cs)
	}
}
