package dayEight

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/spf13/cobra"
)

type ParsedInput struct {
	Instructions string  `@Ident`
	Nodes        []*Node `@@*`
	MappedNodes  map[string]*Node
}

type Node struct {
	Name  string `@Ident "="`
	Left  string `"(" @Ident ","`
	Right string `@Ident ")"`
}

func (n *Node) String() string {
	return fmt.Sprintf("%s = (%s, %s)", n.Name, n.Left, n.Right)
}

func parse(puzzleFile string) *ParsedInput {
	f, err := os.ReadFile(puzzleFile)
	if err != nil {
		slog.Error("failed to parse", err)
		panic(err)
	}

	parser, err := participle.Build[ParsedInput]()
	if err != nil {
		slog.Error("failed to parse", "err", err)
		panic(err)
	}

	input, err := parser.ParseBytes("", f)
	if err != nil {
		slog.Error("failed to parse", "so far", input, "err", err)
		panic(err)
	}

	input.MappedNodes = make(map[string]*Node)
	for _, node := range input.Nodes {
		input.MappedNodes[node.Name] = node
	}

	return input
}

func partOne(puzzleFile string) {
	input := parse(puzzleFile)
	slog.Debug("parsed input", "input", input)

	cur := input.MappedNodes["AAA"]
	steps := 0
	for cur.Name != "ZZZ" {
		switch input.Instructions[steps%len(input.Instructions)] {
		case 'L':
			cur = input.MappedNodes[cur.Left]
		case 'R':
			cur = input.MappedNodes[cur.Right]
		}
		steps++
		slog.Debug("stepping", "cur", cur, "steps", steps)
	}
	slog.Info("Day eight part one", "steps", steps)
}

func partTwo(puzzleFile string) {
	fmt.Println("Day part two", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "dayEight",
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
