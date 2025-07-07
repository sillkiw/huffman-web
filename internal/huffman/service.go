package huffman

import "log/slog"

type Service struct {
	Logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{Logger: logger}
}
