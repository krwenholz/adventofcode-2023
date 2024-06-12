package dayThirteen

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"math/bits"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

/**
121122112
 211221121

 one shift plus half (4) means five is shift point
**/

func fmtBinary(b uint) string {
	return fmt.Sprintf("%064b", b)
}

func partOne(puzzleFile string) {
	slog.Info("Day Thirteen part one", "puzzle file", puzzleFile)
	f, err := os.Open(puzzleFile)
	if err != nil {
		log.Fatal(err)
	}

	sc := bufio.NewScanner(f)

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

	verticalLeftSum := -1
	horizontalAboveSum := -1

	for _, p := range patterns {
		// try horizontal
		initialShifts := make([]int, 0)
		for i := 0; i < len(p[0]); i++ {
			initialShifts = append(initialShifts, i)
		}
		validShifts := [][]int{initialShifts}

		for _, r := range p {
			previousShifts := validShifts[len(validShifts)-1]
			if len(previousShifts) == 0 {
				continue
			}

			theseShifts := []int{}
			tmp, _ := strconv.ParseInt(r, 2, 10)
			bin := uint(tmp)
			tmpR := []rune{}
			for _, c := range r {
				tmpR = append(tmpR, c)
			}
			tmp, _ = strconv.ParseInt(string(tmpR), 2, 10)
			binR := uint(tmp)

			diff, _ := bits.Sub(bin, binR, 0)
			if diff == 0 {
				theseShifts = append(theseShifts, 0)
			}
			for _, s := range previousShifts {
				sBin := bin >> s
				sBinR := binR << s

				diff, _ := bits.Sub(sBin, sBinR, 0)
				slog.Debug("diff", "bin", fmtBinary(sBin), "binR", fmtBinary(sBinR), "diff", fmtBinary(diff))
				if diff == 0 {
					theseShifts = append(theseShifts, s)
				}
			}
			validShifts = append(validShifts, theseShifts)
		}
		slog.Debug("finished computing valid vertical shifts", "validShifts", validShifts)

		if len(validShifts) == len(p[0]) {
			verticalLeftSum += validShifts[len(validShifts)-1][0]
		}

		// try vertical
	}

	slog.Info("Finished day thirteen part one", "summary", verticalLeftSum+horizontalAboveSum*100)
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
