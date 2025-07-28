package app

import "net/http"

func (app *App) Routes() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/upload-encoding", app.uploadFileToEncode)
	mux.HandleFunc("/upload-decoding", app.uploadFileToDecode)
	mux.HandleFunc("/download", app.downloadHandler)
	mux.HandleFunc("/encoding", app.encodingPage)
	mux.HandleFunc("/decoding", app.decodingPage)

	return mux
}
