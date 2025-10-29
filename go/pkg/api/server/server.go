package server

import (
	"database/sql"
	"embed"
	"net/http"
	"time"

	"github.com/bitesinbyte/ferret/pkg/api/types"
	appdb "github.com/bitesinbyte/ferret/pkg/db"
	"github.com/bitesinbyte/ferret/pkg/telemetry"
	"github.com/gin-gonic/gin"
)

var sqlDB = (*sql.DB)(nil)

//go:embed openapi.json
var openapiFS embed.FS

// New returns a configured Gin engine with routes registered.
func New() *gin.Engine {
	mode := gin.Mode()
	if mode == "" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestID())
	r.Use(requestLogger())
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/openapi.json", func(c *gin.Context) { c.FileFromFS("openapi.json", http.FS(openapiFS)) })
	// Open DB once
	if db, err := appdb.OpenFromEnv(); err == nil {
		sqlDB = db
	}

	v1 := r.Group("/v1")
	{
		v1.POST("/auth/signup", signupHandler)
		v1.POST("/auth/login", loginHandler)
		v1.POST("/auth/forgot", forgotHandler)
		v1.POST("/auth/reset", resetHandler)
		// Authenticated endpoints
		v1.Use(authMiddleware())
		v1.GET("/users/:id", userHandler)
		v1.GET("/profile", getProfile)
		v1.PUT("/profile", updateProfile)
		v1.GET("/icp", getICP)
		v1.PUT("/icp", updateICP)
	}
	return r
}

func userHandler(c *gin.Context) {
	_ = telemetry.InitFromEnv(c.Request.Context())
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	// Demo payload
	c.JSON(http.StatusOK, types.UserDTO{ID: id, Email: "user@example.com", DisplayName: "Demo User"})
}

// requestID sets X-Request-ID if missing.
func requestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Writer.Header().Get("X-Request-ID") == "" {
			c.Writer.Header().Set("X-Request-ID", newRID())
		}
		c.Next()
	}
}

// requestLogger logs minimal request info.
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start)
		_ = dur // hook to your logger if desired
	}
}

func newRID() string { return time.Now().UTC().Format("20060102T150405.000000000Z07:00") }
