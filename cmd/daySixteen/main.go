package daySixteen

import (
	"adventofcode/cmd/fileReader"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

/**
. is empty space
/ and \ are mirrors
| and - are splitters

beam enters top left (0, 0) going right
- empty spaces continue in same direction
- mirrors reflect 90 degrees
  - like you expect, let's do the vectors later
- pointy ends of splitters pass through
- flat side of splitters create two beams going two directions of pointy ends
- beams don't interact
- energized = at least one beam passes through, reflects, or splits in that space
**/

type Beam struct {
	id     int
	x      int
	y      int
	deltaX int
	deltaY int
}

func (b *Beam) String() string {
	return fmt.Sprintf("%d: %s %s", b.id, b.LocationString(), b.DirectionString())
}

func (b *Beam) Hash() string {
	return fmt.Sprintf("%s %s", b.LocationString(), b.DirectionString())
}

func (b *Beam) LocationString() string {
	return fmt.Sprintf("%d,%d", b.x, b.y)
}

func (b *Beam) DirectionString() string {
	switch b.deltaX {
	case 0:
		if b.deltaY > 0 {
			return "v"
		} else {
			return "^"
		}
	case -1:
		return "<"
	case 1:
		return ">"
	default:
		panic("WTF")
	}
}

func MirrorForwardSlash(deltaX, deltaY int) (int, int) {
	switch deltaX {
	case 0:
		if deltaY > 0 {
			return -1, 0
		} else {
			return 1, 0
		}
	case -1:
		return 0, 1
	case 1:
		return 0, -1
	default:
		panic("WTF")
	}
}

func MirrorBackwardSlash(deltaX, deltaY int) (int, int) {
	// Remember the top is 0,0 and incrementing y goes _down_
	switch deltaX {
	case 0:
		if deltaY > 0 {
			return 1, 0
		} else {
			return -1, 0
		}
	case -1:
		return 0, -1
	case 1:
		return 0, 1
	default:
		panic("WTF")
	}
}

func (b *Beam) StepBeam(rows []string) []*Beam {
	nextBeams := []*Beam{}

	switch rows[b.y][b.x] {
	case '.':
		nextBeams = append(
			nextBeams,
			&Beam{
				b.id,
				b.x + b.deltaX,
				b.y + b.deltaY,
				b.deltaX,
				b.deltaY,
			},
		)
	case '/':
		deltaX, deltaY := MirrorForwardSlash(b.deltaX, b.deltaY)
		nextBeams = append(
			nextBeams,
			&Beam{
				b.id + 1,
				b.x + deltaX,
				b.y + deltaY,
				deltaX,
				deltaY,
			},
		)
	case '\\':
		deltaX, deltaY := MirrorBackwardSlash(b.deltaX, b.deltaY)
		nextBeams = append(
			nextBeams,
			&Beam{
				b.id + 1,
				b.x + deltaX,
				b.y + deltaY,
				deltaX,
				deltaY,
			},
		)
	case '|':
		// only matters if moving in deltaX
		if b.deltaX == 0 {
			nextBeams = append(
				nextBeams,
				&Beam{
					b.id,
					b.x + b.deltaX,
					b.y + b.deltaY,
					b.deltaX,
					b.deltaY,
				},
			)
		} else {
			nextBeams = append(
				nextBeams,
				&Beam{
					b.id + 1,
					b.x,
					b.y + 1,
					0,
					1,
				},
				&Beam{
					b.id + 2,
					b.x,
					b.y - 1,
					0,
					-1,
				},
			)
		}
	case '-':
		// only matters if moving in deltaY
		if b.deltaY == 0 {
			nextBeams = append(
				nextBeams,
				&Beam{
					b.id,
					b.x + b.deltaX,
					b.y + b.deltaY,
					b.deltaX,
					b.deltaY,
				},
			)
		} else {
			nextBeams = append(
				nextBeams,
				&Beam{
					b.id + 1,
					b.x + 1,
					b.y,
					1,
					0,
				},
				&Beam{
					b.id + 2,
					b.x - 1,
					b.y,
					-1,
					0,
				},
			)

		}
	default:
		panic(fmt.Sprintf("WTF! Unknown character %s", string(rows[b.y][b.x])))
	}

	//slog.Debug("next beams", "nextBeams", nextBeams)
	filteredBeams := []*Beam{}
	for _, nextBeam := range nextBeams {
		if nextBeam.inBounds(rows) {
			filteredBeams = append(filteredBeams, nextBeam)
		}
	}
	return filteredBeams
}

func (b *Beam) inBounds(rows []string) bool {
	return b.y >= 0 && b.y < len(rows) && b.x >= 0 && b.x < len(rows[0])
}

func PrintTouches(rows []string, touchedSpaces map[string]bool) {
	if strings.ToLower(os.Getenv("LOG_LEVEL")) != "debug" {
		return
	}
	fmt.Println("##### touches ######")
	for ts := range touchedSpaces {
		split := strings.Split(ts, ",")
		x, _ := strconv.Atoi(split[0])
		y, _ := strconv.Atoi(split[1])

		rows[y] = rows[y][:x] + "#" + rows[y][x+1:]
	}

	for _, row := range rows {
		fmt.Println(row)
	}
}

func calculateEnergy(rows []string, startBeam *Beam, maxIterations int) int {
	touchedSpaces := map[string]bool{}
	beams := []*Beam{startBeam}
	seenBeams := map[string]bool{} // ???: should I memoize this across runs?

	steps := 0
	for len(beams) > 0 {
		//slog.Debug("running beams", "beams", beams)
		nextBeams := []*Beam{}

		for _, beam := range beams {
			touchedSpaces[beam.LocationString()] = true
			news := beam.StepBeam(rows)
			for _, newBeam := range news {
				if _, ok := seenBeams[newBeam.Hash()]; ok {
					continue
				}
				seenBeams[newBeam.Hash()] = true
				nextBeams = append(nextBeams, newBeam)
			}
		}

		beams = nextBeams
		if steps > maxIterations {
			slog.Debug("breaking because max iterations", "start", startBeam, "steps", steps, "maxIterations", maxIterations)
			break
		}
		steps++
	}

	//PrintTouches(rows, touchedSpaces)
	slog.Debug("calculated energy", "start", startBeam, "touched spaces", touchedSpaces, "steps", steps, "energized spaces", len(touchedSpaces))
	return len(touchedSpaces)
}

func partOne(puzzleFile string, maxIterations int) {
	slog.Info("Day Sixteen part one", "puzzle file", puzzleFile)
	rows := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")

	energizedSpaces := calculateEnergy(rows, &Beam{0, 0, 0, 1, 0}, maxIterations)

	slog.Info("Day Sixteen part one", "energized spaces", energizedSpaces)
}

func partTwo(puzzleFile string, maxIterations int) {
	slog.Info("Day Sixteen part two", "puzzle file", puzzleFile)
	rows := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")

	validStarts := []*Beam{}
	for x := 0; x < len(rows[0]); x++ {
		// top
		validStarts = append(validStarts, &Beam{0, x, 0, 0, 1})
		// bottom
		validStarts = append(validStarts, &Beam{0, x, len(rows) - 1, 0, -1})
	}
	for y := 0; y < len(rows); y++ {
		// left
		validStarts = append(validStarts, &Beam{0, 0, y, 1, 0})
		// right
		validStarts = append(validStarts, &Beam{0, len(rows[0]) - 1, y, -1, 0})
	}
	//validStarts = []*Beam{{0, 0, 8, 1, 0}}
	slog.Debug("all starts", "validStarts", validStarts)

	max := 0
	for _, start := range validStarts {
		energizedSpaces := calculateEnergy(rows, start, maxIterations)
		if energizedSpaces > max {
			max = energizedSpaces
		}
	}

	// need to generate all valid starts and then iterate
	slog.Info("Day Sixteen part two", "max energized spaces", max)
}

var Cmd = &cobra.Command{
	Use: "daySixteen",
	Run: func(cmd *cobra.Command, args []string) {
		puzzleInput, _ := cmd.Flags().GetString("puzzle-input")
		maxIterations, _ := cmd.Flags().GetInt("max-iterations")
		if !cmd.Flag("part-two").Changed {
			partOne(puzzleInput, maxIterations)
		} else {
			partTwo(puzzleInput, maxIterations)
		}
	},
}

func init() {
	Cmd.Flags().Bool("part-two", false, "Whether to run part two of the day's challenge")
	Cmd.Flags().Int("max-iterations", 100, "Max iterations to go through")
}
