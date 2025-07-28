package app

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
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
	app.render(w, http.StatusOK, "decoding_page.html", nil)
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

// processFn takes the raw bytes and returns the transformed payload (or an error).
type processFn func([]byte) ([]byte, error)

func (app *App) uploadFileToEncode(w http.ResponseWriter, r *http.Request) {
	app.uploadAndServe(
		w, r,
		"encode",
		app.HuffSvc.Encode,
		"download_page_encoding.html",
		".bin",
	)
}

func (app *App) uploadFileToDecode(w http.ResponseWriter, r *http.Request) {
	app.uploadAndServe(
		w, r,
		"decode",
		app.HuffSvc.Decode,
		"download_page_decoding.html",
		".txt",
	)
}

// uploadAndServe runs the common logic for both Encode and Decode handlers.
func (app *App) uploadAndServe(
	w http.ResponseWriter,
	r *http.Request,
	action string, // “encode” or “decode” (for logging & template selection)
	fn processFn, // app.HuffSvc.Encode or app.HuffSvc.Decode
	tplName string, // e.g. "download_page_encoding.html"
	outExt string, // ".bin" or ".txt"
) {
	// Parse upload
	file, hdr, err := r.FormFile("upload")
	if err != nil {
		app.Logger.Error("uploadAndServe: FormFile failed",
			slog.String("action", action),
			slog.String("error", err.Error()))
		app.clientError(w, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read data
	data, err := io.ReadAll(file)
	if err != nil {
		app.Logger.Error("uploadAndServe: ReadAll failed",
			slog.String("action", action),
			slog.String("error", err.Error()))
		app.serverError(w, err)
		return
	}

	// Process (encode or decode)
	outData, err := fn(data)
	if err != nil {
		app.Logger.Error("uploadAndServe: processing failed",
			slog.String("action", action),
			slog.String("error", err.Error()))
		app.serverError(w, err)
		return
	}

	//  Make a safe, unique filename
	base := strings.TrimSuffix(filepath.Base(hdr.Filename), filepath.Ext(hdr.Filename))
	id := uuid.New().String()
	name := fmt.Sprintf("%s-%s%s", id, base, outExt)
	app.Logger.Info("uploadAndServe: created file",
		slog.String("action", action),
		slog.String("file", name),
	)

	// Write to disk
	tempPath := filepath.Join("uploads", name)
	if err := os.WriteFile(tempPath, outData, 0o644); err != nil {
		app.Logger.Error("uploadAndServe: WriteFile failed",
			slog.String("action", action),
			slog.String("error", err.Error()))
		app.serverError(w, err)
		return
	}

	// Render download page
	app.render(w, http.StatusOK, tplName, map[string]string{
		"FilePath": tempPath,
	})
}
