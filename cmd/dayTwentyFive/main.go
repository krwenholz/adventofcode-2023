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

func findBridges(cs map[string]*Component, current string, visited map[string]bool, depth int) {
	disc := make([]int, len(cs))
	low := make([]int, len(cs))
	parent := map[string]string{}
	bridges := []string{}
	time := 0

	func dfs func(u int) {
		disc[u] = time
		low[u] = time
		time++

		for _, v := range cs[u].Connections {
			if disc[v] == 0 {
				parent[v] = cs[u].Label
				dfs(v)
				low[u] = min(low[u], low[v])
				if low[v] > disc[u] {
					bridges = append(bridges, u)
				}
			} else if v != parent[u] {
				low[u] = min(low[u], disc[v])
			}
		}
	}
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

	slog.Info("Finished Day TwentyFive part one", "expected", expected)
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
