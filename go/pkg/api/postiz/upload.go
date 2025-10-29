package postiz

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

type UploadService struct{ c *Client }

// Upload multipart file
func (u *UploadService) Upload(ctx context.Context, filename string, content []byte) (*UploadResponse, error) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(fw, bytes.NewReader(content)); err != nil {
		return nil, err
	}
	_ = mw.Close()
	req, err := u.c.newRequest(ctx, http.MethodPost, "/upload", &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	var out UploadResponse
	if err := u.c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Upload from URL
func (u *UploadService) UploadFromURL(ctx context.Context, in UploadFromURLRequest) (*UploadResponse, error) {
	b, _ := json.Marshal(in)
	req, err := u.c.newRequest(ctx, http.MethodPost, "/upload-from-url", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var out UploadResponse
	if err := u.c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
