package dayNine

import (
	"adventofcode/cmd/scanner"
	"log"
	"log/slog"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/spf13/cobra"
)

type Sequence struct {
	Values []int `@Int+`
}

func zeroed(row []int) bool {
	for _, v := range row {
		if v != 0 {
			return false
		}
	}
	return true
}

func (s *Sequence) extrapolate() (int, int) {
	rows := [][]int{s.Values}
	for {
		lastRow := rows[len(rows)-1]
		nextRow := []int{}
		for i := 0; i <= len(lastRow)-2; i++ {
			nextRow = append(nextRow, lastRow[i+1]-lastRow[i])
		}

		rows = append(rows, nextRow)
		if zeroed(nextRow) {
			break
		}
	}

	for i := len(rows) - 2; i >= 0; i-- {
		curNum := rows[i][len(rows[i])-1]
		belowNum := rows[i+1][len(rows[i+1])-1]
		rows[i] = append(rows[i], curNum+belowNum)
	}
	slog.Debug("extrapolated forward", "rows", rows)

	for i := len(rows) - 2; i >= 0; i-- {
		curNum := rows[i][0]
		belowNum := rows[i+1][0]
		rows[i] = append([]int{curNum - belowNum}, rows[i]...)
	}
	slog.Debug("extrapolated backward", "rows", rows)

	return rows[0][0], rows[0][len(rows[0])-1]
}

func newScanner(puzzleFile string) *scanner.PuzzleScanner[Sequence] {
	sequenceLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Int", `(-)?(\d*\.)?\d+`},
		{"Whitespace", `[ \t]+`},
	})
	parser, err := participle.Build[Sequence](
		participle.Lexer(sequenceLexer),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		log.Fatal(err)
	}

	return scanner.NewScanner[Sequence](parser, puzzleFile)
}

func partOne(puzzleFile string) {
	slog.Info("Day Nine part one", "puzzle file", puzzleFile)

	s := newScanner(puzzleFile)
	sum := 0
	for s.Scan() {
		seq := s.Struct()
		_, v := seq.extrapolate()
		sum += v
	}

	slog.Info("Finished day nine part one", "sum", sum)
}

func partTwo(puzzleFile string) {
	slog.Info("Day Nine part two", "puzzle file", puzzleFile)

	s := newScanner(puzzleFile)
	sum := 0
	for s.Scan() {
		seq := s.Struct()
		v, _ := seq.extrapolate()
		sum += v
	}

	slog.Info("Finished day nine part two", "sum", sum)
}

var Cmd = &cobra.Command{
	Use: "dayNine",
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
