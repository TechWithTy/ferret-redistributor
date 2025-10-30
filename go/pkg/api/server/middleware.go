package server

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/bitesinbyte/ferret/pkg/engine/auth"
	"github.com/gin-gonic/gin"
	"os"
)

const ctxUserID = "user_id"

// authMiddleware checks Authorization: Bearer <token> against auth_sessions table.
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tok := strings.TrimSpace(h[len("Bearer "):])
		// Prefer JWT if configured
		if secret := getenv("JWT_SECRET", ""); secret != "" {
			if claims, err := auth.VerifyJWT(tok, secret); err == nil {
				if claims.Subject != "" {
					c.Set(ctxUserID, claims.Subject)
					c.Next()
					return
				}
			}
			// fall through to session check on failure
		}
		// Session fallback (opaque token)
		hash := auth.HashToken(tok)
		if sqlDB == nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "db unavailable"})
			return
		}
		var userID string
		var expiresAt time.Time
		var revoked sql.NullTime
		err := sqlDB.QueryRow(`SELECT user_id, expires_at, revoked_at FROM auth_sessions WHERE token_hash=$1`, hash).Scan(&userID, &expiresAt, &revoked)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		if !revoked.Time.IsZero() || time.Now().After(expiresAt) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "session expired or revoked"})
			return
		}
		c.Set(ctxUserID, userID)
		c.Next()
	}
}

func getenv(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}
