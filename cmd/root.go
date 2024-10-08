package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
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
	start := time.Now()
	defer func() {
		slog.Info("finished", "time", fmt.Sprintf("%.2f", float64(time.Now().Sub(start).Microseconds())/1000.0))
	}()
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
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "trace":
		logLevel = slog.LevelDebug - 1
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
}
