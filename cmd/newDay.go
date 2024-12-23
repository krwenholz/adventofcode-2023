package cmd

import (
	"html/template"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var NewDayCmd = &cobra.Command{
	Use:   "new-day",
	Short: "Create a new day command",
	Run: func(cmd *cobra.Command, args []string) {
		day, _ := cmd.Flags().GetString("day")
		day = strings.ToUpper(day[:1]) + day[1:]

		err := os.Mkdir("cmd/day"+day, 0755)
		if err != nil && !os.IsExist(err) {
			panic(err)
		}
		f, err := os.Create("cmd/day" + day + "/main.go")
		if err != nil {
			panic(err)
		}

		err = cmdTemplate.Execute(f, day)
		if err != nil {
			panic(err)
		}

		dayDirs, err := os.ReadDir("cmd")
		if err != nil {
			panic(err)
		}
		days := []string{}
		for _, d := range dayDirs {
			if d.IsDir() {
				if d.Name() != "dayTemplate" && strings.HasPrefix(d.Name(), "day") {
					days = append(days, strings.TrimPrefix(d.Name(), "day"))
				}
			}
		}
		f, err = os.Create("cmd/dayCommands.go")
		err = registryTemplate.Execute(f, days)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(NewDayCmd)
	NewDayCmd.Flags().String("day", "", "The new day")
}

var cmdTemplate = template.Must(template.New("cmdTemplate").Parse(
	`package day{{.}}

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func partOne(puzzleFile string) {
	slog.Info("Day {{.}} part one", "puzzle file", puzzleFile)
}

func partTwo(puzzleFile string) {
	slog.Info("Day {{.}} part two", "puzzle file", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "day{{.}}",
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
`,
))

var registryTemplate = template.Must(template.New("registryTemplate").Parse(
	`// Code generated by new-day DO NOT EDIT
package cmd

import (
	{{range .}}"adventofcode/cmd/day{{.}}"
	{{end}}

	"github.com/spf13/cobra"
)

func init() {
	for _, c := range []*cobra.Command{
		{{range .}}day{{.}}.Cmd,
		{{end}}
	} {
		rootCmd.AddCommand(c)
	}
}
	`,
))
