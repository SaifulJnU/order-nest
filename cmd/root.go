package cmd

import (
	"os"

	"github.com/order-nest/config"
	appLogger "github.com/order-nest/pkg/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd is the base command for order-nest.
// It is the entry point for the CLI application.
// All subcommands are attached to this root command.
var rootCmd = &cobra.Command{
	Use:   "order-nest",
	Short: "Order Nest: A order management system",
	Long: `Order Nest is a complete order management system for handling
order processing, tracking, and management in production environments.
It supports configuration via CLI flags, environment variables, or YAML files.`,
}

// Execute runs the root command and triggers command parsing and execution.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// It sets up Cobra CLI initialization logic.
func init() {
	// Initialize config handling before running commands
	cobra.OnInitialize(initConfig)

	// Define persistent flags (available to all subcommands)
	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file (default is $HOME/.order-nest.yaml)",
	)

	// Define local flags (specific to root command)
	rootCmd.Flags().BoolP(
		"toggle",
		"t",
		false,
		"Help message for toggle flag",
	)
}

// Priority order: CLI flag > environment variables > default config file.
func initConfig() {
	if cfgFile != "" {
		// Use config file specified via CLI flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory to search for default config file
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Default config search path: $HOME/.order-nest.yaml
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".order-nest")
	}

	// Automatically override config with matching environment variables
	viper.AutomaticEnv()

	// Initialize logger early with service name
	appLogger.Init("order-nest")

	// Read the config file if present
	if err := viper.ReadInConfig(); err == nil {
		appLogger.L().WithField("config_file", viper.ConfigFileUsed()).Info("configuration loaded")

		// Load config into application
		if loadErr := config.Load(); loadErr != nil {
			appLogger.L().WithError(loadErr).Error("failed to parse configuration")
		}
	} else {
		// It is acceptable to run with env-only configuration
		appLogger.L().WithError(err).Warn("no config file found; using environment variables/defaults")
		if loadErr := config.Load(); loadErr != nil {
			appLogger.L().WithError(loadErr).Error("failed to load configuration from environment/defaults")
		}
	}
}
