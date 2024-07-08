package dayNineteen

import (
	"adventofcode/cmd/fileReader"
	"fmt"
	"log/slog"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/spf13/cobra"
)

type Workflow struct {
	Name        string  `@WorkflowName "{"`
	Rule        []*Rule `@@+`
	DefaultRule string  `@Destination "}"`
}

func (w *Workflow) String() string {
	return fmt.Sprintf("Workflow: %s, Rules: %v", w.Name, w.Rule)
}

type Rule struct {
	Category    string `@( "x" | "m" | "a" | "s")`
	Comparator  string `@(">"|"<")`
	Value       int    `@Int ":"`
	Destination string `@Destination ","`
}

func (r *Rule) String() string {
	return fmt.Sprintf("Rule: %v %v %v -> %v", r.Category, r.Comparator, r.Value, r.Destination)
}

type Part struct {
	XRating int `"{x=" @Int`
	MRating int `",m=" @Int`
	ARating int `",a=" @Int`
	SRating int `",s=" @Int "}"`
}

func (p *Part) String() string {
	return fmt.Sprintf("Part: %d %d %d %d", p.XRating, p.MRating, p.ARating, p.SRating)
}

func ParseCommands(puzzleFile string) ([]*Workflow, []*Part) {
	lines := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")

	myLexer := lexer.MustSimple([]lexer.SimpleRule{
		// Order matters here! Int kept stealing the leading cards before I changed the ordering.
		{"WorkflowName", `[a-z]+`},
		{"RatingCategory", `(x|m|a|s)`},
		{"Comparator", `(<|>)`},
		{"Int", `\d+`},
		{"Destination", `([a-z]+|A|R)`},
		{"Colon", `:`},
		{"Brackets", `[{}]`},
		{"Comma", `,`},
	})
	workflowParser := participle.MustBuild[Workflow](
		participle.Lexer(myLexer),
	)
	partParser := participle.MustBuild[Part](
		participle.Lexer(myLexer),
	)

	workflows := []*Workflow{}
	parts := []*Part{}
	ws := true
	for _, l := range lines {
		if l == "" {
			ws = false
			continue
		}

		if ws {
			w, err := workflowParser.ParseBytes("", []byte(l))
			if err != nil {
				slog.Error("failed to parse", "w", w, "err", err)
				panic(err)
			}
			workflows = append(workflows, w)
			continue
		}

		slog.Debug("Parsing part", "line", l)
		p, err := partParser.ParseBytes("", []byte(l))
		if err != nil {
			slog.Error("failed to parse", "p", p, "err", err)
			panic(err)
		}
		parts = append(parts, p)
	}

	return workflows, parts
}
func partOne(puzzleFile string) {
	slog.Info("Day Nineteen part one", "puzzle file", puzzleFile)

	workflows, parts := ParseCommands(puzzleFile)

	slog.Debug("Parsed", "workflows", workflows, "parts", parts)
}

func partTwo(puzzleFile string) {
	slog.Info("Day Nineteen part two", "puzzle file", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "dayNineteen",
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
