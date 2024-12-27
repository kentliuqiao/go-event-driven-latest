package entities

type GenerateFileRequest struct {
	FileID      string `json:"file_id"`
	FileContent string `json:"file_content"`
}

type DownloadFileRequest struct {
	FileID string `json:"file_id"`
}

type FileResponse struct {
}
