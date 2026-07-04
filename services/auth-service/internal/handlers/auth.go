package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/RifqiIrawan/ai-rag-platform/services/auth-service/internal/config"
)

type AuthHandler struct {
	DB  *pgxpool.Pool
	Cfg *config.Config
}

type credentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var creds credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var id string
	err = h.DB.QueryRow(ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`,
		creds.Email, string(hash),
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "could not create user (email may already be registered)"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "email": creds.Email})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var creds credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var id, passwordHash string
	err := h.DB.QueryRow(ctx,
		`SELECT id, password_hash FROM users WHERE email = $1`, creds.Email,
	).Scan(&id, &passwordHash)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lookup failed"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(creds.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	expiryHours, _ := time.ParseDuration(h.Cfg.JWTExpiryHours + "h")
	if expiryHours == 0 {
		expiryHours = 24 * time.Hour
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(expiryHours).Unix(),
	})

	signed, err := token.SignedString([]byte(h.Cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": signed})
}
