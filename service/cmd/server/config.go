package server

import "log/slog"

type Config struct {
	Connection string
	Logger     *slog.Logger
	Host       string
	Port       uint16
}
