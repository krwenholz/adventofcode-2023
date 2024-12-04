package dayTwentyFive

import (
	"bufio"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Component struct {
	Label       string
	Connections []string
}

func findBridges(cs map[string]*Component) [][]string {
	visited := map[string]bool{}
	disc := map[string]int{}
	low := map[string]int{}
	parent := map[string]string{}
	bridges := [][]string{}
	time := 0

	var dfs func(u string)
	dfs = func(u string) {
		visited[u] = true

		disc[u] = time
		low[u] = time
		time++

		// Recurse for all adjacent vertices
		for _, v := range cs[u].Connections {
			// If v is not visited yet, we make it a child of u in DFS tree
			// and then recurse for it
			if !visited[v] {
				if disc[v] == 0 {
					parent[v] = u
					dfs(v)

					// Check if the subtree rooted with v has a connection to one of the ancestors of u
					low[u] = min(low[u], low[v])

					// If the lowest vertex reachable from subtree under v is below u in DFS tree, then u-v is a bridge
					if low[v] > disc[u] {
						bridges = append(bridges, []string{u, v})
					}
				} else if v != parent[u] {
					low[u] = min(low[u], disc[v])
				}
			} else if p, ok := parent[u]; ok && p != v {
				// Update low value of u for parent function calls.
				low[u] = min(low[u], disc[v])
			}
		}
	}

	for u := range cs {
		if !visited[u] {
			dfs(u)
		}
	}

	return bridges
}

func partOne(puzzleFile string) {
	slog.Info("Day TwentyFive part one", "puzzle file", puzzleFile)
	file, err := os.Open(puzzleFile)
	if err != nil {
		slog.Error("Error opening file", "error", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	expected := scanner.Text()

	cs := map[string]*Component{}
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.Split(line, ":")
		label := splits[0]

		if _, ok := cs[label]; !ok {
			c := &Component{
				Label: label,
			}
			cs[label] = c
		}
		c := cs[label]

		connections := strings.Split(splits[1], " ")
		for _, newConLabel := range connections {
			c.Connections = append(c.Connections, newConLabel)

			if conComponent, ok := cs[newConLabel]; ok {
				conComponent.Connections = append(conComponent.Connections, label)
			} else {
				cs[newConLabel] = &Component{
					Label:       newConLabel,
					Connections: []string{label},
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Error reading file", "error", err)
	}

	bridges := findBridges(cs)

	slog.Info("Finished Day TwentyFive part one", "expected", expected, "bridges", bridges)
}

func partTwo(puzzleFile string) {
	slog.Info("Day TwentyFive part two", "puzzle file", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "dayTwentyFive",
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
