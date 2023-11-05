package client

import (
	"log"
	"log/slog"
	"net"
	"strconv"

	"github.com/adlerhurst/eventstore/service/internal/api/eventstore/v1alpha"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Logger *slog.Logger
	Host   string
	Port   uint16
}

var (
	config = Config{
		Logger: slog.Default(),
		Host:   "localhost",
		Port:   8080,
	}

	Command = &cobra.Command{
		Use:              "client",
		Short:            "makes an api call",
		PersistentPreRun: connect,
	}

	client eventstorev1alpha.EventStoreServiceClient
)

func init() {
	viper.SetDefault("config", config)

	Command.AddCommand(pushCmd)
}

func connect(cmd *cobra.Command, args []string) {
	conn, err := grpc.Dial(
		net.JoinHostPort(config.Host, strconv.Itoa(int(config.Port))),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("unable to connect to server %v", err)
	}

	client = eventstorev1alpha.NewEventStoreServiceClient(conn)
}
