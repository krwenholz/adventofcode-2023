package dayTwentyFive

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
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

func twoComponentsAdj(cs map[string][]string) (bool, int) {
	visited := map[string]bool{}

	var dfs func(u string)
	dfs = func(u string) {
		visited[u] = true
		for _, v := range cs[u] {
			if !visited[v] {
				dfs(v)
			}
		}
	}

	components := []int{}
	for u := range cs {
		if !visited[u] {
			dfs(u)
			if len(components) > 0 {
				components = append(components, len(visited)-components[0])
			} else {
				components = append(components, len(visited))
			}
			if len(components) > 2 {
				return false, 0
			}
		}
	}

	if len(components) < 2 {
		return false, 0
	}

	val := 1
	for _, c := range components {
		val *= c
	}
	slog.Debug("Components", "components", components, "val", val)
	return true, val
}

func generateMermaidDiagram(edges [][]string) {
	mermaidFile, err := os.Create("/tmp/mermaid_diagram.mmd")
	if err != nil {
		slog.Error("Error creating mermaid file", "error", err)
		return
	}
	defer mermaidFile.Close()

	_, _ = mermaidFile.WriteString("%%{init: {\"flowchart\": {\"defaultRenderer\": \"elk\"}} }%%\n")
	_, err = mermaidFile.WriteString("flowchart LR;\n")
	if err != nil {
		slog.Error("Error writing to mermaid file", "error", err)
		return
	}

	for _, e := range edges {
		_, err = mermaidFile.WriteString(
			fmt.Sprintf(
				"%s <--> %s\n",
				e[0],
				e[1],
			))
		if err != nil {
			slog.Error("Error writing to mermaid file", "error", err)
			return
		}
	}
	if os.Getenv("LOG_LEVEL") == "debug" || os.Getenv("LOG_LEVEL") == "DEBUG" {
		cmd := exec.Command("mmdc", "-i", "/tmp/mermaid_diagram.mmd", "-o", "/tmp/mermaid_diagram.png")
		output, err := cmd.CombinedOutput()
		if err != nil {
			slog.Error("Error running mmdc command", "error", err, "output", string(output))
			return
		}
		slog.Info("Mermaid diagram generated", "output", string(output), "location", "/tmp/mermaid_diagram.png")
	}
}

func removeEdge(u0, u1 string, cs map[string][]string) {
	cons := []string{}
	for _, v := range cs[u0] {
		if v == u1 {
			continue
		}
		cons = append(cons, v)
	}
	cs[u0] = cons

	cons = []string{}
	for _, v := range cs[u1] {
		if v == u0 {
			continue
		}
		cons = append(cons, v)
	}
	cs[u1] = cons
}

func findComponentizingBridgesAdj(cs map[string][]string) ([][]string, int) {
	edges := [][]string{}
	slog.Debug("Components", "len(cs)", len(cs))

	for u0, cons := range cs {
		for _, u1 := range cons {
			if u0 > u1 {
				// Avoid duplicates
				continue
			}
			edges = append(edges, []string{u0, u1})
		}
	}
	generateMermaidDiagram(edges)
	slog.Debug("Edges", "len(edges)", len(edges))

	for i := range edges {
		e0 := edges[i]

		for j := i + 1; j < len(edges); j++ {
			e1 := edges[j]

			for k := j + 1; k < len(edges); k++ {
				e2 := edges[k]

				// Remove edges and check if we have two components
				tmpGraph := map[string][]string{}
				for u := range cs {
					tmpGraph[u] = cs[u]
				}
				removeEdge(e0[0], e0[1], tmpGraph)
				removeEdge(e1[0], e1[1], tmpGraph)
				removeEdge(e2[0], e2[1], tmpGraph)

				done, val := twoComponentsAdj(tmpGraph)
				if done {
					return [][]string{e0, e1, e2}, val
				}
			}

			slog.Info("Progress", "i", i, "j", j, "len(edges)", len(edges))
		}
	}
	return [][]string{}, 0
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
	csAdj := map[string][]string{}
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.Split(line, ":")
		label := splits[0]
		if label == "" {
			slog.Error("Invalid line", "line", line)
			continue
		}

		if _, ok := cs[label]; !ok {
			c := &Component{
				Label: label,
			}
			cs[label] = c
			csAdj[label] = []string{}
		}
		c := cs[label]

		connections := strings.Split(splits[1], " ")
		for _, newConLabel := range connections {
			if newConLabel == "" {
				continue
			}
			c.Connections = append(c.Connections, newConLabel)
			csAdj[label] = append(csAdj[label], newConLabel)

			if conComponent, ok := cs[newConLabel]; ok {
				conComponent.Connections = append(conComponent.Connections, label)
				csAdj[newConLabel] = append(csAdj[newConLabel], label)
			} else {
				cs[newConLabel] = &Component{
					Label:       newConLabel,
					Connections: []string{label},
				}
				csAdj[newConLabel] = []string{label}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Error reading file", "error", err)
	}

	//var bridges [][]string
	//bridges = findBridges(cs)

	//bridges, val := findComponentizingBridges(cs)
	bridges, val := findComponentizingBridgesAdj(csAdj)

	slog.Info("Finished Day TwentyFive part one", "expected", expected, "bridges", bridges, "val", val)
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
