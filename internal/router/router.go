// package router provides methods to upload, remove and retrieve information of
// files on a server.
package router

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/t-junjie/fs/pkg/storage"
)

// Start initiates the routing for server on http://localhost:8080
func Start() error {
	r := chi.NewRouter()

	r.Route("v1/files", func(r chi.Router) {
		// Provider
		r.Post("/", uploadFile)
		r.Get("/", getFiles)

		// Analyst
		r.Get("/{fileID}", getFileDetails)
		r.Delete("/{fileID}", removeFile)

		r.Post("/sum", computeSum)
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return err
	}

	return nil
}

// // Consider implementing dependency injection for testing HTTP handlers
// type Store struct {
//     data DataInterface
// }

// type DataInterface interface{
//   // presumably some method signature
// }

// func (s *Store) getFiles(w http.ResponseWriter, r *http.Request) {...}

// save file data on server until database is implemented
var files = make(map[uuid.UUID][][]int)
var filesInfo = make(map[uuid.UUID]storage.FileInfo)

type FilesResponse struct {
	Files []storage.FileInfo `json:"files"`
}

type FileResponse struct {
	Files storage.FileInfo `json:"files"`
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		writeResponseHeader(w, http.StatusInternalServerError)
		e := APIErrorResponse{ErrorMsg: "unable to process file"}
		json.NewEncoder(w).Encode(e)
		return
	}

	defer file.Close()

	// [filename ext]
	ext := strings.Split(header.Filename, ".")[1]
	if ext != "csv" {
		writeResponseHeader(w, http.StatusBadRequest)
		e := APIErrorResponse{ErrorMsg: "unable to process file"}
		json.NewEncoder(w).Encode(e)
		return
	}

	tbl, err := storage.ConvertToTable(file)
	if err != nil {
		writeResponseHeader(w, http.StatusInternalServerError)
		e := APIErrorResponse{ErrorMsg: "unable to process file"}
		json.NewEncoder(w).Encode(e)
		return
	}

	rowCount, colCount := storage.Size(tbl)

	id := uuid.New()
	info := storage.FileInfo{
		Id:       &id,
		Name:     header.Filename,
		Size:     ByteCountSI(header.Size),
		RowCount: rowCount,
		ColCount: colCount,
	}

	files[id] = tbl
	filesInfo[id] = info

	writeResponseHeader(w, http.StatusOK)
	json.NewEncoder(w).Encode(FileResponse{info})
}

func getFiles(w http.ResponseWriter, r *http.Request) {
	info := make([]storage.FileInfo, 0)
	for _, fi := range filesInfo {
		info = append(info, fi)
	}

	writeResponseHeader(w, http.StatusOK)
	if len(info) == 0 {
		json.NewEncoder(w).Encode(FilesResponse{info})
		return
	}

	json.NewEncoder(w).Encode(FilesResponse{info})
}

func getFileDetails(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "fileID")
	id, err := uuid.Parse(fileID)

	if err != nil {
		writeResponseHeader(w, http.StatusBadRequest)
		e := APIErrorResponse{ErrorMsg: "invalid uuid"}
		json.NewEncoder(w).Encode(e)
		return
	}

	info, ok := filesInfo[id]
	if !ok {
		writeResponseHeader(w, http.StatusNotFound)
		e := APIErrorResponse{ErrorMsg: "file not found"}
		json.NewEncoder(w).Encode(e)
		return
	}

	writeResponseHeader(w, http.StatusOK)
	json.NewEncoder(w).Encode(FileResponse{info})
}

func removeFile(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "fileID")
	id, err := uuid.Parse(fileID)

	if err != nil {
		writeResponseHeader(w, http.StatusBadRequest)
		e := APIErrorResponse{ErrorMsg: "invalid uuid"}
		json.NewEncoder(w).Encode(e)
		return
	}

	info, ok := filesInfo[id]
	if !ok {
		writeResponseHeader(w, http.StatusNotFound)
		e := APIErrorResponse{ErrorMsg: "file not found"}
		json.NewEncoder(w).Encode(e)
		return
	}

	delete(files, id)
	delete(filesInfo, id)

	writeResponseHeader(w, http.StatusOK)
	json.NewEncoder(w).Encode(FileResponse{info})
}

func computeSum(w http.ResponseWriter, r *http.Request) {
	// reject content if not json
	cTypeHead := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(cTypeHead)
	if err != nil || mediatype != "application/json" {
		writeResponseHeader(w, http.StatusBadRequest)
		r := struct {
			Empty *string `json:"empty,omitempty"`
		}{}
		json.NewEncoder(w).Encode(r)
		return
	}

	var cells []storage.Cell

	err = json.NewDecoder(r.Body).Decode(&cells)
	if err != nil {
		writeResponseHeader(w, http.StatusBadRequest)
		r := struct {
			Empty *string `json:"empty,omitempty"`
		}{}
		json.NewEncoder(w).Encode(r)
		return
	}

	// create a map to hold cell values sorted by the map key
	query := make(map[uuid.UUID][]storage.Cell)
	for _, cell := range cells {
		query[cell.Id] = append(query[cell.Id], cell)
	}

	var results []int
	for id, cells := range query {
		tbl := files[id]
		r := storage.Items(tbl, cells)
		results = append(results, r...)
	}

	sum := 0
	for _, s := range results {
		sum += s
	}

	writeResponseHeader(w, http.StatusOK)
	s := struct {
		Sum int `json:"sum"`
	}{Sum: sum}
	json.NewEncoder(w).Encode(s)
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

func writeResponseHeader(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
}
