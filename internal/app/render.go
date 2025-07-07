package app

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (app *App) render(w http.ResponseWriter, status int, name string, data any) {
	tmpl, ok := app.TemplateCache[name]
	if !ok {
		app.Logger.Error("Template not found", slog.String("name", name))
		app.serverError(w, fmt.Errorf("template %s not found", name))
		return
	}

	w.WriteHeader(status)

	err := tmpl.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}
