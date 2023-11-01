package main

import (
	// "context"
	// "log"
	// "net"
	"log"
	"log/slog"
	"os"

	// "github.com/adlerhurst/eventstore/cockroachdb"
	// "github.com/adlerhurst/eventstore/service/internal/api"
	// eventstorev1alpha "github.com/adlerhurst/eventstore/service/internal/api/eventstore/v1alpha"
	// "github.com/jackc/pgx/v5/pgxpool"
	"github.com/adlerhurst/eventstore/service/cmd/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	// "google.golang.org/grpc"
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file (default is $HOME/.eventstore/config.yaml)")
	rootCmd.AddCommand(server.Command)
	rootCmd.Version = version
}

func main() {
	err := rootCmd.Execute()
	log.Fatal("failed to execute command", err)
}

func initConfig() {
	viper.SetEnvPrefix("eventstore")

	setConfig()

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Info("no config found fall back to default")
	}
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
