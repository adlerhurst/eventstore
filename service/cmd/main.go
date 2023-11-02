package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/adlerhurst/eventstore/service/cmd/client"
	"github.com/adlerhurst/eventstore/service/cmd/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configPath string
	logger     = slog.Default()

	rootCmd = &cobra.Command{
		Use:   "eventstore",
		Short: "eventstore is the only storage you need",
		Long:  `Eventstore is the best implementation of an eventstore`,
	}

	version = "dev"
)

func init() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})
	logger = slog.New(logHandler)
	slog.SetDefault(logger)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file (default is $HOME/.eventstore/config.yaml)")
	rootCmd.AddCommand(server.Command)
	rootCmd.AddCommand(client.Command)
	rootCmd.Version = version
}

func main() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func initConfig() {
	viper.SetEnvPrefix("eventstore")

	setConfig()

	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	var configNotFoundErr viper.ConfigFileNotFoundError
	if errors.As(err, &configNotFoundErr) {
		logger.Info("no config found fall back to default")
		return
	}
	logger.Error("failed to read config", "cause", err)
	os.Exit(1)
}

func setConfig() {
	if configPath != "" {
		viper.SetConfigFile(configPath)
		return
	}
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".cobra" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".eventstore")
}
