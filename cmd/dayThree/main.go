package dayThree

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"unicode"
)

type SchematicEntry struct {
	PartNumber int
	Symbol     string
	Xs         []int
	Y          int
}

func (s *SchematicEntry) IsSymbol() bool {
	return s.Symbol != "" && s.Symbol != "."
}

func (s *SchematicEntry) IsValidPart(schematic map[string]*SchematicEntry) bool {
	xs := append([]int{s.Xs[0] - 1}, s.Xs...)
	xs = append(xs, s.Xs[len(s.Xs)-1]+1)
	for _, x := range xs {
		for _, row := range []int{s.Y - 1, s.Y, s.Y + 1} {
			if e, ok := schematic[id(x, row)]; ok && e.IsSymbol() {
				return true
			}
		}
	}

	return false
}

func (s *SchematicEntry) String() string {
	return fmt.Sprintf("PartNumber: %d, Symbol: %s, Xs: %v, Y: %d", s.PartNumber, s.Symbol, s.Xs, s.Y)
}

func (s *SchematicEntry) Id() string {
	return id(s.Xs[0], s.Y)
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

func PartOne(puzzleFile string) {
	schematic := BuildSchematic(puzzleFile)
	slog.Debug("built schematic", "entries", schematic)

	//fmt.Println(parts)
	parts := map[string]*SchematicEntry{}
	for _, e := range schematic {
		if e.IsSymbol() {
			continue
		}
		if _, ok := parts[e.Id()]; !ok {
			if e.IsValidPart(schematic) {
				parts[e.Id()] = e
			}
		}
	}

	partsSum := 0
	for _, p := range parts {
		partsSum += p.PartNumber
	}
	slog.Info("final sum", "sum", partsSum)
}

func PartTwo(puzzleFile string) {
}
