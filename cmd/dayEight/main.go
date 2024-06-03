package dayEight

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type Node struct {
	Name  string
	Left  string
	Right string
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

	if err != nil {
		slog.Error("failed to parse", "err", err)
		panic(err)
	}

	scanner.Scan()
	instructions := scanner.Text()
	scanner.Scan() // Skip the blank line
	slog.Debug("parsed instructions", "instructions", instructions)

	nodes := map[string]*Node{}
	for scanner.Scan() {
		node := &Node{
			Name:  scanner.Text()[:3],
			Left:  scanner.Text()[7:10],
			Right: scanner.Text()[12:15],
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

	curs := map[string]*Node{}
	zsVisited := map[string]map[string][]int{}
	for _, node := range nodes {
		if node.Name[2] == 'A' {
			curs[node.Name] = node
			zsVisited[node.Name] = map[string][]int{}
		}
	}
	slog.Info("initial curs", "curs", curs)

	steps := 0
	for {
		i := instructions[steps%len(instructions)]
		for origin, cur := range curs {
			switch i {
			case 'L':
				cur = nodes[cur.Left]
			case 'R':
				cur = nodes[cur.Right]
			}
			curs[origin] = cur
		}
		steps++

		numZs := 0
		for origin, cur := range curs {
			if cur.Name[2] == 'Z' {
				numZs++
				if _, ok := zsVisited[origin][cur.Name]; !ok {
					zsVisited[origin][cur.Name] = []int{}
				}
				zsVisited[origin][cur.Name] = append(zsVisited[origin][cur.Name], steps)
			}
		}
		slog.Debug("stepping", "curs", curs, "steps", steps, "numZs", numZs)
		/**if numZs > 1 {
			slog.Info("stepping on zs", "curs", curs, "steps", steps, "numZs", numZs)
		}**/

		if numZs == len(curs) || allCycled(zsVisited) {
			break
		}
	}
	j, _ := json.MarshalIndent(zsVisited, "", "  ")
	os.WriteFile(fmt.Sprintf("inputs/zsVisited%d.json", time.Now().Nanosecond()), j, 0644)

	slog.Info("Day eight part two", "steps", steps, "curs", curs)
}

// Turns out I had to look up the answer on Reddit. Once we're in a cycle, we can find the LCM of
// the cycle lengths to find the number of steps to get every ghost in sync.
func allCycled(zsVisited map[string]map[string][]int) bool {
	for _, zs := range zsVisited {
		hasCycle := false
		for _, visits := range zs {
			if len(visits) > 1 {
				hasCycle = true
			}
		}
		if !hasCycle {
			return false
		}
	}
	return true
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
