package app

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func (app *App) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.Logger.Info("page not found",
			slog.String("path", r.URL.Path),
		)
		app.notFound(w)
		return
	}
	app.Logger.Info("request started",
		slog.String("ip", r.RemoteAddr),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("userAgent", r.UserAgent()),
	)
	app.render(w, http.StatusOK, "home_page.html", nil)
}

func (app *App) encodingPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	app.render(w, http.StatusOK, "encoding_page.html", nil)
}

func (app *App) decodingPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	app.render(w, http.StatusOK, "encoding_page.html", nil)
}

func (app *App) uploadFile(w http.ResponseWriter, r *http.Request) {
	f, fileHandler, err := r.FormFile("upload")
	if err != nil {
		app.Logger.Error("Failed to receive uploaded file", slog.String("error", err.Error()))
		app.clientError(w, http.StatusBadRequest)
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		app.Logger.Error("Failed to read uploaded file", slog.String("error", err.Error()))
		app.serverError(w, err)
		return
	}

	encoded := app.HuffSvc.Encode(data)

	tempPath := filepath.Join("uploads", fileHandler.Filename+".bin")
	err = os.WriteFile(tempPath, encoded, 0644)
	if err != nil {
		app.Logger.Error("Failed to write to temp file", slog.String("error", err.Error()))
		app.serverError(w, err)
		return
	}
	app.Logger.Info("Create temp file with encoded data")

	app.render(w, http.StatusOK, "download_page.html", map[string]string{
		"FilePath": tempPath,
	})
}

func (app *App) downloadHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("file")
	if path == "" || !filepath.HasPrefix(path, "uploads") {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=result.huff")
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, path)
	app.Logger.Info("Send encoded file to user")
}
