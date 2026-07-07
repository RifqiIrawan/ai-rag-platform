package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"

	"github.com/RifqiIrawan/ai-rag-platform/services/document-service/internal/clients"
	"github.com/RifqiIrawan/ai-rag-platform/services/document-service/internal/config"
)

type DocumentHandler struct {
	DB        *pgxpool.Pool
	Minio     *minio.Client
	OCR       *clients.OCRClient
	Embedding *clients.EmbeddingClient
	Cfg       *config.Config
}

// Upload stores the file in MinIO, records its metadata in Postgres, and
// kicks off async ingestion (OCR + embedding) so the document becomes
// searchable via rag-service without blocking the upload response.
func (h *DocumentHandler) Upload(c *gin.Context) {
	ownerID := c.GetHeader("X-User-Id")
	if ownerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Id header is required"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	objectKey := fmt.Sprintf("%s-%s", uuid.NewString(), header.Filename)
	contentType := header.Header.Get("Content-Type")

	_, err = h.Minio.PutObject(ctx, h.Cfg.MinioBucket, objectKey, bytes.NewReader(content), int64(len(content)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store file: " + err.Error()})
		return
	}

	var id string
	err = h.DB.QueryRow(ctx,
		`INSERT INTO documents (owner_id, filename, object_key, content_type, size_bytes)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		ownerID, header.Filename, objectKey, contentType, header.Size,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record metadata: " + err.Error()})
		return
	}

	go h.ingest(id, header.Filename, contentType, content)

	c.JSON(http.StatusCreated, gin.H{"id": id, "filename": header.Filename, "object_key": objectKey})
}

// ingest extracts text from the uploaded file (OCR for images, raw bytes for
// plain text) and indexes it into Qdrant via embedding-service, updating the
// document's status as it progresses. Runs detached from the request.
func (h *DocumentHandler) ingest(documentID, filename, contentType string, content []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()

	h.setStatus(ctx, documentID, "processing")

	var text string
	if strings.HasPrefix(contentType, "text/") {
		text = string(content)
	} else {
		extracted, err := h.OCR.ExtractText(ctx, filename, content)
		if err != nil {
			log.Printf("document %s: ocr failed: %v", documentID, err)
			h.setStatus(ctx, documentID, "failed")
			return
		}
		text = extracted
	}

	if strings.TrimSpace(text) == "" {
		log.Printf("document %s: no text extracted, skipping indexing", documentID)
		h.setStatus(ctx, documentID, "indexed")
		return
	}

	if err := h.Embedding.Index(ctx, h.Cfg.QdrantCollection, text); err != nil {
		log.Printf("document %s: embedding failed: %v", documentID, err)
		h.setStatus(ctx, documentID, "failed")
		return
	}

	h.setStatus(ctx, documentID, "indexed")
}

func (h *DocumentHandler) setStatus(ctx context.Context, documentID, status string) {
	if _, err := h.DB.Exec(ctx, `UPDATE documents SET status = $1, updated_at = now() WHERE id = $2`, status, documentID); err != nil {
		log.Printf("document %s: failed to set status %q: %v", documentID, status, err)
	}
}

func (h *DocumentHandler) List(c *gin.Context) {
	ownerID := c.GetHeader("X-User-Id")
	if ownerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Id header is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	rows, err := h.DB.Query(ctx,
		`SELECT id, filename, content_type, size_bytes, status, created_at
		 FROM documents WHERE owner_id = $1 ORDER BY created_at DESC`, ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type docSummary struct {
		ID          string    `json:"id"`
		Filename    string    `json:"filename"`
		ContentType string    `json:"content_type"`
		SizeBytes   int64     `json:"size_bytes"`
		Status      string    `json:"status"`
		CreatedAt   time.Time `json:"created_at"`
	}

	docs := []docSummary{}
	for rows.Next() {
		var d docSummary
		if err := rows.Scan(&d.ID, &d.Filename, &d.ContentType, &d.SizeBytes, &d.Status, &d.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		docs = append(docs, d)
	}

	c.JSON(http.StatusOK, gin.H{"documents": docs})
}
