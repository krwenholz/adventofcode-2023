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

func (s *Sequence) predictNextValue() int {
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
	slog.Debug("extrapolated", "rows", rows)

	return rows[0][len(rows[0])-1]
}

func partOne(puzzleFile string) {
	slog.Info("Day Nine part one", "puzzle file", puzzleFile)

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

	scanner := scanner.NewScanner[Sequence](parser, puzzleFile)
	sum := 0
	for scanner.Scan() {
		seq := scanner.Struct()
		sum += seq.predictNextValue()
	}

	slog.Info("Finished day nine part one", "sum", sum)
}

func partTwo(puzzleFile string) {
	slog.Info("Day Nine part two", "puzzle file", puzzleFile)
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
