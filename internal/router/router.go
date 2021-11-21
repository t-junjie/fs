package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
var files = make(map[FileInfo]File)

// type File [][]int
type File interface{}
type FileInfo struct {
	ID             int
	FileName, Size string
}

type Message struct {
	Message string     `json:"message,omitempty"`
	Files   []FileInfo `json:"files,omitempty"`
}

func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	// Accept a multipart form of up to 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		msg := Message{Message: "File size exceeded"}
		writeResponse(w, http.StatusInternalServerError, msg)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		msg := Message{Message: "Could not get the uploaded file"}
		writeResponse(w, http.StatusInternalServerError, msg)
	}

	defer file.Close()
	fmt.Printf("filename:%s, size:%s\n", header.Filename, ByteCountSI(header.Size))

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		msg := Message{Message: "Fail to copy file data to buffer"}
		writeResponse(w, http.StatusInternalServerError, msg)
	}

	// TODO: Find a way to create binary data to test
	// var dataString [][]int
	// err = binary.Read(buf, binary.LittleEndian, dataString)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte("binary.Read failed"))
	// 	fmt.Println(err)
	// }

	fInfo := FileInfo{
		ID:       rand.Intn(100),
		FileName: header.Filename,
		Size:     ByteCountSI(header.Size),
	}

	files[fInfo] = buf
	msg := Message{
		Message: "success",
		Files:   []FileInfo{fInfo},
	}
	writeResponse(w, http.StatusOK, msg)

}

func getFiles(w http.ResponseWriter, r *http.Request) {
	var info []FileInfo
	for i := range files {
		info = append(info, i)
	}

	msg := Message{
		Message: "success",
		Files:   info}

	writeResponse(w, http.StatusOK, msg)
}

func getFileDetails(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "fileID")
	id, err := strconv.Atoi(fileID)
	if err != nil {
		// return 500
	}

	var f File
	for i, file := range files {
		if id == i.ID {
			f = file
		}
	}

	fmt.Println(f)

}

func removeFile(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "fileID")
	id, err := strconv.Atoi(fileID)
	if err != nil {
		// return 500
	}

	var info FileInfo
	for i := range files {
		if id == i.ID {
			info = i
			delete(files, i)
		}
	}

	msg := Message{
		Message: "success",
		Files:   []FileInfo{info},
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
