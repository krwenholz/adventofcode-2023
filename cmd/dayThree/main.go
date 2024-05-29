package dayThree

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"unicode"

	"github.com/spf13/cobra"
)

type SchematicEntry struct {
	PartNumber  int
	Symbol      string
	Xs          []int
	Y           int
	isValidPart bool
	GearParts   map[string]int
}

func (e *SchematicEntry) IsSymbol() bool {
	return e.Symbol != "" && e.Symbol != "."
}

func (e *SchematicEntry) AddGearPart(p *SchematicEntry) {
	if _, ok := e.GearParts[e.Id()]; !ok {
		e.GearParts[p.Id()] = p.PartNumber
	}
}

func (s *SchematicEntry) IsValidPart(schematic map[string]*SchematicEntry) bool {
	if s.isValidPart {
		return true
	}
	if s.IsSymbol() {
		return false
	}
	xs := append([]int{s.Xs[0] - 1}, s.Xs...)
	xs = append(xs, s.Xs[len(s.Xs)-1]+1)
	for _, x := range xs {
		for _, row := range []int{s.Y - 1, s.Y, s.Y + 1} {
			if e, ok := schematic[id(x, row)]; ok && e.IsSymbol() {
				s.isValidPart = true
				return true
			}
		}
	}

	return false
}

func (s *SchematicEntry) IsGear(schematic map[string]*SchematicEntry) bool {
	if s.Symbol != "*" {
		return false
	}

	xs := append([]int{s.Xs[0] - 1}, s.Xs...)
	xs = append(xs, s.Xs[len(s.Xs)-1]+1)
	for _, x := range xs {
		for _, row := range []int{s.Y - 1, s.Y, s.Y + 1} {
			if e, ok := schematic[id(x, row)]; ok && e.IsValidPart(schematic) {
				slog.Debug("adding a new gear part", "s.Id", s.Id(), "gear part", e)
				s.AddGearPart(e)
			}
		}
	}

	slog.Debug(
		"calculated gear parts",
		"s.Id", s.Id(),
		"gear parts", s.GearParts,
		"length of gears", len(s.GearParts),
	)
	return len(s.GearParts) == 2
}

// NAUGHTY: this relies on IsGear being called first, should prepopulate gearParts at construction
// meh
func (e *SchematicEntry) GearRatio() int {
	ratio := 1
	for _, partNumber := range e.GearParts {
		ratio *= partNumber
	}
	return ratio
}

func (e *SchematicEntry) String() string {
	s, err := json.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}

	return string(s)
}

func (e *SchematicEntry) Id() string {
	return id(e.Xs[0], e.Y)
}

func id(x, y int) string {
	return fmt.Sprintf("%d, %d", x, y)
}

func isSpace(c rune) bool {
	return c == '.'
}

func BuildSchematic(path string) map[string]*SchematicEntry {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	rows := []string{}

	for scanner.Scan() {
		// Not my favorite, but I'm getting lazy today
		rows = append(rows, scanner.Text())
	}

	schematicEntries := []*SchematicEntry{}
	xs := []int{}

	for i, row := range rows {
		candidatePart := 0

		for j, c := range row {
			atEnd := j == len(row)-1

			if unicode.IsDigit(c) {
				candidatePart = candidatePart*10 + int(c-'0')
				xs = append(xs, j)

				if !atEnd {
					continue
				}
			}

			e := &SchematicEntry{
				PartNumber: 0,
				Symbol:     string(c),
				Xs:         []int{j},
				Y:          i,
				GearParts:  map[string]int{},
			}
			if e.IsSymbol() {
				schematicEntries = append(schematicEntries, e)
			}

			if candidatePart != 0 || atEnd {
				schematicEntries = append(
					schematicEntries,
					&SchematicEntry{
						PartNumber: candidatePart,
						Symbol:     "",
						Xs:         xs,
						Y:          i,
						GearParts:  map[string]int{},
					})

				candidatePart = 0
				xs = []int{}
			}
		}
	}

	schematic := map[string]*SchematicEntry{}
	for _, e := range schematicEntries {
		for _, x := range e.Xs {
			schematic[id(x, e.Y)] = e
		}
	}

	return schematic
}

func SumGearRatios(path string) []int {
	return nil
}

func partOne(puzzleFile string) {
	schematic := BuildSchematic(puzzleFile)
	slog.Debug("built schematic", "entries", schematic)

	parts := map[string]*SchematicEntry{}
	partsSum := 0
	for _, e := range schematic {
		if e.IsSymbol() {
			continue
		}
		if _, ok := parts[e.Id()]; !ok {
			if e.IsValidPart(schematic) {
				parts[e.Id()] = e
				partsSum += e.PartNumber
			}
		}
	}

	slog.Debug("final sum", "parts", parts, "sum", partsSum)
	slog.Info("final sum", "sum", partsSum)
}

func partTwo(puzzleFile string) {
	schematic := BuildSchematic(puzzleFile)
	slog.Debug("built schematic", "entries", schematic)

	gears := map[string]*SchematicEntry{}
	gearRatiosSum := 0
	for _, e := range schematic {
		if !e.IsGear(schematic) {
			continue
		}

		if _, ok := gears[e.Id()]; !ok {
			gears[e.Id()] = e
			gearRatiosSum += e.GearRatio()
		}
	}

	slog.Debug("final sum", "parts", gears, "sum", gearRatiosSum)
	slog.Info("final sum", "sum", gearRatiosSum)
}

var Cmd = &cobra.Command{
	Use: "dayThree",
	Run: func(cmd *cobra.Command, args []string) {
		puzzleInput, _ := cmd.Flags().GetString("puzzle-input")
		if !cmd.Flag("part-two").Changed {
			partOne(puzzleInput)
		} else {
			partTwo(puzzleInput)
		}
	},
}
