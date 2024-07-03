package dayEighteen

import (
	"adventofcode/cmd/fileReader"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/spf13/cobra"
)

type Coordinate struct {
	Row int
	Col int
}

func (c *Coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", c.Row, c.Col)
}

func (c *Coordinate) Equals(other *Coordinate) bool {
	return c.Row == other.Row && c.Col == other.Col
}

type DigCommand struct {
	Dir   string `@Direction`
	Dist  int    `@Int`
	Color string `@Color`
}

func (d DigCommand) String() string {
	return fmt.Sprintf("%s %d #%s", d.Dir, d.Dist, d.Color)
}

type Space struct {
	Shape string
	Color string
}

func parseCommands(puzzleFile string) []*DigCommand {
	rawCommands := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")

	commandLexer := lexer.MustSimple([]lexer.SimpleRule{
		// Order matters here! Int kept stealing the leading cards before I changed the ordering.
		{"Direction", `[RLDU]`},
		{"Color", `[a-z0-9]{6}`},
		{"Int", `(\d*\.)?\d+`},
		{"Whitespace", `[ \t]+`},
		{"Parens", `(\(|\))`},
		{"Pound", `#`},
	})
	parser := participle.MustBuild[DigCommand](
		participle.Lexer(commandLexer),
		participle.Elide("Whitespace", "Parens", "Pound"),
	)

	commands := []*DigCommand{}
	for _, rawC := range rawCommands {
		c, err := parser.ParseBytes("", []byte(rawC))
		if err != nil {
			slog.Error("failed to parse", "c", c, "err", err)
			panic(err)
		}
		commands = append(commands, c)
	}

	return commands
}

func cornerShape(prev, next *DigCommand) string {
	switch prev.Dir {
	case "R":
		switch next.Dir {
		case "U":
			return "J"
		case "D":
			return "7"
		case "R":
			return "-"
		case "L":
			return "-"
		}
	case "L":
		switch next.Dir {
		case "U":
			return "L"
		case "D":
			return "F"
		case "R":
			return "-"
		case "L":
			return "-"
		}
	case "U":
		switch next.Dir {
		case "U":
			return "|"
		case "D":
			return "|"
		case "R":
			return "F"
		case "L":
			return "7"
		}
	case "D":
		switch next.Dir {
		case "U":
			return "|"
		case "D":
			return "|"
		case "R":
			return "L"
		case "L":
			return "J"
		}
	}
	return ""
}

func partOne(puzzleFile string) {
	slog.Info("Day Eighteen part one", "puzzle file", puzzleFile)

	commands := parseCommands(puzzleFile)

	maxX, maxY, minX, minY := 0, 0, 0, 0
	grid := map[int]map[int]*Space{}
	pos := &Coordinate{0, 0}
	for cI, c := range commands {
		for i := 0; i < c.Dist; i++ {
			switch c.Dir {
			case "R":
				pos.Col++
			case "L":
				pos.Col--
			case "U":
				pos.Row++
			case "D":
				pos.Row--
			}

			if pos.Col > maxX {
				maxX = pos.Col
			}
			if pos.Col < minX {
				minX = pos.Col
			}
			if pos.Row > maxY {
				maxY = pos.Row
			}
			if pos.Row < minY {
				minY = pos.Row
			}

			if grid[pos.Row] == nil {
				grid[pos.Row] = map[int]*Space{}
			}

			var shape string
			if i == c.Dist-1 {
				shape = cornerShape(c, commands[(cI+1)%len(commands)])
			} else {
				switch c.Dir {
				case "R":
					shape = "-"
				case "L":
					shape = "-"
				case "U":
					shape = "|"
				case "D":
					shape = "|"
				}
			}

			grid[pos.Row][pos.Col] = &Space{shape, c.Color}
		}
	}

	filledPositions := 0
	printGrid := ""
	for row := maxY; row >= minY; row-- {
		filledThisRow := 0
		wallCount := 0
		for col := minX; col <= maxX; col++ {
			if s, ok := grid[row][col]; ok {
				filledThisRow++
				printGrid += s.Shape
				switch s.Shape {
				case "|":
					wallCount++
				case "L":
					wallCount++
				case "J":
					wallCount++
				}
			} else {
				if wallCount%2 == 1 {
					filledThisRow++
				}
				printGrid += "."
			}
		}

		filledPositions += filledThisRow

		printGrid += fmt.Sprintf(" (%d) \n", filledThisRow)
		slog.Debug("row filled", "row", row, "filledThisRow", filledThisRow)
	}

	os.WriteFile("/tmp/dayEighteenGrid.txt", []byte(printGrid), 0644)

	slog.Info("finished digging", "filled positions", filledPositions)
}

