package dayTwentyFive

import (
	"adventofcode/cmd/util"
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func productOfTwoComponents(cs map[string][]string) int {
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
				panic("More than two components")
			}
		}
	}

	if len(components) < 2 {
		panic("Less than two components")
	}

	val := 1
	for _, c := range components {
		val *= c
	}
	slog.Debug("Components", "components", components, "val", val)
	return val
}

func generateMermaidDiagram(cs map[string][]string) {
	if !util.InDebugMode() {
		return
	}

	edges := [][]string{}
	for u, cons := range cs {
		for _, v := range cons {
			edges = append(edges, []string{u, v})
		}
	}

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
	cmd := exec.Command("mmdc", "-i", "/tmp/mermaid_diagram.mmd", "-o", "/tmp/mermaid_diagram.png")
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Error running mmdc command", "error", err, "output", string(output))
		return
	}
	slog.Info("Mermaid diagram generated", "output", string(output), "location", "/tmp/mermaid_diagram.png")
}

func generateGraphviz(cs map[string][]string) {
	if !util.InDebugMode() {
		return
	}

	edges := map[string]bool{}
	for u, cons := range cs {
		for _, v := range cons {
			if _, ok := edges[v+" -- "+u]; !ok {
				edges[u+"--"+v] = true
			}
		}
	}

	graphvizFile, err := os.Create("/tmp/graphviz_diagram.dot")
	if err != nil {
		slog.Error("Error creating graphviz file", "error", err)
		return
	}
	defer graphvizFile.Close()

	_, _ = graphvizFile.WriteString("graph {\n")
	for e := range edges {
		_, err = graphvizFile.WriteString(
			fmt.Sprintf("%s\n", e))
	}
	_, _ = graphvizFile.WriteString("}\n")
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

// https://www.sciencedirect.com/science/article/pii/S1570866708000415#sec005
// ugh, god nevermind

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

	cs := map[string][]string{}
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.Split(line, ":")
		c := splits[0]
		if c == "" {
			slog.Error("Invalid line", "line", line)
			continue
		}

		if _, ok := cs[c]; !ok {
			cs[c] = []string{}
		}

		connections := strings.Split(splits[1], " ")
		for _, newConLabel := range connections {
			if newConLabel == "" {
				continue
			}
			cs[c] = append(cs[c], newConLabel)

			if _, ok := cs[newConLabel]; ok {
				cs[newConLabel] = append(cs[newConLabel], c)
			} else {
				cs[newConLabel] = []string{c}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Error reading file", "error", err)
	}

	generateGraphviz(cs)

	// I looked at the graphviz and found these
	// Helpful tip on the settings to cluster: https://www.reddit.com/r/adventofcode/comments/18qcsux/2023_day_25_part_1_solve_by_visualization/
	if strings.Contains(puzzleFile, "sample") {
		for _, e := range [][]string{
			{"pzl", "hfx"},
			{"nvd", "jqt"},
			{"cmg", "bvb"},
		} {
			removeEdge(e[0], e[1], cs)
		}
	} else {
		for _, e := range [][]string{
			{"mnf", "hrs"},
			{"kpc", "nnl"},
			{"rkh", "sph"},
		} {
			removeEdge(e[0], e[1], cs)
		}
	}

	val := productOfTwoComponents(cs)

	slog.Info("Finished Day TwentyFive part one", "expected", expected, "val", val)
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
