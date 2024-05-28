package dayFour

import (
	"adventofcode/cmd/scanner"
	"fmt"
	"log"
	"log/slog"

	"github.com/alecthomas/participle/v2"
)

// Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
type Card struct {
	Id               int   `"Card" @Int ":"`
	WinningNumbers   []int `@Int+ "|"`
	PossessedNumbers []int `@Int+`
}

func (c *Card) String() string {
	return fmt.Sprintf("Card %d: %v | %v", c.Id, c.WinningNumbers, c.PossessedNumbers)
}

func (c *Card) Points() int {
	winners := make(map[int]bool)
	for _, num := range c.WinningNumbers {
		winners[num] = true
	}

	points := 0
	for _, num := range c.PossessedNumbers {
		if _, ok := winners[num]; ok {
			if points == 0 {
				points = 1
			} else {
				points *= 2
			}
		}
	}

	slog.Debug("points calculation", "card", c, "points", points)
	return points
}

func (c *Card) Copies() int {
	winners := make(map[int]bool)
	for _, num := range c.WinningNumbers {
		winners[num] = true
	}

	count := 0
	for _, num := range c.PossessedNumbers {
		if _, ok := winners[num]; ok {
			count += 1
		}
	}

	slog.Debug("copies calculation", "card", c, "copies", count)
	return count
}

func PartOne(puzzleFile string) {
	fmt.Println("Day template part one", puzzleFile)
	parser, err := participle.Build[Card]()
	if err != nil {
		log.Fatal(err)
	}

	scanner := scanner.NewScanner[Card](parser, puzzleFile)
	points := 0
	for scanner.Scan() {
		points += scanner.Struct().Points()
	}

	slog.Info("Total points", "points", points)
}

func PartTwo(puzzleFile string) {
	fmt.Println("Day template part one", puzzleFile)
	parser, err := participle.Build[Card]()
	if err != nil {
		log.Fatal(err)
	}

	scanner := scanner.NewScanner[Card](parser, puzzleFile)
	// Store the number of copies granted to followers as we visit each
	// card. These are cumulative via multiplication.
	copyTally := map[int]int{}
	total := 0
	for scanner.Scan() {
		thisCard := scanner.Struct()
		countOfThisCard := 1
		if c, ok := copyTally[thisCard.Id]; ok {
			countOfThisCard += c
		}
		total += countOfThisCard

		slog.Debug(
			"Copying and counting cards",
			"card", thisCard,
			"countOfThisCard", countOfThisCard,
			"count", total,
			"copyTally", copyTally,
			"copies", thisCard.Copies(),
		)
		for i := 0; i < countOfThisCard; i++ {
			cardsToCopy := thisCard.Copies()
			for j := 1; j <= cardsToCopy; j++ {
				copyingCard := thisCard.Id + j
				if _, ok := copyTally[copyingCard]; !ok {
					copyTally[copyingCard] = 0
				}
				copyTally[copyingCard] = copyTally[copyingCard] + 1
			}
		}
	}

	slog.Debug("Total cards", "count", total, "copyTally", copyTally)
	slog.Info("Total cards", "count", total)
}
