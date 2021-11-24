package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/t-junjie/fs/pkg/storage"
)

func TestUploadFile(t *testing.T) {}
func TestGetFiles(t *testing.T) {
	// happy path without files on server
	t.Run("returns all files on server", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/v1/files", nil)
		res := httptest.NewRecorder()

		getFiles(res, req)

		// check status code and response body
		if status := res.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %d, want %d", status, http.StatusOK)
		}

		expected := FilesResponse{[]storage.FileInfo{}}
		var actual FilesResponse
		json.NewDecoder(res.Body).Decode(&actual)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("handler returned unexpected body: got %v, want %v", actual, expected)
		}
	})
}

func TestGetFileDetails(t *testing.T) {}
func TestRemoveFile(t *testing.T)     {}
func TestComputeSum(t *testing.T)     {}
