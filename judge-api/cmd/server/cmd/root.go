// cmd/root.go

package cmd

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "your_program",
	Short: "A brief description of your program",
	Long:  `A longer description of your program and what it does`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	if err := Execute(); err != nil {
		os.Exit(1)
	}
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add flags for configuration
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error, fatal)")

	rootCmd.AddCommand(serveCmd)

	// Set up the context
	cobra.OnInitialize(func() {
		// Set up the logger
		logger := logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{})
		logger.SetLevel(logrus.InfoLevel)
		logger.AddHook(&writer.Hook{
			Writer:    os.Stdout,
			LogLevels: []logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel},
		})

		// Create a new context with the logger
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", logger)

		// Set the context on the root command
		rootCmd.SetContext(ctx)
	})
}
