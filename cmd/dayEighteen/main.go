package dayEighteen

import (
	"adventofcode/cmd/fileReader"
	"fmt"
	"log/slog"
	"math"
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

func ParseCommands(puzzleFile string) []*DigCommand {
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

type Map struct {
	MaxX            int
	MinX            int
	MaxY            int
	MinY            int
	VerticesOrdered []*Coordinate
	VerticesMapped  map[string]bool
}

func (m *Map) AddVertex(c *Coordinate) {
	m.VerticesOrdered = append(m.VerticesOrdered, c)
	m.VerticesMapped[c.String()] = true

	if c.Col > m.MaxX {
		m.MaxX = c.Col
	}
	if c.Col < m.MinX {
		m.MinX = c.Col
	}
	if c.Row > m.MaxY {
		m.MaxY = c.Row
	}
	if c.Row < m.MinY {
		m.MinY = c.Row
	}
}

func BuildMap(commands []*DigCommand) *Map {
	maxX, maxY, minX, minY := 0, 0, 0, 0
	theMap := &Map{maxX, minX, maxY, minY, []*Coordinate{}, map[string]bool{}}
	pos := &Coordinate{0, 0}
	for _, c := range commands {
		switch c.Dir {
		case "R":
			pos.Col += c.Dist
		case "L":
			pos.Col -= c.Dist
		case "U":
			pos.Row += c.Dist
		case "D":
			pos.Row -= c.Dist
		}

		theMap.AddVertex(&Coordinate{pos.Row, pos.Col})
	}

	return theMap
}

// Hello [Shoelace Formula](https://en.wikipedia.org/wiki/Shoelace_formula#Shoelace_formula)!
func CalculateArea(theMap *Map) float64 {
	area := 0.0
	borderLen := 0.0
	for i := 0; i < len(theMap.VerticesOrdered); i++ {
		c1 := theMap.VerticesOrdered[i]
		c2 := theMap.VerticesOrdered[(i+1)%len(theMap.VerticesOrdered)]
		area += float64(c1.Row*c2.Col - c2.Row*c1.Col)
		borderLen += math.Abs(float64(c1.Col - c2.Col + c2.Row - c1.Row))
	}

	return area/2 + borderLen/2 + 1
}

func PrintableGrid(theMap *Map) string {
	printableGrid := []string{}
	for row := theMap.MaxY; row >= theMap.MinY; row-- {
		b := strings.Builder{}
		b.Grow(theMap.MaxX - theMap.MinX + 1)
		for col := theMap.MinX; col <= theMap.MaxX; col++ {
			c := &Coordinate{row, col}
			if theMap.VerticesMapped[c.String()] {
				b.WriteString("#")
			} else {
				b.WriteString(".")
			}
		}
		printableGrid = append(printableGrid, b.String())
	}

	return strings.Join(printableGrid, "\n")
}

func partOne(puzzleFile string) {
	slog.Info("Day Eighteen part one", "puzzle file", puzzleFile)

	commands := ParseCommands(puzzleFile)

	theMap := BuildMap(commands)
	slog.Debug("got a map!", "theMap", theMap)

	printGrid := PrintableGrid(theMap)

	filledPositions := CalculateArea(theMap)

	os.WriteFile("/tmp/dayEighteenGrid.txt", []byte(printGrid), 0644)

	slog.Info("finished digging", "filled positions", filledPositions)
}

func partTwo(puzzleFile string) {
	slog.Info("Day Eighteen part two", "puzzle file", puzzleFile)

	bustedCommands := ParseCommands(puzzleFile)

	fixedCommands := []*DigCommand{}
	for _, bustedCommand := range bustedCommands {
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
		fixedCommands = append(fixedCommands, c)
	}

	theMap := BuildMap(fixedCommands)
	printableGrid := PrintableGrid(theMap)
	filledPositions := CalculateArea(theMap)

	slog.Info("finished digging", "filled positions", filledPositions)

	os.WriteFile("/tmp/dayEighteenGrid.txt", []byte(printableGrid), 0644)
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
