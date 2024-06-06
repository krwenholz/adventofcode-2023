package dayEleven

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math"
	"os"

	"github.com/spf13/cobra"
	"gonum.org/v1/gonum/stat/combin"
)

const ONE_MILLION = 1000000

/**
build list of all galaxy coordinates
determine rows and columns with no galaxies using that list
for each galaxy pair in the list, calculate the distance between them
**/

type Galaxy struct {
	Id              int
	X, Y            int
	ClosestGalaxy   *Galaxy
	ClosestDistance int
}

func (g *Galaxy) String() string {
	j, _ := json.Marshal(g)
	return string(j)
}

type Observation struct {
	Height, Width int
	Galaxies      map[string]*Galaxy
	GalaxiesById  map[int]*Galaxy
	EmptyRows     map[int]bool
	EmptyCols     map[int]bool
}

func (o *Observation) String() string {
	j, _ := json.Marshal(o)
	return string(j)
}

func (o *Observation) PutGalaxy(g *Galaxy) {
	key := fmt.Sprintf("%d,%d", g.X, g.Y)
	o.Galaxies[key] = g
	o.GalaxiesById[g.Id] = g
}

func (o *Observation) GetGalaxy(x, y int) *Galaxy {
	key := fmt.Sprintf("%d,%d", x, y)
	g, _ := o.Galaxies[key]
	return g
}

func (o *Observation) GetGalaxyById(id int) *Galaxy {
	g, _ := o.GalaxiesById[id]
	return g
}

func parse(path string) *Observation {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)

	if err != nil {
		slog.Error("failed to parse", "err", err)
		panic(err)
	}

	observation := &Observation{
		Galaxies:     map[string]*Galaxy{},
		GalaxiesById: map[int]*Galaxy{},
		EmptyRows:    map[int]bool{},
		EmptyCols:    map[int]bool{},
	}
	y := 0
	rowsWithGalaxies := map[int]bool{}
	colsWithGalaxies := map[int]bool{}
	for scanner.Scan() {
		for x, c := range scanner.Text() {
			observation.Width = x
			if c == '.' {
				continue
			}

			observation.PutGalaxy(&Galaxy{
				Id: len(observation.Galaxies),
				X:  x,
				Y:  y,
			})
			rowsWithGalaxies[y] = true
			colsWithGalaxies[x] = true
		}
		observation.Height = y
		y++
	}

	slog.Debug("rows with galaxies", "rows", rowsWithGalaxies)
	for row := 0; row <= observation.Height; row++ {
		if !rowsWithGalaxies[row] {
			observation.EmptyRows[row] = true
		}
	}
	slog.Debug("cols with galaxies", "cols", colsWithGalaxies)
	for col := 0; col <= observation.Width; col++ {
		if !colsWithGalaxies[col] {
			observation.EmptyCols[col] = true
		}
	}

	return observation
}

func partOne(puzzleFile string) {
	observation := parse(puzzleFile)
	os.WriteFile("inputs/dayElevenObservations.json", []byte(observation.String()), 0644)

	combinationIndices := combin.Combinations(len(observation.Galaxies), 2)
	sum := 0
	for _, indices := range combinationIndices {
		g1 := observation.GetGalaxyById(indices[0])
		g2 := observation.GetGalaxyById(indices[1])
		maxX := int(math.Max(float64(g1.X), float64(g2.X)))
		minX := int(math.Min(float64(g1.X), float64(g2.X)))
		maxY := int(math.Max(float64(g1.Y), float64(g2.Y)))
		minY := int(math.Min(float64(g1.Y), float64(g2.Y)))

		// Normal space (manhattan) distance
		distance := maxX - minX + maxY - minY
		slog.Debug("combination distance", "combination", indices, "distance", distance)
		// Add one for every empty row or column between the galaxies
		for row := minY + 1; row < maxY; row++ {
			if observation.EmptyRows[row] {
				distance++
			}
		}
		for col := minX + 1; col < maxX; col++ {
			if observation.EmptyCols[col] {
				distance++
			}
		}
		if distance < g1.ClosestDistance || g1.ClosestGalaxy == nil {
			g1.ClosestDistance = distance
			g1.ClosestGalaxy = g2
		}
		if distance < g2.ClosestDistance || g2.ClosestGalaxy == nil {
			g2.ClosestDistance = distance
			g2.ClosestGalaxy = g1
		}

		sum += distance
	}

	slog.Info("Day Eleven part one", "sum", sum)
}

func partTwo(puzzleFile string) {
	observation := parse(puzzleFile)
	os.WriteFile("inputs/dayElevenObservations.json", []byte(observation.String()), 0644)

	combinationIndices := combin.Combinations(len(observation.Galaxies), 2)
	sum := 0
	for _, indices := range combinationIndices {
		g1 := observation.GetGalaxyById(indices[0])
		g2 := observation.GetGalaxyById(indices[1])
		maxX := int(math.Max(float64(g1.X), float64(g2.X)))
		minX := int(math.Min(float64(g1.X), float64(g2.X)))
		maxY := int(math.Max(float64(g1.Y), float64(g2.Y)))
		minY := int(math.Min(float64(g1.Y), float64(g2.Y)))

		// Normal space (manhattan) distance
		distance := maxX - minX + maxY - minY
		slog.Debug("combination distance", "combination", indices, "distance", distance)
		// Add one for every empty row or column between the galaxies
		expandedDistances := 0
		for row := minY + 1; row < maxY; row++ {
			if observation.EmptyRows[row] {
				expandedDistances++
			}
		}
		for col := minX + 1; col < maxX; col++ {
			if observation.EmptyCols[col] {
				expandedDistances++
			}
		}

		distance = (distance - expandedDistances) + expandedDistances*ONE_MILLION

		if distance < g1.ClosestDistance || g1.ClosestGalaxy == nil {
			g1.ClosestDistance = distance
			g1.ClosestGalaxy = g2
		}
		if distance < g2.ClosestDistance || g2.ClosestGalaxy == nil {
			g2.ClosestDistance = distance
			g2.ClosestGalaxy = g1
		}

		sum += distance
	}

	slog.Info("Day Eleven part two", "sum", sum)
}

var Cmd = &cobra.Command{
	Use: "dayEleven",
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
