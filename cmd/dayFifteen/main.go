package dayFifteen

import (
	"adventofcode/cmd/fileReader"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

/**
HASH
0. start at 0
1. for each character
   a. determine the ASCII code
   b. increase cur by ASCII code
   c. set cur to cur*17
   d. set cur to cur % 256

**/

func partOne(puzzleFile string) {
	slog.Info("Day Fifteen part one", "puzzle file", puzzleFile)
	steps := strings.Split(
		strings.ReplaceAll(
			fileReader.ReadFileContents(puzzleFile),
			"\n",
			"",
		),
		",",
	)

	sum := 0
	for _, step := range steps {
		h := hash(step)
		slog.Debug("computed hash", "step", step, "hash", h)
		sum += h
	}

	slog.Info("sum", "sum", sum)
}

func hash(step string) int {
	h := 0
	for _, char := range step {
		ascii := int(char)
		h += ascii
		h *= 17
		h %= 256
	}
	return h
}

type Lens struct {
	Label       string
	FocalLength int
}

func (l *Lens) String() string {
	return fmt.Sprintf("[%s %d]", l.Label, l.FocalLength)
}

type Box struct {
	id            int
	lenses        map[string]int
	orderedLenses []*Lens
}

func (b *Box) String() string {
	return fmt.Sprintf("Box %d: %v", b.id, b.orderedLenses)
}

func (b *Box) focusingPowers() []int {
	ps := []int{}
	for i, l := range b.orderedLenses {
		ps = append(ps, (1+b.id)*(1+i)*l.FocalLength)
	}

	return ps
}

func (b *Box) addLens(l *Lens) {
	if i, ok := b.lenses[l.Label]; ok {
		b.orderedLenses[i] = l
		return
	}

	b.orderedLenses = append(b.orderedLenses, l)
	b.lenses[l.Label] = len(b.orderedLenses) - 1
}

func (b *Box) removeLens(label string) {
	if i, ok := b.lenses[label]; ok {
		b.orderedLenses = append(b.orderedLenses[:i], b.orderedLenses[i+1:]...)
		delete(b.lenses, label)

		for j := i; j < len(b.orderedLenses); j++ {
			b.lenses[b.orderedLenses[j].Label] = j
		}
	}
}

func partTwo(puzzleFile string) {
	slog.Info("Day Fifteen part two", "puzzle file", puzzleFile)
	/**
	steps are now
	- a sequence of letters for a label
	- hash of the seq is the box for the step
	- following chars are the operation, = or -

	- -> go to the box
	  - remove the lens with the label if present
	  - nothing if not present
	  - shift remaining lenses forward (so it's a queue?)

	= -> add the indicated lens to the box
	  - if label is already present, replace it
	  - else add label to the end of the lenses
	**/
	steps := strings.Split(
		strings.ReplaceAll(
			fileReader.ReadFileContents(puzzleFile),
			"\n",
			"",
		),
		",",
	)

	boxes := make(map[int]*Box)

	for _, step := range steps {
		removeSplits := strings.Split(step, "-")
		if len(removeSplits) > 1 {
			lensLabel := removeSplits[0]
			boxNumber := hash(lensLabel)

			if box, ok := boxes[boxNumber]; ok {
				box.removeLens(lensLabel)
			}
		}

		addSplits := strings.Split(step, "=")
		if len(addSplits) > 1 {
			lensLabel := addSplits[0]
			boxNumber := hash(lensLabel)
			focalLength, _ := strconv.Atoi(addSplits[1])

			if box, ok := boxes[boxNumber]; ok {
				box.addLens(&Lens{
					Label:       lensLabel,
					FocalLength: focalLength,
				})
			} else {
				box := &Box{
					id:            boxNumber,
					lenses:        make(map[string]int),
					orderedLenses: []*Lens{},
				}
				box.addLens(&Lens{
					Label:       lensLabel,
					FocalLength: focalLength,
				})
				boxes[boxNumber] = box
			}
		}
		slog.Debug("Added a step", "step", step, "boxes", boxes)
	}

	sum := 0
	for _, box := range boxes {
		ps := box.focusingPowers()
		slog.Debug("Box focusing powers", "box", box.id, "powers", ps)
		for _, p := range ps {
			sum += p
		}
	}

	slog.Info("Finished day fifteen part two", "sum", sum)
}

var Cmd = &cobra.Command{
	Use: "dayFifteen",
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
