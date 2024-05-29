package dayFive

import (
	"fmt"
	"log/slog"
	"math"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/spf13/cobra"
)

type SeedMap struct {
	Seeds      []int  `"seeds" (@Int)+ EOL EOL`
	Maps       []*Map `@@+`
	MappedMaps map[string]*Map
}

type Map struct {
	SrcType     string         `@Ident ToSeparator`
	DstType     string         `@Ident "map" EOL`
	MappedRange []*MappedRange `@@+ EOL*`
}

type MappedRange struct {
	Dst   int `@Int`
	Src   int `@Int`
	Range int `@Int EOL?`
}

func (m *Map) String() string {
	return fmt.Sprintf("%s %s %v", m.SrcType, m.DstType, m.MappedRange)
}

func (m *Map) Get(src int) int {
	for _, r := range m.MappedRange {
		if r.Src <= src && src < r.Src+r.Range {
			return r.Dst + (src - r.Src)
		}
	}
	return src
}

func (m *MappedRange) String() string {
	return fmt.Sprintf("%d -> %d (range %d)", m.Src, m.Dst, m.Range)
}

func parse(puzzleFile string) *SeedMap {
	f, err := os.ReadFile(puzzleFile)
	if err != nil {
		slog.Error("failed to parse", err)
		panic(err)
	}

	seedMapLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Int", `(\d*\.)?\d+`},
		{"Ident", `[a-zA-Z_]\w*`},
		{"EOL", `\n`},
		{"Colon", `:`},
		{"Whitespace", `[ \t]+`},
		{"ToSeparator", `-to-`},
	})
	parser, err := participle.Build[SeedMap](
		participle.Lexer(seedMapLexer),
		participle.Elide("Colon", "Whitespace"),
	)
	if err != nil {
		slog.Error("failed to parse", "err", err)
		panic(err)
	}

	seedMap, err := parser.ParseBytes("", f)
	if err != nil {
		slog.Error("failed to parse", "so far", seedMap, "err", err)
		panic(err)
	}
	seedMap.MappedMaps = make(map[string]*Map)

	slog.Debug("parsed", "seedMap", seedMap)

	for _, m := range seedMap.Maps {
		// Make nice maps
		seedMap.MappedMaps[m.SrcType] = m
	}

	return seedMap
}

func partOne(puzzleFile string) {
	seedMap := parse(puzzleFile)

	minimumLocation := math.MaxInt

	for _, s := range seedMap.Seeds {
		src := "seed"
		curNumber := s
		for {
			if dst, ok := seedMap.MappedMaps[src]; ok {
				curNumber = dst.Get(curNumber)
				src = dst.DstType
			} else {
				break
			}
		}

		if curNumber < minimumLocation {
			minimumLocation = curNumber
		}
	}

	slog.Info("Day five part one", "minimumLocation", minimumLocation)
}

func partTwo(puzzleFile string) {
	fmt.Println("Day part two", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "dayFive",
	Run: func(cmd *cobra.Command, args []string) {
		puzzleInput, _ := cmd.Flags().GetString("puzzle-input")
		if !cmd.Flag("part-two").Changed {
			partOne(puzzleInput)
		} else {
			partTwo(puzzleInput)
		}
	},
}
