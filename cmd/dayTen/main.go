package dayTen

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

type Pipe int

const (
	Vertical Pipe = iota
	Horizontal
	NinetyDegreeNorthEast
	NinetyDegreeNorthWest
	NinetyDegreeSouthWest
	NinetyDegreeSouthEast
	Ground
	Start
)

func (k Pipe) String() string {
	switch k {
	case Vertical:
		return "|"
	case Horizontal:
		return "-"
	case NinetyDegreeNorthEast:
		return "L"
	case NinetyDegreeNorthWest:
		return "J"
	case NinetyDegreeSouthWest:
		return "7"
	case NinetyDegreeSouthEast:
		return "F"
	case Ground:
		return "."
	case Start:
		return "S"
	default:
		return ""
	}
}

func NewPipe(r rune) Pipe {
	switch r {
	case '|':
		return Vertical
	case '-':
		return Horizontal
	case 'L':
		return NinetyDegreeNorthEast
	case 'J':
		return NinetyDegreeNorthWest
	case '7':
		return NinetyDegreeSouthWest
	case 'F':
		return NinetyDegreeSouthEast
	case '.':
		return Ground
	case 'S':
		return Start
	default:
		panic("invalid pipe type")
	}
}

type Grid struct {
	rows                    [][]*Position
	StartX, StartY          int
	MaxDistance, MaxX, MaxY int
}

func (g *Grid) Get(x, y int) *Position {
	// Add a safety check here so we don't have to everywhere else
	if x < 0 || y < 0 || x >= len(g.rows[0]) || y >= len(g.rows) {
		return nil
	}

	return g.rows[y][x]
}

func (g *Grid) String() string {
	str := ""
	for _, row := range g.rows {
		for _, p := range row {
			str += p.String() + ", "
		}
		str += "\n"
	}
	return str
}

func (g *Grid) DistancesString() string {
	str := ""
	for _, row := range g.rows {
		for _, p := range row {
			str += fmt.Sprint(p.DistanceFromStart) + " "
		}
		str += "\n"
	}
	return str
}

type Position struct {
	X, Y              int
	Type              Pipe
	DistanceFromStart int
}

func (p *Position) String() string {
	return fmt.Sprintf("Position{X: %d, Y: %d, Type: %s, Dist: %d}", p.X, p.Y, p.Type, p.DistanceFromStart)
}

func (p *Position) SetDistanceFromStart(d int) {
	p.DistanceFromStart = d
}

func CanNorth(p *Position) bool {
	if p == nil {
		return false
	}
	return p.Type == Start || p.Type == Vertical || p.Type == NinetyDegreeNorthEast || p.Type == NinetyDegreeNorthWest
}

func CanSouth(p *Position) bool {
	if p == nil {
		return false
	}
	return p.Type == Start || p.Type == Vertical || p.Type == NinetyDegreeSouthEast || p.Type == NinetyDegreeSouthWest
}

func CanEast(p *Position) bool {
	if p == nil {
		return false
	}
	return p.Type == Start || p.Type == Horizontal || p.Type == NinetyDegreeSouthEast || p.Type == NinetyDegreeNorthEast
}

func CanWest(p *Position) bool {
	if p == nil {
		return false
	}
	return p.Type == Start || p.Type == Horizontal || p.Type == NinetyDegreeSouthWest || p.Type == NinetyDegreeNorthWest
}

func (p *Position) Connections(g *Grid) []*Position {
	connections := []*Position{}
	var north, south, east, west *Position
	switch p.Type {
	case Vertical:
		north = g.Get(p.X, p.Y-1)
		south = g.Get(p.X, p.Y+1)
	case Horizontal:
		west = g.Get(p.X-1, p.Y)
		east = g.Get(p.X+1, p.Y)
	case NinetyDegreeNorthEast:
		north = g.Get(p.X, p.Y-1)
		east = g.Get(p.X+1, p.Y)
	case NinetyDegreeNorthWest:
		north = g.Get(p.X, p.Y-1)
		west = g.Get(p.X-1, p.Y)
	case NinetyDegreeSouthWest:
		south = g.Get(p.X, p.Y+1)
		west = g.Get(p.X-1, p.Y)
	case NinetyDegreeSouthEast:
		south = g.Get(p.X, p.Y+1)
		east = g.Get(p.X+1, p.Y)
	case Start:
		west = g.Get(p.X-1, p.Y)
		east = g.Get(p.X+1, p.Y)
		north = g.Get(p.X, p.Y-1)
		south = g.Get(p.X, p.Y+1)
	case Ground:
	}

	if CanNorth(south) {
		connections = append(connections, south)
	}
	if CanSouth(north) {
		connections = append(connections, north)
	}
	if CanEast(west) {
		connections = append(connections, west)
	}
	if CanWest(east) {
		connections = append(connections, east)
	}

	ret := []*Position{}
	for _, p := range connections {
		if p != nil && p.Type != Ground {
			ret = append(ret, p)
		}
	}
	slog.Debug("Connections", "p", p, "connections", connections, "ret", ret)
	return ret
}

func parse(path string) *Grid {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)

	g := &Grid{
		rows: [][]*Position{},
	}

	y := 0
	for scanner.Scan() {
		text := scanner.Text()
		row := []*Position{}
		for x, r := range text {
			newPos := &Position{
				X:                 x,
				Y:                 y,
				Type:              NewPipe(r),
				DistanceFromStart: -1,
			}
			row = append(row, newPos)
			if newPos.Type == Start {
				g.StartX = x
				g.StartY = y
				newPos.DistanceFromStart = 0
			}
		}
		g.rows = append(g.rows, row)
		y++
	}

	return g
}

func (g *Grid) GetStart() *Position {
	return g.Get(g.StartX, g.StartY)
}

func calculateDistance(g *Grid, p *Position, distanceFromStart int) int {
	// depth first traversal of each child node until we find start, keep track of distance so far
	slog.Debug("Calculating distance", "p", p, "distanceFromStart", distanceFromStart)
	if p.DistanceFromStart != -1 {
		return p.DistanceFromStart
	}

	// Start with our best guess, based on parent suggested distance
	p.DistanceFromStart = distanceFromStart

	// Fetch the distances for any valid neighbor
	neighborDistances := []int{}
	for _, conn := range p.Connections(g) {
		neighborDistances = append(neighborDistances, calculateDistance(g, conn, distanceFromStart+1))
	}

	// Check if our current (parent suggested) distance can be beat by traversing through
	// a neighbor
	for _, d := range neighborDistances {
		if d < p.DistanceFromStart {
			p.DistanceFromStart = d + 1
		}
	}

	if g.MaxDistance < p.DistanceFromStart {
		g.MaxDistance = p.DistanceFromStart
		g.MaxX = p.X
		g.MaxY = p.Y
	}

	return p.DistanceFromStart
}

func partOne(puzzleFile string) {
	// dayTen.simple.input is 4
	// dayTen.complex.input is 8
	slog.Info("Day Ten part one", "puzzle file", puzzleFile)

	grid := parse(puzzleFile)
	for _, c := range grid.GetStart().Connections(grid) {
		calculateDistance(grid, c, 1)
	}

	slog.Debug("distance calculated", "grid", grid.String())

	slog.Info("Day Ten part one", "max distance", grid.MaxDistance, "max x", grid.MaxX, "max y", grid.MaxY)
}

func partTwo(puzzleFile string) {
	slog.Info("Day Ten part two", "puzzle file", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "dayTen",
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
