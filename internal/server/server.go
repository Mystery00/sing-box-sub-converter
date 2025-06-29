package server

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"sing-box-sub-converter/internal/config"
	"sing-box-sub-converter/internal/converter"
	"sing-box-sub-converter/internal/template"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
}

// globalExceptionHandler is a middleware that recovers from any panics and returns a 500 error
func globalExceptionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				slog.Error("Panic recovered", "error", err)

				// Return a 500 error with the error message
				errMsg := ""
				switch e := err.(type) {
				case error:
					errMsg = e.Error()
				case string:
					errMsg = e
				default:
					errMsg = "Unknown error occurred"
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": errMsg})
			}
		}()

		c.Next()
	}
}

// NewServer creates a new HTTP server
func NewServer() *Server {
	gin.SetMode(gin.ReleaseMode)

	// Create a new router without default middleware
	router := gin.New()

	// Add the logger and recovery middleware
	router.Use(gin.Logger())

	// Add our custom global exception handler
	router.Use(globalExceptionHandler())

	server := &Server{
		router: router,
	}

	// Set up routes
	server.setupRoutes()

	return server
}

// setupRoutes sets up the HTTP routes
func (s *Server) setupRoutes() {
	s.router.GET("/api/generate", s.handleGenerate)
	s.router.GET("/api/quickstart/*url", s.handleQuickstart)
}

// Run starts the HTTP server
func (s *Server) Run() error {
	port, b := os.LookupEnv("SERVER_PORT")
	if !b {
		port = "5000"
	}
	slog.Info("Starting server", "port", port)
	return s.router.Run(":" + port)
}

// handleGenerate handles the /api/generate endpoint
func (s *Server) handleGenerate(c *gin.Context) {
	// Get template file from query parameter
	templateFile := c.Query("file")
	if templateFile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing template file parameter"})
		return
	}

	// 获取配置模板
	configs, err := template.GetConfigTemplate(templateFile)
	if err != nil {
		slog.Error("Failed to load template", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to load template"})
		return
	}

	// Get configuration
	cfg := config.GetConfig()

	// 处理订阅
	handleGenerateConfigForSubscription(c, configs, cfg.Subscribes)
}

func handleGenerateConfigForSubscription(c *gin.Context, configs map[string]any, subscribes []config.Subscription) {
	nodes, err := converter.ProcessSubscribes(subscribes)
	if err != nil {
		slog.Error("Failed to process subscribes", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process subscribes"})
		return
	}

	// 节点信息添加到模板
	finalConfig, err := template.MergeToConfig(configs, nodes)
	if err != nil {
		slog.Error("Failed to merge config", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to merge config"})
		return
	}
	c.JSON(http.StatusOK, finalConfig)
}

// handleQuickstart handles the /api/quickstart endpoint
func (s *Server) handleQuickstart(c *gin.Context) {
	// Extract subscription URL from path
	fullPath := c.Param("url")
	if fullPath == "" || fullPath == "/" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing subscription URL"})
		return
	}

	// Remove leading slash
	subURL := fullPath[1:]

	// Get template file from query parameter
	templateFile := c.Query("file")
	if templateFile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing template file parameter"})
		return
	}

	// 获取配置模板
	configs, err := template.GetConfigTemplate(templateFile)
	if err != nil {
		slog.Error("Failed to load template", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to load template"})
		return
	}

	subscribes := make([]config.Subscription, 0)
	subscribes = append(subscribes, config.Subscription{
		URL:       subURL,
		Tag:       "single",
		Prefix:    "",
		UserAgent: "",
	})
	// 处理订阅
	handleGenerateConfigForSubscription(c, configs, subscribes)
}
