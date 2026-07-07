package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type OCRClient struct {
	baseURL string
	http    *http.Client
}

func NewOCRClient(baseURL string) *OCRClient {
	return &OCRClient{
		baseURL: baseURL,
		// ocr-service lazy-loads PaddleOCR on its first ever request, which
		// can take significantly longer than a typical request timeout.
		http: &http.Client{Timeout: 5 * time.Minute},
	}
}

// Extract sends the file to ocr-service and returns the recognized lines of text.
func (o *OCRClient) Extract(ctx context.Context, filename string, content []byte) ([]string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(content); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/api/v1/ocr/extract", &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := o.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ocr-service returned status %d", resp.StatusCode)
	}

	var parsed struct {
		Lines []string `json:"lines"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}
	return parsed.Lines, nil
}

// ExtractText is a convenience wrapper that joins the recognized lines.
func (o *OCRClient) ExtractText(ctx context.Context, filename string, content []byte) (string, error) {
	lines, err := o.Extract(ctx, filename, content)
	if err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
}
