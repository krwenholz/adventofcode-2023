package dayEight

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/spf13/cobra"
)

type Node struct {
	Name  string `@Node " = " "("`
	Left  string `@Node ", "`
	Right string `@Node ")" EOL?`
}

func (n *Node) String() string {
	return fmt.Sprintf("%s = (%s, %s)", n.Name, n.Left, n.Right)
}

func parse(path string) (string, map[string]*Node) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)

	nodeLexer := lexer.MustSimple([]lexer.SimpleRule{
		// Order matters here! Int kept stealing the leading cards before I changed the ordering.
		{"Node", `[A-Z0-9]{3}`},
		{"EOL", `\n`},
		{"Colon", `:`},
		{"Whitespace", `[ \t]+`},
	})
	parser, err := participle.Build[Node](
		participle.Lexer(nodeLexer),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		slog.Error("failed to parse", "err", err)
		panic(err)
	}

	instructions := scanner.Text()

	nodes := map[string]*Node{}
	for scanner.Scan() {
		node, err := parser.ParseString(path, scanner.Text())
		if err != nil {
			slog.Error("failed to parse", "err", err)
			panic(err)
		}
		nodes[node.Name] = node
	}

	return instructions, nodes
}

func partOne(puzzleFile string) {
	instructions, nodes := parse(puzzleFile)
	slog.Debug("parsed input", "instructions", instructions, "nodes", nodes)

	cur := nodes["AAA"]
	steps := 0
	for cur.Name != "ZZZ" {
		switch instructions[steps%len(instructions)] {
		case 'L':
			cur = nodes[cur.Left]
		case 'R':
			cur = nodes[cur.Right]
		}
		steps++
		slog.Debug("stepping", "cur", cur, "steps", steps)
	}
	slog.Info("Day eight part one", "steps", steps)
}

func partTwo(puzzleFile string) {
	instructions, nodes := parse(puzzleFile)
	slog.Debug("parsed input", "input", instructions, "nodes", nodes)

	curs := []*Node{}
	for _, node := range nodes {
		if node.Name[2] == 'A' {
			curs = append(curs, node)
		}
	}
	steps := 0
	for {
		newCurs := []*Node{}
		for _, cur := range curs {
			switch instructions[steps%len(instructions)] {
			case 'L':
				cur = nodes[cur.Left]
			case 'R':
				cur = nodes[cur.Right]
			}
			newCurs = append(newCurs, cur)
		}
		curs = newCurs
		steps++
		slog.Debug("stepping", "curs", curs, "steps", steps)

		allZs := true
		for _, cur := range curs {
			if cur.Name[2] != 'Z' {
				allZs = false
				break
			}
		}
		if allZs {
			break
		}
	}
	slog.Info("Day eight part one", "steps", steps)
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
