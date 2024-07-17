package dayTwentyOne

import (
	"adventofcode/cmd/coordinates"
	"adventofcode/cmd/fileReader"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func NextPositions(g []string, c *coordinates.Coordinate) []*coordinates.Coordinate {
	next := []*coordinates.Coordinate{}
	for _, m := range coordinates.GridMoves() {
		n := c.Move(m)
		if !(0 <= n.Row && n.Row < len(g) && 0 <= n.Col && n.Col < len(g[0])) {
			continue
		}
		if g[n.Row][n.Col] == '#' {
			continue
		}
		next = append(next, n)
	}
	return next
}

type Step struct {
	Pos        *coordinates.Coordinate
	StepNumber int
}

func (s *Step) Hash() string {
	return fmt.Sprintf("%s-%d", s.Pos.String(), s.StepNumber%2)
}

func ReachablePlots(g []string, start *coordinates.Coordinate, steps int) int {
	curs := []*Step{{start, 0}}
	seen := map[string]bool{}
	finalPlots := map[string]bool{}
	for len(curs) > 0 {
		/**
		we'll grab next positions and potentially add them to curs to step from
		all positions can be final by
		**/
		next := curs[0]
		curs = curs[1:]
		nextStepNumber := next.StepNumber + 1
		if nextStepNumber > steps {
			continue
		}
		for _, n := range NextPositions(g, next.Pos) {
			if _, ok := seen[n.String()]; !ok {
				seen[n.String()] = true
				curs = append(curs, &Step{n, nextStepNumber})
			}
			// Can only enter a plot at the end if we're at the step count or have an even number of steps left
			if nextStepNumber == steps || (steps-nextStepNumber)%2 == 0 {
				finalPlots[n.String()] = true
			}
		}
	}

	PrintGrid(g, finalPlots)
	return len(finalPlots)
}

func PrintGrid(g []string, plots map[string]bool) {
	if strings.ToLower(os.Getenv("LOG_LEVEL")) != "debug" {
		return
	}
	for i, row := range g {
		for j, c := range row {
			if c == '#' {
				print("#")
			} else if _, ok := plots[(&coordinates.Coordinate{Row: i, Col: j}).String()]; ok {
				print("O")
			} else if c == 'S' {
				print("S")
			} else {
				print(".")
			}
		}
		println()
	}
}

func partOne(puzzleFile string, stepCount int) {
	slog.Info("Day TwentyOne part one", "puzzle file", puzzleFile)
	g := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")

	start := &coordinates.Coordinate{Row: 0, Col: 0}
	for i, row := range g {
		for j, c := range row {
			if c == 'S' {
				start = &coordinates.Coordinate{Row: i, Col: j}
				break
			}
		}
	}

	plots := ReachablePlots(g, start, stepCount)

	slog.Info("Day TwentyOne part one", "reachable plots", plots)
}

// Returns next valid coordinates and any coordinates that would leave the current grid
func NextInfinitePositions(g []string, c *coordinates.Coordinate) ([]*coordinates.Coordinate, []*coordinates.Coordinate) {
	next := []*coordinates.Coordinate{}
	starts := []*coordinates.Coordinate{}
	for _, m := range coordinates.GridMoves() {
		n := c.Move(m)

		isStart := false
		if n.Row < 0 {
			n = &coordinates.Coordinate{Row: len(g) - 1, Col: n.Col}
			isStart = true
		} else if n.Row >= len(g) {
			n = &coordinates.Coordinate{Row: 0, Col: n.Col}
			isStart = true
		} else if n.Col < 0 {
			n = &coordinates.Coordinate{Row: n.Row, Col: len(g[0]) - 1}
			isStart = true
		} else if n.Col >= len(g[0]) {
			n = &coordinates.Coordinate{Row: n.Row, Col: 0}
			isStart = true
		}

		if g[n.Row][n.Col] == '#' {
			continue
		}

		if isStart {
			starts = append(starts, n)
			continue
		}

		next = append(next, n)
	}
	return next, starts
}

func partTwo(puzzleFile string, steps int) {
	slog.Info("Day TwentyOne part two", "puzzle file", puzzleFile)
	g := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")

	start := &coordinates.Coordinate{Row: 0, Col: 0}
	for i, row := range g {
		for j, c := range row {
			if c == 'S' {
				start = &coordinates.Coordinate{Row: i, Col: j}
				break
			}
		}
	}

	starts := []*Step{{start, 0}}
	startsSeen := map[string]int{starts[0].Hash(): 1}
	startsToPlots := map[string]int{}
	for len(starts) > 0 {
		/**
		for every start, we'll run it through, collecting new starts we haven't seen
		or incrementing their count
		**/
		s := starts[0]
		starts = starts[1:]

		seen := map[string]bool{}
		finalPlots := map[string]bool{}
		curs := []*Step{s}

		for len(curs) > 0 {
			next := curs[0]
			curs = curs[1:]
			nextStepNumber := next.StepNumber + 1
			if nextStepNumber > steps {
				continue
			}
			nexts, newStarts := NextInfinitePositions(g, next.Pos)
			for _, n := range nexts {
				if _, ok := seen[n.String()]; !ok {
					seen[n.String()] = true
					curs = append(curs, &Step{n, nextStepNumber})
				}
				// Can only enter a plot at the end if we're at the step count or have an even number of steps left
				if nextStepNumber == steps || (steps-nextStepNumber)%2 == 0 {
					finalPlots[n.String()] = true
				}
			}

			for _, n := range newStarts {
				nextStep := &Step{n, nextStepNumber}
				if _, ok := startsSeen[nextStep.Hash()]; !ok {
					startsSeen[nextStep.Hash()] = 1
					starts = append(starts, nextStep)
					slog.Debug("Adding new start", "start", nextStep)
				} else {
					startsSeen[nextStep.Hash()]++
				}
			}
		}

		startsToPlots[s.Hash()] = len(finalPlots)
	}

	slog.Debug("Finished simulating", "starts seen", startsSeen, "starts to plots", startsToPlots)

	totalPlots := 0
	for start, plots := range startsToPlots {
		totalPlots += plots * startsSeen[start]
	}

	slog.Info("Day TwentyOne part two", "steps", steps, "reachable plots", totalPlots)
}

var Cmd = &cobra.Command{
	Use: "dayTwentyOne",
	Run: func(cmd *cobra.Command, args []string) {
		puzzleInput, _ := cmd.Flags().GetString("puzzle-input")
		stepCount, _ := cmd.Flags().GetInt("step-count")
		if !cmd.Flag("part-two").Changed {
			partOne(puzzleInput, stepCount)
		} else {
			partTwo(puzzleInput, stepCount)
		}
	},
}

func init() {
	Cmd.Flags().Bool("part-two", false, "Whether to run part two of the day's challenge")
	Cmd.Flags().Int("step-count", 4, "Steps to take")
}
