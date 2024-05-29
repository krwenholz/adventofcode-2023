package dayTwo

import (
	"adventofcode/cmd/scanner"
	"encoding/json"
	"fmt"
	"log"

	"github.com/alecthomas/participle/v2"
	"github.com/spf13/cobra"
)

// Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
type Game struct {
	Id     int      `"Game" @Int ":"`
	Rounds []*Round `(@@(";")?)+`
}

type Round struct {
	Dice []*Dice `(@@(",")?)+`
}

type Dice struct {
	Count int    ` @Int`
	Color string ` @("blue" | "red" | "green")`
}

func (g *Game) String() string {
	s, err := json.Marshal(g)
	if err != nil {
		log.Fatal(err)
	}

	return string(s)
}

func (g *Game) isValid(maxBlue, maxRed, maxGreen int) bool {
	for _, disp := range g.Rounds {
		for _, d := range disp.Dice {
			switch d.Color {
			case "blue":
				if d.Count > maxBlue {
					return false
				}
			case "red":
				if d.Count > maxRed {
					return false
				}
			case "green":
				if d.Count > maxGreen {
					return false
				}
			}
		}
	}
	return true
}

func (g *Game) power() int {
	var minBlue, minRed, minGreen int
	for _, r := range g.Rounds {
		for _, d := range r.Dice {
			switch d.Color {
			case "blue":
				if d.Count > minBlue {
					minBlue = d.Count
				}
			case "red":
				if d.Count > minRed {
					minRed = d.Count
				}
			case "green":
				if d.Count > minGreen {
					minGreen = d.Count
				}
			}
		}
	}
	return minBlue * minRed * minGreen
}

func partOne(puzzleFile string) {
	fmt.Println("Day template part one", puzzleFile)
	parser, err := participle.Build[Game]()
	if err != nil {
		log.Fatal(err)
	}

	scanner := scanner.NewScanner[Game](parser, puzzleFile)

	validGames := []*Game{}
	cumValidIdSum := 0

	for scanner.Scan() {
		game := scanner.Struct()

		if game.isValid(14, 12, 13) {
			validGames = append(validGames, game)
			cumValidIdSum += game.Id
		}
	}

	//j, _ := json.MarshalIndent(validGames, "", "  ")
	//fmt.Println(string(j))
	fmt.Println(len(validGames), cumValidIdSum)
}

func partTwo(puzzleFile string) {
	parser, err := participle.Build[Game]()
	if err != nil {
		log.Fatal(err)
	}

	scanner := scanner.NewScanner[Game](parser, puzzleFile)

	cumPowers := 0

	for scanner.Scan() {
		cumPowers += scanner.Struct().power()
	}

	fmt.Println(cumPowers)
}

// dayTwoCmd represents the dayTwo command
var Cmd = &cobra.Command{
	Use: "dayTwo",
	Run: func(cmd *cobra.Command, args []string) {
		puzzleInput, _ := cmd.Flags().GetString("puzzle-input")
		if !cmd.Flag("part-two").Changed {
			partOne(puzzleInput)
		} else {
			partTwo(puzzleInput)
		}
	},
}
