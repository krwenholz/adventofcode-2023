package cmd

import (
	"log/slog"
	"os"

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

	rootCmd.PersistentFlags().StringVar(&puzzleInput, "puzzle-input", "", "the puzzle input")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	var logLevel slog.Level
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		logLevel = slog.LevelDebug
	default:
		logLevel = slog.LevelInfo
	}
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(h))
}
