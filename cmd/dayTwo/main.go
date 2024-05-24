package dayTwo

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/participle/v2"
)

// Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
type Game struct {
	GameNumber int        `@("Game" @Int) ":"`
	Displays   []*Display `@@+(";")?`
}

type Display struct {
	Dice []*Dice `@@+(",")?`
}

type Dice struct {
	Count int    ` @Int`
	Color string ` @("blue" | "red" | "green")`
}

func PartOne(puzzleFile string) {
	fmt.Println("Day template part one", puzzleFile)
	f, err := os.Open(puzzleFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	parser, err := participle.Build[Game]()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		game, err := parser.ParseBytes(puzzleFile, scanner.Bytes())
		if err != nil {
			fmt.Println(scanner.Text())
			log.Fatal(game, err)
		}
		fmt.Println(game)
	}
}

func PartTwo(puzzleFile string) {
	fmt.Println("Day template part two", puzzleFile)
}
