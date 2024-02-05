package eventstorev1alpha

import (
	"log/slog"
)

type Config struct {
	Logger *slog.Logger
}

var (
	DefaultConfig = Config{
		Logger: slog.Default(),
	}
)
