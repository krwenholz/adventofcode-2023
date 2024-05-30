package daySix

import (
	"log/slog"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/spf13/cobra"
)

type Input struct {
	Times     []int `"Time" @Int+ EOL`
	Distances []int `"Distance" @Int+`
}

func parse(puzzleFile string) *Input {
	f, err := os.ReadFile(puzzleFile)
	if err != nil {
		slog.Error("failed to parse", err)
		panic(err)
	}

	inputLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Int", `(\d*\.)?\d+`},
		{"Ident", `[a-zA-Z_]\w*`},
		{"EOL", `\n`},
		{"Colon", `:`},
		{"Whitespace", `[ \t]+`},
		{"ToSeparator", `-to-`},
	})
	parser, err := participle.Build[Input](
		participle.Lexer(inputLexer),
		participle.Elide("Colon", "Whitespace"),
	)
	if err != nil {
		slog.Error("failed to parse", "err", err)
		panic(err)
	}

	input, err := parser.ParseBytes("", f)
	if err != nil {
		slog.Error("failed to parse", "so far", input, "err", err)
		panic(err)
	}

	return input
}

func distanceForTime(timeHeld, timeRunning int) int {
	return timeHeld * timeRunning
}

func validHolds(times []int, distances []int) int {
	validHolds := make(map[int]int)

	for i, totalTime := range times {
		bestDistance := distances[i]
		firstGoodHold := 0
		for timeHeld := 0; timeHeld < totalTime; timeHeld++ {
			timeRunning := totalTime - timeHeld
			distance := distanceForTime(timeHeld, timeRunning)
			if distance > bestDistance {
				firstGoodHold = timeHeld - 1
				break
			}
		}
		// The times held that don't work are symmetric on either side of a curve of good times
		validHolds[i] = totalTime - (firstGoodHold * 2) - 1
	}

	answer := 1
	for _, v := range validHolds {
		answer = answer * v
	}

	slog.Debug("calculated holds", "answer", validHolds, "answer", answer)
	return answer
}

func partOne(puzzleFile string) {
	input := parse(puzzleFile)
	answer := validHolds(input.Times, input.Distances)
	slog.Info("Day six part one", "input", input, "answer", answer)
}

func partTwo(puzzleFile string) {
	slog.Info("Day six part two")
}

var Cmd = &cobra.Command{
	Use: "daySix",
	Run: func(cmd *cobra.Command, args []string) {
		puzzleInput, _ := cmd.Flags().GetString("puzzle-input")
		if !cmd.Flag("part-two").Changed {
			partOne(puzzleInput)
		} else {
			partTwo(puzzleInput)
		}
	},
}
