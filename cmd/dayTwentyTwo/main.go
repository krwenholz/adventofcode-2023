package dayTwentyTwo

import (
	"adventofcode/cmd/fileReader"
	"adventofcode/cmd/util"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/spf13/cobra"
)

type Coordinate struct {
	X int `@Int Comma`
	Y int `@Int Comma`
	Z int `@Int`
}

func (c *Coordinate) String() string {
	return fmt.Sprintf("(%v, %v, %v)", c.X, c.Y, c.Z)
}

type Brick struct {
	Coords      []*Coordinate `@@ Tilde @@`
	SupportedBy []string
	Supporting  []string
	Id          string
}

func (b *Brick) String() string {
	return fmt.Sprintf("%s(%v): supported by: %v supports: %v", b.Id, b.Coords, b.SupportedBy, b.Supporting)
}

func (b *Brick) FinishInit(id int) {
	b.Id = ""
	for id >= 0 {
		b.Id += string('A' + (id % 26))
		id -= 26
	}
	b.SupportedBy = []string{}
	b.Supporting = []string{}
}

func (b *Brick) TopZ() int {
	return b.Coords[1].Z
}

func (b *Brick) BottomZ() int {
	return b.Coords[0].Z
}

func (b *Brick) Fall(deltaZ int) {
	for _, c := range b.Coords {
		c.Z -= deltaZ
	}
}

// Assume we have two line segments, s and t, return true if they overlap
func overlap(s1, s2, t1, t2 int) bool {
	sContains := s1 <= t1 && t2 <= s2
	tContains := t1 <= s1 && s2 <= t2
	overlapLeft := s1 <= t1 && t1 <= s2
	overlapRight := s1 <= t2 && t2 <= s2
	return sContains || tContains || overlapLeft || overlapRight
}

func overlapX(a1, a2, b1, b2 *Coordinate) bool {
	aL, aR := util.Order(a1.X, a2.X)
	bL, bR := util.Order(b1.X, b2.X)
	return overlap(aL, aR, bL, bR)
}

func overlapY(a1, a2, b1, b2 *Coordinate) bool {
	aL, aR := util.Order(a1.Y, a2.Y)
	bL, bR := util.Order(b1.Y, b2.Y)
	return overlap(aL, aR, bL, bR)
}

func (b *Brick) Supports(b2 *Brick) bool {
	if b.TopZ() >= b2.BottomZ() {
		return false
	}

	// I need to check if b has any overlap in the x,y plane with b2
	if !overlapX(b.Coords[0], b.Coords[1], b2.Coords[0], b2.Coords[1]) {
		return false
	}
	if !overlapY(b.Coords[0], b.Coords[1], b2.Coords[0], b2.Coords[1]) {
		return false
	}

	return true
}

func ParseBricks(puzzleFile string) []*Brick {
	lines := strings.Split(fileReader.ReadFileContents(puzzleFile), "\n")

	myLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Tilde", `~`},
		{"Comma", `,`},
		{"Int", `[0-9]+`},
	})
	brickParser := participle.MustBuild[Brick](
		participle.Lexer(myLexer),
	)
	bricks := []*Brick{}

	hasAnswerLine := len(lines[0]) < 4
	if hasAnswerLine {
		slog.Info("Expected answer", "count", lines[0])
		lines = lines[1:]
	}
	for i, l := range lines {
		b, err := brickParser.ParseBytes("", []byte(l))
		b.FinishInit(i)
		if err != nil {
			slog.Error("failed to parse", "l", l, "b", b, "err", err)
			panic(err)
		}

		bricks = append(bricks, b)
	}

	slices.SortFunc[[]*Brick](bricks, func(a, b *Brick) int {
		return a.BottomZ() - b.BottomZ()
	})
	return bricks
}

func applyGravity(bricks []*Brick) []*Brick {
	topOfBricksAtRest := map[int][]*Brick{}
	for _, b := range bricks {
		supported := false
		for b.BottomZ() > 1 {
			if bs, ok := topOfBricksAtRest[b.BottomZ()-1]; ok {
				for _, b2 := range bs {
					if b2.Supports(b) {
						supported = true
						break
					}
				}
			}
			if supported {
				break
			}

			b.Coords[0].Z--
			b.Coords[1].Z--
		}

		if _, ok := topOfBricksAtRest[b.TopZ()]; !ok {
			topOfBricksAtRest[b.TopZ()] = []*Brick{}
		}
		topOfBricksAtRest[b.TopZ()] = append(topOfBricksAtRest[b.TopZ()], b)

		/**
		for dec := i - 1; dec >= 0; dec-- {
			b2 := bricks[j]
			if b2.Supports(b) {
				unsupported = false

				slog.Debug("Falling", "b", b, "b2", b2)
				b.Fall(b.BottomZ() - b2.TopZ() - 1)

				break
			}
		}

		// unsupported bricks fall to the ground
		if unsupported {
			slog.Debug("Falling to ground", "b", b)
			b.Fall(b.BottomZ() - 1)
		}
			**/
	}

	slices.SortFunc[[]*Brick](bricks, func(a, b *Brick) int {
		return a.BottomZ() - b.BottomZ()
	})

	return bricks
}

func findSupports(bricks []*Brick) []*Brick {
	for i, b := range bricks {
		for j := i - 1; j >= 0; j-- {
			b2 := bricks[j]
			if b2.Supports(b) {
				if b.BottomZ()-b2.TopZ() > 1 {
					// Although we're sorted, we can't break because tall blocks may actually be further down
					// the list
					continue
				}

				b.SupportedBy = append(b.SupportedBy, b2.Id)
				b2.Supporting = append(b2.Supporting, b.Id)

				// We only care about being supported by _more_ than one brick, so exit early
				if len(b.SupportedBy) == 2 {
					break
				}
			}

		}
	}

	return bricks
}

/*
wrong answers:
- 399
- 423
- 439
- 452
- 460
- 407??
*/
func partOne(puzzleFile string) {
	slog.Info("Day TwentyTwo part one", "puzzle file", puzzleFile)

	bricks := ParseBricks(puzzleFile)
	bricks = applyGravity(bricks)
	bricks = findSupports(bricks)

	bricksMapped := map[string]*Brick{}
	for _, b := range bricks {
		bricksMapped[b.Id] = b
	}

	disintegrable := []string{}
	for _, b := range bricks {
		canDisintegrate := true
		for _, supportedId := range b.Supporting {
			if len(bricksMapped[supportedId].SupportedBy) != 2 {
				canDisintegrate = false
				break
			}
		}

		if canDisintegrate {
			disintegrable = append(disintegrable, b.Id)
		}
	}

	printBricks(bricks)
	slog.Debug("Disintegrable", "disintegrable", disintegrable)
	slog.Info("Disintegrable", "disintegrable", len(disintegrable))
}

func printBricks(bricks []*Brick) {
	if strings.ToLower(os.Getenv("LOG_LEVEL")) != "debug" {
		return
	}
	for i := len(bricks) - 1; i >= 0; i-- {
		b := bricks[i]
		fmt.Println(b)
	}
}

func partTwo(puzzleFile string) {
	slog.Info("Day TwentyTwo part two", "puzzle file", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "dayTwentyTwo",
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