func partTwo(puzzleFile string) {
	slog.Info("Day Eighteen part two", "puzzle file", puzzleFile)

	commands := parseCommands(puzzleFile)

	maxX, maxY, minX, minY := 0, 0, 0, 0
	grid := map[int]map[int]*Space{}
	pos := &Coordinate{0, 0}
	for cI, bustedCommand := range commands {
		dist, _ := strconv.ParseInt(bustedCommand.Color[:5], 16, 32)
		var dir string
		switch bustedCommand.Color[5:] {
		case "0":
			dir = "R"
		case "1":
			dir = "D"
		case "2":
			dir = "L"
		case "3":
			dir = "U"
		}

		c := &DigCommand{dir, int(dist), bustedCommand.Color}
		slog.Debug("fixed command", "bustedCommand", bustedCommand, "c", c)

		for i := 0; i < c.Dist; i++ {
			switch c.Dir {
			case "R":
				pos.Col++
			case "L":
				pos.Col--
			case "U":
				pos.Row++
			case "D":
				pos.Row--
			}

			if pos.Col > maxX {
				maxX = pos.Col
			}
			if pos.Col < minX {
				minX = pos.Col
			}
			if pos.Row > maxY {
				maxY = pos.Row
			}
			if pos.Row < minY {
				minY = pos.Row
			}

			if grid[pos.Row] == nil {
				grid[pos.Row] = map[int]*Space{}
			}

			var shape string
			if i == c.Dist-1 {
				shape = cornerShape(c, commands[(cI+1)%len(commands)])
			} else {
				switch c.Dir {
				case "R":
					shape = "-"
				case "L":
					shape = "-"
				case "U":
					shape = "|"
				case "D":
					shape = "|"
				}
			}

			grid[pos.Row][pos.Col] = &Space{shape, c.Color}
		}
	}

	slog.Debug("Finished building grid", "size", len(grid), "width", maxX-minX, "height", maxY-minY)

	filledPositions := 0
	//printGrid := &strings.Builder{}
	printGrid := &NoopStringBuilder{}
	printGrid.Grow(maxY - minY + maxX - minX)
	for row := maxY; row >= minY; row-- {
		filledThisRow := 0
		wallCount := 0
		for col := minX; col <= maxX; col++ {
			if s, ok := grid[row][col]; ok {
				filledThisRow++
				printGrid.WriteString(s.Shape)
				switch s.Shape {
				case "|":
					wallCount++
				case "L":
					wallCount++
				case "J":
					wallCount++
				}
			} else {
				if wallCount%2 == 1 {
					filledThisRow++
				}
				printGrid.WriteString(".")
			}
		}

		filledPositions += filledThisRow

		printGrid.WriteString(fmt.Sprintf(" (%d) \n", filledThisRow))

		if row%500 == 0 {
			slog.Debug("row filling", "row", row, "percent", float64((maxY-row)*100)/float64((maxY-minY)))
		}
	}

	slog.Info("finished digging", "filled positions", filledPositions)

	os.WriteFile("/tmp/dayEighteenGrid.txt", []byte(printGrid.String()), 0644)
}

type NoopStringBuilder struct{}

func (sb *NoopStringBuilder) WriteString(p string) (n int, err error) {
	return len(p), nil
}

func (sb *NoopStringBuilder) Grow(i int) {}

func (sb *NoopStringBuilder) String() string {
	return ""
}

var Cmd = &cobra.Command{
	Use: "dayEighteen",
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
