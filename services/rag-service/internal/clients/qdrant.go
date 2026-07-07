package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type QdrantClient struct {
	baseURL string
	http    *http.Client
}

func NewQdrantClient(baseURL string) *QdrantClient {
	return &QdrantClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

// Healthy checks Qdrant reachability via its collections endpoint.
func (q *QdrantClient) Healthy(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, q.baseURL+"/collections", nil)
	if err != nil {
		return err
	}
	resp, err := q.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("qdrant returned status %d", resp.StatusCode)
	}
	return nil
}

// CollectionExists reports whether the named collection exists.
func (q *QdrantClient) CollectionExists(ctx context.Context, collection string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, q.baseURL+"/collections/"+collection, nil)
	if err != nil {
		return false, err
	}
	resp, err := q.http.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("qdrant returned status %d", resp.StatusCode)
	}
}

type SearchResult struct {
	ID      any             `json:"id"`
	Score   float64         `json:"score"`
	Payload json.RawMessage `json:"payload"`
}

// Search performs a vector similarity search against a collection.
// Real retrieval logic (embedding the query, ranking, filtering) is a
// follow-up sprint item; this wires the transport.
func (q *QdrantClient) Search(ctx context.Context, collection string, vector []float32, limit int) ([]SearchResult, error) {
	body, err := json.Marshal(map[string]any{
		"vector":       vector,
		"limit":        limit,
		"with_payload": true,
	})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/collections/%s/points/search", q.baseURL, collection)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := q.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("qdrant search returned status %d", resp.StatusCode)
	}

	var parsed struct {
		Result []SearchResult `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}
	return parsed.Result, nil
}
