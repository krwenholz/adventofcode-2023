package dayTwentyTwo

import (
	"adventofcode/cmd/fileReader"
	"adventofcode/cmd/util"
	"fmt"
	"log/slog"
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
	Coords           []*Coordinate `@@ Tilde @@`
	SupportingBricks []string
	SupportsBricks   []string
	Id               string
}

func (b *Brick) String() string {
	return fmt.Sprintf("%s(%v): supported by: %v supports: %v", b.Id, b.Coords, b.SupportingBricks, b.SupportsBricks)
}

func (b *Brick) FinishInit(id int) {
	b.Id = ""
	for id >= 0 {
		b.Id += string('A' + (id % 26))
		id -= 26
	}
	b.SupportingBricks = []string{}
	b.SupportsBricks = []string{}
	b.Order()
}

func (b *Brick) Order() {
	c1 := b.Coords[0]
	c2 := b.Coords[1]
	if c1.Z < c2.Z {
		b.Coords[0] = c2
		b.Coords[1] = c1
	}
}

func (b *Brick) TopZ() int {
	return b.Coords[0].Z
}

func (b *Brick) BottomZ() int {
	return b.Coords[1].Z
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

/*
todo: 399 was too low and 423 was too high
*/
func partOne(puzzleFile string) {
	slog.Info("Day TwentyTwo part one", "puzzle file", puzzleFile)

	bricks := ParseBricks(puzzleFile)

	// collapse by minimizing the z values of any given brick
	// the next lowest supporting brick is as far as we can fall
	// we iterate _upwards_ to lower the lowest blocks first
	bricksMapped := map[string]*Brick{}
	for i, b := range bricks {
		// first find any supporting brick
		for j := i - 1; j >= 0; j-- {
			b2 := bricks[j]
			if b2.Supports(b) {
				// Only consider bricks that are at most 1 unit apart after the first support
				if len(b.SupportingBricks) != 0 && b.BottomZ()-b2.TopZ() > 1 {
					break
				}

				b.SupportingBricks = append(b.SupportingBricks, b2.Id)
				b2.SupportsBricks = append(b2.SupportsBricks, b.Id)
				if len(b.SupportingBricks) == 2 {
					break
				}
				if len(b.SupportingBricks) == 1 {
					slog.Debug("Falling", "b", b, "b2", b2)
					b.Fall(b.BottomZ() - b2.TopZ() - 1)
				}
			}

		}

		// unsupported bricks fall to the ground
		if len(b.SupportingBricks) == 0 {
			slog.Debug("Falling to ground", "b", b)
			b.Fall(b.BottomZ() - 1)
		}
		bricksMapped[b.Id] = b
	}

	disintegrable := []string{}
	for _, b := range bricks {
		canDisintegrate := true
		for _, supportedId := range b.SupportsBricks {
			if len(bricksMapped[supportedId].SupportingBricks) != 2 {
				canDisintegrate = false
				break
			}
		}

		if canDisintegrate {
			disintegrable = append(disintegrable, b.Id)
		}
	}

	slog.Debug("Bricks", "bricks", bricks, "disintegrable", disintegrable)
	slog.Info("Disintegrable", "disintegrable", len(disintegrable))
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
