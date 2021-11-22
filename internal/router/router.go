// package router provides methods to upload, remove and retrieve information of
// files on a server.
package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/t-junjie/fs/pkg/storage"
)

// Start initiates the routing for server on http://localhost:8080
func Start() error {
	r := chi.NewRouter()

	r.Route("/files", func(r chi.Router) {
		// Provider
		r.Post("/", handleFileUpload)
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

// save file data on server until database is implemented
var files = make(map[storage.FileInfo][][]int)

type Message struct {
	Files []storage.FileInfo `json:"files"`
}

func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		msg := Message{}
		writeResponse(w, http.StatusBadRequest, msg)
	}

	defer file.Close()

	// [filename ext]
	ext := strings.Split(header.Filename, ".")[1]
	if ext != "csv" {
		msg := Message{}
		writeResponse(w, http.StatusBadRequest, msg)
	}

	tbl, err := storage.ConvertToTable(file)
	if err != nil {
		msg := Message{}
		writeResponse(w, http.StatusInternalServerError, msg)
	}

	rowCount, colCount := storage.Size(tbl)

	info := storage.FileInfo{
		Id:       uuid.New(),
		Name:     header.Filename,
		Size:     ByteCountSI(header.Size),
		RowCount: rowCount,
		ColCount: colCount,
	}

	files[info] = tbl

	msg := Message{Files: []storage.FileInfo{info}}
	writeResponse(w, http.StatusOK, msg)

}

func getFiles(w http.ResponseWriter, r *http.Request) {
	info := make([]storage.FileInfo, 0)
	for i := range files {
		info = append(info, i)
	}

	msg := Message{Files: info}

	writeResponse(w, http.StatusOK, msg)
}

func getFileDetails(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "fileID")

	info := make([]storage.FileInfo, 0)

	id, err := uuid.Parse(fileID)
	if err == nil {
		for i := range files {
			if id == i.Id {
				info = append(info, i)
				break
			}
		}
	}

	msg := Message{Files: info}
	writeResponse(w, http.StatusOK, msg)

}

func removeFile(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "fileID")
	id := uuid.MustParse(fileID)

	var info storage.FileInfo
	for i := range files {
		if id == i.Id {
			info = i
			delete(files, i)
		}
	}

	msg := Message{
		Files: []storage.FileInfo{info},
	}
	writeResponse(w, http.StatusOK, msg)
}

func computeSum(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("computeSum endpoint reached"))
}

// ByteCountSI returns a human-readable string representation of the file size
func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func writeResponse(w http.ResponseWriter, statusCode int, message Message) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}
