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
	Name        string  `@Ident "{"`
	Rule        []*Rule `@@*`
	DefaultRule string  `@Ident "}"`
}

func (w *Workflow) String() string {
	return fmt.Sprintf("Workflow: %s, Rules: %v", w.Name, w.Rule)
}

func (w *Workflow) Eval(p *Part) string {
	for _, r := range w.Rule {
		n := r.Eval(p)
		if n != "" {
			slog.Debug("Matched workflow rule", "rule", r, "part", p, "next", n)
			return n
		}
	}
	return w.DefaultRule
}

type Rule struct {
	Category    string `@( "x" | "m" | "a" | "s")`
	Comparator  string `@(">"|"<")`
	Value       int    `@Int ":"`
	Destination string `@Ident ","`
}

func (r *Rule) String() string {
	return fmt.Sprintf("Rule: %s %s %d -> %s", r.Category, r.Comparator, r.Value, r.Destination)
}

func (r *Rule) Eval(p *Part) string {
	var pVal int
	switch r.Category {
	case "x":
		pVal = p.XRating
	case "m":
		pVal = p.MRating
	case "a":
		pVal = p.ARating
	case "s":
		pVal = p.SRating
	}

	switch r.Comparator {
	case ">":
		if pVal > r.Value {
			return r.Destination
		}
	case "<":
		if pVal < r.Value {
			return r.Destination
		}
	}

	return ""
}

type Flower struct {
	orderedWorkflows []*Workflow
	mappedWorkflows  map[string]*Workflow
}

func (f *Flower) String() string {
	return fmt.Sprintf("Flower: %v", f.orderedWorkflows)
}

func (f *Flower) AddWorkflow(w *Workflow) {
	f.orderedWorkflows = append(f.orderedWorkflows, w)
	f.mappedWorkflows[w.Name] = w
}

func (f *Flower) Process(p *Part) bool {
	w := f.mappedWorkflows["in"]
	for {
		n := w.Eval(p)
		switch n {
		case "A":
			return true
		case "R":
			return false
		default:
			w = f.mappedWorkflows[n]
		}
	}
}

type Part struct {
	XRating int `"{" "x" "=" @Int`
	MRating int `"," "m" "=" @Int`
	ARating int `"," "a" "=" @Int`
	SRating int `"," "s" "=" @Int "}"`
}

func (p *Part) String() string {
	return fmt.Sprintf("Part: %d %d %d %d", p.XRating, p.MRating, p.ARating, p.SRating)
}

func (p *Part) TotalRating() int {
	return p.XRating + p.MRating + p.ARating + p.SRating
}

func ParseCommands(puzzleFile string) (*Flower, []*Part) {
	lines := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")

	myLexer := lexer.MustSimple([]lexer.SimpleRule{
		// Order matters here! Int kept stealing the leading cards before I changed the ordering.
		{"Ident", `[a-zAR]+`},
		{"Int", `\d+`},
		{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
	})
	workflowParser := participle.MustBuild[Workflow](
		participle.Lexer(myLexer),
	)
	partParser := participle.MustBuild[Part](
		participle.Lexer(myLexer),
	)

	flower := &Flower{
		orderedWorkflows: []*Workflow{},
		mappedWorkflows:  map[string]*Workflow{},
	}
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
				slog.Error("failed to parse", "l", l, "w", w, "err", err)
				panic(err)
			}
			flower.AddWorkflow(w)
			continue
		}

		p, err := partParser.ParseBytes("", []byte(l))
		if err != nil {
			slog.Error("failed to parse", "l", l, "p", p, "err", err)
			panic(err)
		}
		parts = append(parts, p)
	}

	return flower, parts
}

func partOne(puzzleFile string) {
	slog.Info("Day Nineteen part one", "puzzle file", puzzleFile)

	flower, parts := ParseCommands(puzzleFile)

	slog.Debug("Parsed", "workflows", flower, "parts", parts)

	acceptedRatings := 0
	for _, p := range parts {
		if flower.Process(p) {
			acceptedRatings += p.TotalRating()
		}
	}

	slog.Info("Evaluated all parts", "accepted ratings", acceptedRatings)
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
