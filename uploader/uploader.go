package uploader

import (
	"context"
	"fmt"
	"github.com/duffpl/google-photos-api-client/internal"
	"github.com/gabriel-vasile/mimetype"
	"net/http"
	"os"
	"path"
)

type MediaUploader interface {
	UploadFile(filePath string, ctx context.Context) (string, error)
}

type HttpMediaUploader struct {
	client *internal.HttpClient
}

// Uploads file specified by path. Returns upload token
func (h HttpMediaUploader) UploadFile(filePath string, ctx context.Context) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot open file: %w", err)
	}
	mimeType, err := mimetype.DetectFile(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot detect mime type: %w", err)
	}
	token := ""
	err = h.client.PostFile("v1/uploads", nil, f, &token, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/octet-streams")
		req.Header.Set("X-Goog-Upload-File-Name", path.Base(filePath))
		req.Header.Set("X-Goog-Upload-Content-Type", mimeType.String())
		req.Header.Set("X-Goog-Upload-Protocol", "raw")
	}, ctx)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	return token, nil
}

func NewHttpMediaUploader(authenticatedClient *http.Client) HttpMediaUploader {
	return HttpMediaUploader{
		client: internal.NewHttpClient(authenticatedClient),
	}
}
