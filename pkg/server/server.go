// package server provides methods to upload, remove and retrieve information of
// files on a server.
package server

import "context"

type FileService interface {
	Uploader
	Remover
	Getter
	Computer
}

type Uploader interface {
	Upload(ctx context.Context, binFileName string) error
}

type Remover interface {
	Remove(ctx context.Context, binFileName string) error
}

type Getter interface {
	GetFileNames(ctx context.Context) ([]string, error)
	GetFileDetails(ctx context.Context, binFileName string) (*fileDetails, error)
}

type Computer interface {
	Sum(ctx context.Context, q []fileQuery) (int, error)
}

type fileDetails struct {
	Name                  string
	RowCount, ColumnCount int
}

type fileQuery struct {
	FileIdentifier string
	Row, Column    int
}

func Upload(ctx context.Context, binFileName string) error {
	return nil
}

func Remove(ctx context.Context, binFileName string) error {
	return nil
}

func GetFileDetails(ctx context.Context, binFileName string) (*fileDetails, error) {
	return nil, nil
}

func GetFileNames(ctx context.Context, binFileName string) ([]string, error) {
	return nil, nil
}

func Sum(ctx context.Context, binFileName string) (int, error) {
	return 0, nil
}
