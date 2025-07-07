package app

import "net/http"

func (app *App) Routes() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/upload", app.uploadFile)
	mux.HandleFunc("/download", app.downloadHandler)
	mux.HandleFunc("/encoding", app.encodingPage)

	return mux
}
