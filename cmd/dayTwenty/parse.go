package dayTwenty

import (
	"adventofcode/cmd/fileReader"
	"log/slog"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

func ParseModules(puzzleFile string) map[string]*Module {
	lines := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")

	myLexer := lexer.MustSimple([]lexer.SimpleRule{
		// Order matters here! Int kept stealing the leading cards before I changed the ordering.
		{"Ident", `[a-zAR]+`},
		{"ModuleKind", `[%&]`},
		{"Pointer", ` -> `},
		{"ReceiverSeparator", `, `},
	})
	moduleParser := participle.MustBuild[Module](
		participle.Lexer(myLexer),
	)

	modules := map[string]*Module{}

	for _, l := range lines[1:] {
		m, err := moduleParser.ParseBytes("", []byte(l))
		if err != nil {
			slog.Error("failed to parse", "l", l, "m", m, "err", err)
			panic(err)
		}

		if m.ModuleKind == "%" {
			m.flipFlopState = false
		}
		modules[m.Name] = m
	}

	// Initialize all conjunction module input states
	cs := map[string]*Module{}
	for _, m := range modules {
		if m.ModuleKind == "&" {
			m.conjunctionState = map[string]PulseType{}
			cs[m.Name] = m
		}
	}
	for _, m := range modules {
		for _, r := range m.Receivers {
			if _, ok := cs[r]; ok {
				cs[r].conjunctionState[m.Name] = Low
			}
		}
	}

	return modules
}
