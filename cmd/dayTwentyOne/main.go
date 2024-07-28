package dayTwentyOne

import (
	"adventofcode/cmd/coordinates"
	"adventofcode/cmd/fileReader"
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

func NextInfinitePositions(g []string, c *coordinates.Coordinate) []*coordinates.Coordinate {
	next := []*coordinates.Coordinate{}
	for _, m := range coordinates.GridMoves() {
		n := c.Move(m)
		row := n.Row % len(g)
		if row < 0 {
			row = len(g) + row
		}
		col := n.Col % len(g[0])
		if col < 0 {
			col = len(g[0]) + col
		}
		if g[row][col] == '#' {
			continue
		}
		next = append(next, n)
	}
	return next
}

/*
*
Okay, so this is correct but too slow. I spent several commits (look at the history) attempting
various caching and speeding up techniques, but [Reddit](https://www.reddit.com/r/adventofcode/comments/18orn0s/2023_day_21_part_2_links_between_days/)
has some great points about looking at this as a _polynomial_. You can even use the day nine solution
to extrapolate to an answer. The "mathy" trick is to find the polynomial. Lagrange interpolation is a good
fit for this.

Using [this site](https://www.dcode.fr/lagrange-interpolating-polynomial) we can find the polynomial.
I ran mine on several points and got the following polynomial:
14888*x^2/17161 + 26154*x/17161 − 213738/17161
*
*/
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
		for _, n := range NextInfinitePositions(g, next.Pos) {
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

	slog.Info("Day TwentyOne part two", "reachable plots", len(finalPlots))
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
