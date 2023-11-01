package server

import (
	"context"
	"log"
	"log/slog"
	"net"
	"strconv"

	"github.com/adlerhurst/eventstore/service/internal/api"
	"github.com/adlerhurst/eventstore/service/internal/api/eventstore/v1alpha"
	"github.com/adlerhurst/eventstore/v2/cockroachdb"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	config = Config{
		Connection: "postgresql://root@localhost:26257/eventstore?sslmode=disable",
		Logger:     slog.Default(),
		Host:       "localhost",
		Port:       8080,
	}

	Command = &cobra.Command{
		Use:   "server",
		Short: "starts the server",
		Run:   run,
	}
)

func init() {
	viper.SetDefault("config", config)
}

func run(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	poolConfig, err := pgxpool.ParseConfig(config.Connection)
	if err != nil {
		log.Fatalf("unable to parse conn string: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("unable to create database pool: %v", err)
	}

	es := cockroachdb.New(&cockroachdb.Config{
		Pool: pool,
	})

	server := api.NewServer(ctx, es)

	listener, err := net.Listen("tcp", net.JoinHostPort(config.Host, strconv.Itoa(int(config.Port))))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	config.Logger.Info("listen on", "addr", listener.Addr().String())

	grpcServer := grpc.NewServer()
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	eventstorev1alpha.RegisterEventStoreServiceServer(grpcServer, server)
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
