package app

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func (app *App) serverError(w http.ResponseWriter, err error) {
	stack := debug.Stack()
	app.Logger.Error(
		"internal server error",
		slog.String("error", err.Error()),
		slog.String("stack", string(stack)),
	)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *App) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *App) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
