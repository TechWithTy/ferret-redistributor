package server

import (
    "database/sql"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/bitesinbyte/ferret/pkg/api/types"
    "github.com/bitesinbyte/ferret/pkg/engine/auth"
    "github.com/gin-gonic/gin"
)

func signupHandler(c *gin.Context) {
	var req types.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if sqlDB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db unavailable"})
		return
	}
	// Hash password (requires secure build tag)
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "password hashing not enabled"})
		return
	}
	// Create user and identity
	tx, err := sqlDB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()
	uid := newID()
	if _, err := tx.Exec(`INSERT INTO users(id, org_id, email, display_name, created_at, updated_at) VALUES($1,$2,$3,$4,NOW(),NOW())`, uid, nullOr(req.OrgID), req.Email, req.DisplayName); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	iid := newID()
	if _, err := tx.Exec(`INSERT INTO auth_identities(id, user_id, provider, identifier, secret_hash, is_primary, created_at, updated_at) VALUES($1,$2,'email',$3,$4,TRUE,NOW(),NOW())`, iid, uid, req.Email, hash); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user_id": uid})
}

func loginHandler(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if sqlDB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db unavailable"})
		return
	}
	var uid string
	var hash string
	err := sqlDB.QueryRow(`SELECT user_id, secret_hash FROM auth_identities WHERE provider='email' AND identifier=$1`, req.Email).Scan(&uid, &hash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !auth.CheckPassword(hash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	// If JWT configured, issue JWT; else issue opaque session
	if secret := os.Getenv("JWT_SECRET"); strings.TrimSpace(secret) != "" {
		ttl := 24 * time.Hour
		if v := os.Getenv("JWT_TTL"); v != "" {
			if d, err := time.ParseDuration(v); err == nil {
				ttl = d
			}
		}
		iss := os.Getenv("JWT_ISSUER")
		now := time.Now()
		tok, err := auth.SignJWT(auth.Claims{Issuer: iss, Subject: uid, IssuedAt: now.Unix(), Expires: now.Add(ttl).Unix()}, secret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, types.LoginResponse{AccessToken: tok, TokenType: "bearer", ExpiresIn: int(ttl.Seconds())})
		return
	}
	// Opaque session fallback
	tok, _ := auth.GenerateToken(32)
	th := auth.HashToken(tok)
	sid := newID()
	exp := time.Now().Add(24 * time.Hour)
	if _, err := sqlDB.Exec(`INSERT INTO auth_sessions(id, user_id, token_hash, created_at, last_seen_at, expires_at) VALUES($1,$2,$3,NOW(),NOW(),$4)`, sid, uid, th, exp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, types.LoginResponse{AccessToken: tok, TokenType: "bearer", ExpiresIn: int((24 * time.Hour).Seconds())})
}

func forgotHandler(c *gin.Context) {
	var req types.ForgotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if sqlDB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db unavailable"})
		return
	}
	// Lookup user id by email identity
	var uid string
	_ = sqlDB.QueryRow(`SELECT user_id FROM auth_identities WHERE provider='email' AND identifier=$1`, req.Email).Scan(&uid)
	// Always return 200 to avoid user discovery
	if uid != "" {
		tok, _ := auth.GenerateToken(32)
		exp := time.Now().Add(1 * time.Hour)
		rid := newID()
		_, _ = sqlDB.Exec(`INSERT INTO auth_password_resets(id, user_id, email, token, expires_at, created_at) VALUES($1,$2,$3,$4,$5,NOW())`, rid, uid, req.Email, tok, exp)
		// TODO: email delivery hook
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func resetHandler(c *gin.Context) {
	var req types.ResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if sqlDB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db unavailable"})
		return
	}
	var uid string
	var email string
	var expiresAt time.Time
	var usedAt sql.NullTime
	err := sqlDB.QueryRow(`SELECT user_id, email, expires_at, used_at FROM auth_password_resets WHERE token=$1`, req.Token).Scan(&uid, &email, &expiresAt, &usedAt)
	if err != nil || time.Now().After(expiresAt) || (usedAt.Valid && !usedAt.Time.IsZero()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}
	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "password hashing not enabled"})
		return
	}
	tx, err := sqlDB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`UPDATE auth_identities SET secret_hash=$1, updated_at=NOW() WHERE user_id=$2 AND provider='email' AND identifier=$3`, hash, uid, email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if _, err := tx.Exec(`UPDATE auth_password_resets SET used_at=NOW() WHERE token=$1`, req.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func newID() string { return time.Now().UTC().Format("20060102150405.000000000") }
func nullOr(s string) any {
	if s == "" {
		return nil
	}
	return s
}
