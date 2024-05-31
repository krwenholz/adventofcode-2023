package daySeven

import (
	"adventofcode/cmd/scanner"
	"container/heap"
	"fmt"
	"log"
	"log/slog"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/spf13/cobra"
)

type Kind int

const (
	Unknown Kind = iota
	FiveOfAKind
	FourOfAKind
	FullHouse
	ThreeOfAKind
	TwoPair
	OnePair
	HighCard
)

func (k Kind) String() string {
	switch k {
	case Unknown:
		return "Unknown"
	case FiveOfAKind:
		return "FiveOfAKind"
	case FourOfAKind:
		return "FourOfAKind"
	case FullHouse:
		return "FullHouse"
	case ThreeOfAKind:
		return "ThreeOfAKind"
	case TwoPair:
		return "TwoPair"
	case OnePair:
		return "OnePair"
	case HighCard:
		return "HighCard"
	default:
		return ""
	}
}

var cardOrdering = map[byte]int{
	'A': 14,
	'K': 13,
	'Q': 12,
	'J': 11,
	'T': 10,
	'9': 9,
	'8': 8,
	'7': 7,
	'6': 6,
	'5': 5,
	'4': 4,
	'3': 3,
	'2': 2,
}

type Hand struct {
	Cards string `@Cards`
	Bid   int    `@Int`
	kind  Kind
}

func complexHandKind(cardCounts map[rune]int) Kind {
	fourOfAKind := 0
	threeOfAKind := 0
	twoOfAKind := 0
	for _, count := range cardCounts {
		switch count {
		case 4:
			fourOfAKind++
		case 3:
			threeOfAKind++
		case 2:
			twoOfAKind++
		}
	}

	if fourOfAKind > 0 {
		return FourOfAKind
	} else if threeOfAKind > 0 && twoOfAKind > 0 {
		return FullHouse
	} else if threeOfAKind > 0 {
		return ThreeOfAKind
	}

	return TwoPair
}

func (h *Hand) Kind() Kind {
	if h.kind != Unknown {
		return h.kind
	}
	cardCounts := map[rune]int{}
	for _, card := range h.Cards {
		if _, ok := cardCounts[card]; !ok {
			cardCounts[card] = 0
		}
		cardCounts[card] = cardCounts[card] + 1
	}
	switch len(cardCounts) {
	case 1:
		h.kind = FiveOfAKind
	case 2:
		h.kind = complexHandKind(cardCounts)
	case 3:
		h.kind = complexHandKind(cardCounts)
	case 4:
		h.kind = OnePair
	default:
		h.kind = HighCard
	}
	return h.kind
}

func (h *Hand) Less(other *Hand) bool {
	if h.Kind() == other.Kind() {
		for i := range h.Cards {
			hOrder := cardOrdering[h.Cards[i]]
			otherOrder := cardOrdering[other.Cards[i]]
			if hOrder == otherOrder {
				continue
			}

			return hOrder < otherOrder
		}
	}
	// kind values are in reverse order
	if h.Kind() < other.Kind() {
		return false
	}
	return true
}

func (h *Hand) String() string {
	return fmt.Sprintf("Hand{Cards: %s, Bid: %d, Kind: %s}", h.Cards, h.Bid, h.Kind())
}

type HandHeap []*Hand

func (h HandHeap) Len() int { return len(h) }
func (h HandHeap) Less(i, j int) bool {
	slog.Debug("Comparing hands", "i", h[i], "j", h[j], "result", h[i].Less(h[j]))
	return h[i].Less(h[j])
}
func (h HandHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *HandHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*Hand))
}

func (h *HandHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func partOne(puzzleFile string) {
	fmt.Println("Day template part one", puzzleFile)
	handLexer := lexer.MustSimple([]lexer.SimpleRule{
		// Order matters here! Int kept stealing the leading cards before I changed the ordering.
		{"Cards", `[AKQJT98765432]{5}`},
		{"Int", `(\d*\.)?\d+`},
		{"EOL", `\n`},
		{"Colon", `:`},
		{"Whitespace", `[ \t]+`},
	})
	parser, err := participle.Build[Hand](
		participle.Lexer(handLexer),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		log.Fatal(err)
	}

	scanner := scanner.NewScanner[Hand](parser, puzzleFile)

	handHeap := &HandHeap{}
	for scanner.Scan() {
		hand := scanner.Struct()
		heap.Push(handHeap, hand)
	}

	ordered := []*Hand{}
	totalWinnings := 0
	rank := 0
	for handHeap.Len() > 0 {
		rank++
		h := heap.Pop(handHeap).(*Hand)
		totalWinnings += h.Bid * rank
		ordered = append(ordered, h)
	}

	slog.Debug("finished computing!", "ordered hands", ordered)
	slog.Info("Day seven part one", "total winnings", totalWinnings)
}

func partTwo(puzzleFile string) {
	fmt.Println("Day part two", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "daySeven",
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
	Cmd.Flags().BoolP("part-two", "p", false, "Whether to run part two of the day's challenge")
}
