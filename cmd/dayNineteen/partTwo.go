package dayNineteen

import (
	"fmt"
	"log/slog"
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
	switch r.Category {
	case "x":
		if r.Comparator == ">" {
			v.MinX = min(v.MinX, r.Value)
		} else {
			v.MaxX = max(v.MaxX, r.Value)
		}
	case "m":
		if r.Comparator == ">" {
			v.MinM = min(v.MinM, r.Value)
		} else {
			v.MaxM = max(v.MaxM, r.Value)
		}
	case "a":
		if r.Comparator == ">" {
			v.MinA = min(v.MinA, r.Value)
		} else {
			v.MaxA = max(v.MaxA, r.Value)
		}
	case "s":
		if r.Comparator == ">" {
			v.MinS = min(v.MinS, r.Value)
		} else {
			v.MaxS = max(v.MaxS, r.Value)
		}
	}

	return v
}

func (v *ValidRange) EvadeRule(r *Rule) *ValidRange {
	switch r.Category {
	case "x":
		if r.Comparator == "<" {
			v.MinX = max(v.MinX, r.Value)
		} else {
			v.MaxX = min(v.MaxX, r.Value)
		}
	case "m":
		if r.Comparator == "<" {
			v.MinM = max(v.MinM, r.Value)
		} else {
			v.MaxM = min(v.MaxM, r.Value)
		}
	case "a":
		if r.Comparator == "<" {
			v.MinA = max(v.MinA, r.Value)
		} else {
			v.MaxA = min(v.MaxA, r.Value)
		}
	case "s":
		if r.Comparator == "<" {
			v.MinS = max(v.MinS, r.Value)
		} else {
			v.MaxS = min(v.MaxS, r.Value)
		}
	}

	return v
}

func (f *Flower) FindValidRanges() []*ValidRange {
	for !f.mappedWorkflows["in"].AllPathsExplored() {
		for _, wf := range f.orderedWorkflows {
			if wf.AllPathsExplored() {
				// Been here! Let's skip.
				continue
			}

			pathsExplored := 1 // default is a gimme

			validRanges := []*ValidRange{}
			// TODO: handle default rule
			if wf.DefaultRule == "A" {
				validRanges = append(validRanges, DefaultValidRange())
			}
			if dest, ok := f.mappedWorkflows[wf.DefaultRule]; ok {
				validRanges = append(validRanges, dest.ValidRanges...)
			}

			// Walk the rules backward, applying constraints as we go
			for i := len(wf.Rules) - 1; i >= 0; i-- {
				r := wf.Rules[i]
				nr := []*ValidRange{}
				if r.Destination == "A" {
					nr = append(nr, DefaultValidRange())
					pathsExplored++
				}
				if r.Destination == "R" {
					pathsExplored++
				}
				if dest, ok := f.mappedWorkflows[r.Destination]; ok {
					if dest.AllPathsExplored() {
						nr = append(nr, dest.ValidRanges...)
						pathsExplored++
					}
				}

				tmp := []*ValidRange{}
				for _, vr := range validRanges {
					tmp = append(tmp, vr.EvadeRule(r))
				}
				for _, vr := range nr {
					tmp = append(tmp, vr.ConstrainWithRule(r))
				}
				validRanges = tmp
			}

			slog.Debug("Found valid ranges", "wf", wf.Name, "pathsExplored", pathsExplored, "ranges", validRanges)
			wf.ValidRanges = validRanges
			wf.pathsExplored = pathsExplored
			f.mappedWorkflows[wf.Name] = wf
		}
	}

	return f.mappedWorkflows["in"].ValidRanges
}

func AllCombinations(vrs []*ValidRange) int {
	combinations := 0
	for _, vr := range vrs {
		combinations += ((min(vr.MaxX, 4000) - max(vr.MinX, 1)) *
			(min(vr.MaxM, 4000) - max(vr.MinM, 1)) *
			(min(vr.MaxA, 4000) - max(vr.MinA, 1)) *
			(min(vr.MaxS, 4000) - max(vr.MinS, 1)))
		slog.Debug("Completed combinations for a valid range", "range", vr, "combinations", combinations)
	}
	return combinations
}

func partTwo(puzzleFile string) {
	slog.Info("Day Nineteen part two", "puzzle file", puzzleFile)

	flower, _ := ParseCommands(puzzleFile)

	vrs := flower.FindValidRanges()
	slog.Debug("Found valid ranges", "ranges", vrs)

	slog.Info("Found combinations", "ranges", vrs, "combinations", AllCombinations(vrs))
}
