package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Start() error {
	r := chi.NewRouter()

	r.Route("/files", func(r chi.Router) {
		// Provider
		r.Post("/", uploadFile)
		r.Get("/", getFiles)

		// Analyst
		r.Get("/{fileID}", getFileDetails)
		r.Delete("/{fileID}", removeFile)
		r.Post("/{fileID}", computeSum)
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return err
	}

	return nil
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("uploadFile endpoint reached"))
}

func getFiles(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getFilesendpoint reached"))
}

func getFileDetails(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getFileDetails endpoint reached"))
}

func removeFile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("removeFile endpoint reached"))
}

func computeSum(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("computeSum endpoint reached"))
}
