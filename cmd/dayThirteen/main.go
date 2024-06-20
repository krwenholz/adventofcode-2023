package dayThirteen

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"math"
	"math/bits"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

/**
101100110
 011001101

 one shift plus half (4) means five is shift point
**/

func fmtBinary(b uint32) string {
	return fmt.Sprintf("%032b", b)
}

func lmask(length, mirrorPoint int) uint32 {
	if mirrorPoint > length/2 {
		return uint32(math.Pow(2, float64(length-mirrorPoint)) - 1)
	}
	return uint32(math.Pow(2, 32) - 1)
}

func rmask(length, mirrorPoint int) uint32 {
	if mirrorPoint <= length/2 {
		return uint32(math.Pow(2, float64(mirrorPoint)) - 1)
	}
	return uint32(math.Pow(2, float64(length-mirrorPoint)) - 1)
}

func mirror(original uint32, length, mirrorPoint int) uint32 {
	left := original >> (length - mirrorPoint)
	left &= lmask(length, mirrorPoint)

	reversed := bits.Reverse32(original)
	reversed >>= (32 - (length - mirrorPoint))
	reversed &= rmask(length, mirrorPoint)

	diff := left ^ reversed

	slog.Debug("diff",
		"original",
		fmtBinary(original),
		"left",
		fmtBinary(left),
		"reversed",
		fmtBinary(reversed),
		"originalL",
		length,
		"mirrorPoint",
		mirrorPoint,
		"diff",
		fmtBinary(diff),
	)

	return diff //|| bits.TrailingZeros(diff) >= originalL-mirrorPoint
}

type Shift struct {
	S    int
	DAcc int
}

func (s *Shift) String() string {
	return fmt.Sprintf("[S: %d, DAcc: %d]", s.S, s.DAcc)
}

func splitIndex(p []string, allowedDifferences int) int {
	initialShifts := make([]*Shift, 0)
	for i := 1; i < len(p[0]); i++ {
		initialShifts = append(initialShifts, &Shift{i, 0})
	}
	validShifts := [][]*Shift{initialShifts}

	for _, r := range p {
		previousShifts := validShifts[len(validShifts)-1]
		if len(previousShifts) == 0 {
			continue
		}

		theseShifts := []*Shift{}
		tmp, _ := strconv.ParseInt(r, 2, 32)
		bin := uint32(tmp)

		for _, s := range previousShifts {
			newS := &Shift{
				S:    s.S,
				DAcc: s.DAcc + bits.OnesCount32(mirror(bin, len(r), s.S)),
			}
			if newS.DAcc <= allowedDifferences {
				theseShifts = append(theseShifts, newS)
			}
		}
		validShifts = append(validShifts, theseShifts)
	}

	slog.Debug(
		"finished computing shifts",
		"validShifts", validShifts,
	)

	if len(validShifts) >= len(p) && len(validShifts[len(validShifts)-1]) > 0 {
		for _, s := range validShifts[len(validShifts)-1] {
			if s.DAcc == allowedDifferences {
				return s.S
			}
		}
	}
	return 0
}

func verticalLeftSplit(pattern []string, allowedDifferences int) int {
	leftSplitIdx := splitIndex(pattern, allowedDifferences)
	if leftSplitIdx > 0 {
		slog.Debug(
			"finished computing valid vertical split",
			"left split index", leftSplitIdx)
	}
	return leftSplitIdx
}

func horizontalAboveSplit(pattern []string, allowedDifferences int) int {
	rotatedPattern := make([]string, len(pattern[0]))
	for i := 0; i < len(pattern[0]); i++ {
		for j := 0; j < len(pattern); j++ {
			rotatedPattern[i] = rotatedPattern[i] + string(pattern[j][i])
		}
	}
	slog.Debug("rotated pattern", "p", rotatedPattern)
	aboveSplitIdx := splitIndex(rotatedPattern, allowedDifferences)
	if aboveSplitIdx > 0 {
		slog.Debug(
			"finished computing valid horizontal split",
			"above split index", aboveSplitIdx)
	}
	return aboveSplitIdx
}

func partOne(puzzleFile string) {
	slog.Info("Day Thirteen part one", "puzzle file", puzzleFile)
	f, err := os.Open(puzzleFile)
	if err != nil {
		log.Fatal(err)
	}

	sc := bufio.NewScanner(f)
	sc.Scan()
	ans := sc.Text()

	patterns := make([][]string, 1)
	for sc.Scan() {
		row := ""
		t := sc.Text()
		if t == "" {
			patterns = append(patterns, []string{})
			continue
		}
		for _, c := range t {
			if c == '#' {
				row += "1"
			} else {
				row += "0"
			}
		}
		patterns[len(patterns)-1] = append(patterns[len(patterns)-1], row)
	}

	verticalLeftSum := 0
	horizontalAboveSum := 0

	for i, p := range patterns {
		// try vertical
		l := verticalLeftSplit(p, 0)
		verticalLeftSum += l

		// try horizontal
		h := horizontalAboveSplit(p, 0)
		horizontalAboveSum += h

		slog.Debug(
			"finished computing split",
			"pattern", i,
			"vert", l,
			"horizontal", h,
		)
	}

	slog.Info("Finished day thirteen part one", "expected", ans, "value", verticalLeftSum+horizontalAboveSum*100)
}

func partTwo(puzzleFile string) {
	slog.Info("Day Thirteen part two", "puzzle file", puzzleFile)
	f, err := os.Open(puzzleFile)
	if err != nil {
		log.Fatal(err)
	}

	sc := bufio.NewScanner(f)
	sc.Scan()
	ans := sc.Text()

	patterns := make([][]string, 1)
	for sc.Scan() {
		row := ""
		t := sc.Text()
		if t == "" {
			patterns = append(patterns, []string{})
			continue
		}
		for _, c := range t {
			if c == '#' {
				row += "1"
			} else {
				row += "0"
			}
		}
		patterns[len(patterns)-1] = append(patterns[len(patterns)-1], row)
	}

	verticalLeftSum := 0
	horizontalAboveSum := 0

	for i, p := range patterns {
		// try vertical
		l := verticalLeftSplit(p, 1)
		verticalLeftSum += l

		// try horizontal
		h := horizontalAboveSplit(p, 1)
		horizontalAboveSum += h

		slog.Debug(
			"finished computing split",
			"pattern", i,
			"vert", l,
			"horizontal", h,
		)
	}

	slog.Info("Finished day thirteen part two", "expected", ans, "value", verticalLeftSum+horizontalAboveSum*100)
}

var Cmd = &cobra.Command{
	Use: "dayThirteen",
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
