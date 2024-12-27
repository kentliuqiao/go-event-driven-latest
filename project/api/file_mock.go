package api

import (
	"context"
	"sync"
	"tickets/entities"
)

type FileMock struct {
	mu    sync.Mutex
	Files map[string]string
}

func (f *FileMock) UploadFile(ctx context.Context, req entities.GenerateFileRequest) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.Files == nil {
		f.Files = make(map[string]string)
	}

	f.Files[req.FileID] = req.FileContent

	return nil
}

func (f *FileMock) DownloadFile(ctx context.Context, req entities.DownloadFileRequest) ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for fileID, content := range f.Files {
		if fileID == req.FileID {
			return []byte(content), nil
		}
	}

	return nil, nil
}
