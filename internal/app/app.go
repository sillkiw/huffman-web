package app

import (
	"log/slog"
	"text/template"

	"github.com/sillkiw/huffman-web/internal/huffman"
	"github.com/sillkiw/huffman-web/internal/templates"
)

type App struct {
	Logger        *slog.Logger
	HuffSvc       *huffman.Service
	TemplateCache map[string]*template.Template
}

func NewApp(logger *slog.Logger) *App {
	svc := huffman.NewService(logger)
	templateCache, err := templates.NewTemplateCache("ui/html")
	if err != nil {
		logger.Error("Failed to load templates", slog.String("error", err.Error()))
	}
	return &App{Logger: logger, HuffSvc: svc, TemplateCache: templateCache}
}
