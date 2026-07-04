package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"

	"github.com/RifqiIrawan/ai-rag-platform/services/document-service/internal/config"
)

type DocumentHandler struct {
	DB    *pgxpool.Pool
	Minio *minio.Client
	Cfg   *config.Config
}

// Upload stores the file in MinIO and records its metadata in Postgres.
// ownerID is currently hardcoded/placeholder until auth middleware is wired
// in a later sprint (see rag-service/auth-service integration).
func (h *DocumentHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	objectKey := fmt.Sprintf("%s-%s", uuid.NewString(), header.Filename)

	_, err = h.Minio.PutObject(ctx, h.Cfg.MinioBucket, objectKey, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store file: " + err.Error()})
		return
	}

	ownerID := c.GetHeader("X-User-Id")
	if ownerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Id header is required"})
		return
	}

	var id string
	err = h.DB.QueryRow(ctx,
		`INSERT INTO documents (owner_id, filename, object_key, content_type, size_bytes)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		ownerID, header.Filename, objectKey, header.Header.Get("Content-Type"), header.Size,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record metadata: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "filename": header.Filename, "object_key": objectKey})
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
