package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type EmbeddingClient struct {
	baseURL string
	http    *http.Client
}

func NewEmbeddingClient(baseURL string) *EmbeddingClient {
	return &EmbeddingClient{
		baseURL: baseURL,
		// embedding-service lazy-loads the bge-m3 model on its first ever
		// request, which can take several minutes (multi-GB download); a
		// short timeout here would spuriously fail every document uploaded
		// before that warm-up completes.
		http: &http.Client{Timeout: 5 * time.Minute},
	}
}

// Index embeds text and upserts it into the given Qdrant collection via embedding-service.
func (e *EmbeddingClient) Index(ctx context.Context, collection, text string) error {
	body, err := json.Marshal(map[string]string{"collection": collection, "text": text})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.baseURL+"/api/v1/embeddings/generate", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("embedding-service returned status %d", resp.StatusCode)
	}
	return nil
}
