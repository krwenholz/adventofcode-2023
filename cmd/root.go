package cmd

import (
	"adventofcode/cmd/dayFive"
	"adventofcode/cmd/dayFour"
	"adventofcode/cmd/daySix"
	"adventofcode/cmd/dayThree"
	"adventofcode/cmd/dayTwo"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

var puzzleInput string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "adventofcode",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().String("puzzle-input", "", "the puzzle input")
	rootCmd.MarkFlagRequired("puzzle-input")

	rootCmd.PersistentFlags().Bool("part-two", false, "Whether to run part two of the day's challenge")

	// Logging configuration
	var logLevel slog.Level
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		logLevel = slog.LevelDebug
	default:
		logLevel = slog.LevelInfo
	}
	//h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	//slog.SetDefault(slog.New(h))
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      logLevel,
			TimeFormat: time.TimeOnly,
		}),
	))

	for _, c := range []*cobra.Command{
		dayTwo.Cmd,
		dayThree.Cmd,
		dayFour.Cmd,
		dayFive.Cmd,
		daySix.Cmd,
	} {
		rootCmd.AddCommand(c)
	}
}
