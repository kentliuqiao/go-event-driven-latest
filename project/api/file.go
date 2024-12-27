package api

import (
	"context"
	"fmt"
	"net/http"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

type FilesAPIClient struct {
	client *clients.Clients
}

func NewFilesAPIClient(client *clients.Clients) *FilesAPIClient {
	if client == nil {
		panic("NewFileAPIClient: client is nil")
	}

	return &FilesAPIClient{client: client}
}

func (c FilesAPIClient) UploadFile(ctx context.Context, req entities.GenerateFileRequest) error {
	resp, err := c.client.Files.PutFilesFileIdContentWithTextBodyWithResponse(ctx, req.FileID, req.FileContent)
	if err != nil {
		return err
	}

	if resp.StatusCode() == http.StatusConflict {
		log.FromContext(ctx).Infof("file %s already exists", req.FileID)
		return nil
	}
	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("unexpected status code while uploading file %s: %d", req.FileID, resp.StatusCode())
	}

	return nil
}

func (c FilesAPIClient) DownloadFile(ctx context.Context, req entities.DownloadFileRequest) ([]byte, error) {
	resp, err := c.client.Files.GetFilesFileIdContentWithResponse(ctx, req.FileID)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code while downloading file %s: %d", req.FileID, resp.StatusCode())
	}

	return resp.Body, nil
}
