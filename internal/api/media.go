package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// MediaService handles media upload and download operations on the Google Chat API.
type MediaService struct {
	client *Client
}

// NewMediaService creates a new MediaService bound to the given API client.
func NewMediaService(client *Client) *MediaService {
	return &MediaService{client: client}
}

// Upload uploads a file as an attachment to the specified parent space.
// POST /v1/{parent}/attachments:upload
func (s *MediaService) Upload(ctx context.Context, parent string, filePath string) (json.RawMessage, error) {
	parent = NormalizeName(parent, "spaces/")

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file %s: %w", filePath, err)
	}
	defer f.Close()

	// Detect the content type from the file extension, falling back to
	// application/octet-stream.
	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Build a multipart request body containing the file.
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("creating multipart form file: %w", err)
	}

	if _, err := io.Copy(part, f); err != nil {
		return nil, fmt.Errorf("copying file data: %w", err)
	}

	// Add the filename metadata field.
	if err := writer.WriteField("filename", filepath.Base(filePath)); err != nil {
		return nil, fmt.Errorf("writing filename field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("closing multipart writer: %w", err)
	}

	path := parent + "/attachments:upload"
	return s.client.Upload(ctx, path, nil, &buf, writer.FormDataContentType())
}

// Download downloads media content by resource name.
// GET /v1/media/{resourceName}?alt=media
// Returns the response body as a ReadCloser, the Content-Type header, and any error.
func (s *MediaService) Download(ctx context.Context, resourceName string) (io.ReadCloser, string, error) {
	path := "media/" + resourceName
	// The Download method on Client builds the full URL. We need to append
	// the alt=media query parameter. Since Client.Download does not accept
	// query params, we encode them directly into the path.
	fullPath := path + "?alt=media"

	// Use a direct HTTP GET to handle the alt=media query parameter.
	reqURL := s.client.BaseURL + "/" + fullPath

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("creating download request: %w", err)
	}

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("executing download request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		apiErr := parseAPIErrorFromBody(resp.StatusCode, body)
		if apiErr != nil {
			return nil, "", apiErr
		}
		return nil, "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	ct := resp.Header.Get("Content-Type")
	return resp.Body, ct, nil
}
