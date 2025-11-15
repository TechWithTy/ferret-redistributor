package postiz

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
)

// UploadFile uploads binary content using multipart/form-data.
func (c *Client) UploadFile(ctx context.Context, filename string, r io.Reader) (FileAsset, error) {
	if filename == "" {
		return FileAsset{}, ErrMissingFilename
	}
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return FileAsset{}, err
	}
	if _, err := io.Copy(part, r); err != nil {
		return FileAsset{}, err
	}
	if err := writer.Close(); err != nil {
		return FileAsset{}, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/upload", buf)
	if err != nil {
		return FileAsset{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var res FileAsset
	if err := c.do(req, &res); err != nil {
		return FileAsset{}, err
	}
	return res, nil
}

// UploadFromURL registers a file that already exists remotely.
func (c *Client) UploadFromURL(ctx context.Context, url string) (FileAsset, error) {
	if url == "" {
		return FileAsset{}, ErrMissingURL
	}
	payload := UploadFromURLRequest{URL: url}
	var res FileAsset
	if err := c.doJSON(ctx, http.MethodPost, "/upload-from-url", payload, &res); err != nil {
		return FileAsset{}, err
	}
	return res, nil
}
