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

func fmtBinary(b uint16) string {
	return fmt.Sprintf("%016b", b)
}

func lmask(length, mirrorPoint int) uint16 {
	if mirrorPoint > length/2 {
		return uint16(math.Pow(2, float64(length-mirrorPoint)) - 1)
	}
	return uint16(math.Pow(2, 16) - 1)
}

func rmask(length, mirrorPoint int) uint16 {
	if mirrorPoint <= length/2 {
		return uint16(math.Pow(2, float64(mirrorPoint)) - 1)
	}
	return uint16(math.Pow(2, float64(length-mirrorPoint)) - 1)
}

func mirror(original uint16, length, mirrorPoint int) bool {
	left := original >> (length - mirrorPoint)
	left &= lmask(length, mirrorPoint)

	reversed := bits.Reverse16(original)
	fmt.Println(fmtBinary(reversed))
	reversed >>= (16 - (length - mirrorPoint))
	fmt.Println(fmtBinary(reversed))
	reversed &= rmask(length, mirrorPoint)
	fmt.Println(fmtBinary(reversed))

	diff := left - reversed

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

	return diff == 0 //|| bits.TrailingZeros(diff) >= originalL-mirrorPoint
}

func splitIndex(p []string) int {
	initialShifts := make([]int, 0)
	for i := len(p[0]) / 2; i < len(p[0]); i++ {
		initialShifts = append(initialShifts, i)
	}
	validShifts := [][]int{initialShifts}

	for _, r := range p {
		previousShifts := validShifts[len(validShifts)-1]
		if len(previousShifts) == 0 {
			continue
		}

		theseShifts := []int{}
		firstShift := previousShifts[0]
		tmp, _ := strconv.ParseInt(r, 2, 10)
		bin := uint16(tmp)

		if mirror(bin, len(r), firstShift) {
			theseShifts = append(theseShifts, firstShift)
		}

		for _, s := range previousShifts[1:] {
			if mirror(bin, len(r), s) {
				theseShifts = append(theseShifts, s)
			}
		}
		validShifts = append(validShifts, theseShifts)
	}

	slog.Debug(
		"finished computing shifts",
		"validShifts", validShifts,
	)

	if len(validShifts) >= len(p) && len(validShifts[len(validShifts)-1]) > 0 {
		splitIdx := validShifts[len(validShifts)-1][0]
		return splitIdx
	}
	return 0
}

func verticalLeftSplit(pattern []string) int {
	leftSplitIdx := splitIndex(pattern)
	if leftSplitIdx > 0 {
		slog.Debug(
			"finished computing valid vertical split",
			"left split index", leftSplitIdx)
	}
	return leftSplitIdx
}

func horizontalAboveSplit(pattern []string) int {
	rotatedPattern := make([]string, len(pattern[0]))
	for i := 0; i < len(pattern[0]); i++ {
		for j := 0; j < len(pattern); j++ {
			rotatedPattern[i] = rotatedPattern[i] + string(pattern[j][i])
		}
	}
	slog.Debug("rotated pattern", "p", rotatedPattern)
	aboveSplitIdx := splitIndex(rotatedPattern)
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
		l := verticalLeftSplit(p)
		verticalLeftSum += l

		// try horizontal
		h := horizontalAboveSplit(p)
		horizontalAboveSum += h

		slog.Debug(
			"finished computing split",
			"pattern", i,
			"vert", l,
			"horizontal", h,
		)
	}

	slog.Info("Finished day thirteen part one", "expected", ans, "summary", verticalLeftSum+horizontalAboveSum*100)
}

func partTwo(puzzleFile string) {
	slog.Info("Day Thirteen part two", "puzzle file", puzzleFile)
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
