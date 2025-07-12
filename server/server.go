package server

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sing-box-sub-converter/config"
	"sing-box-sub-converter/converter"
	"sing-box-sub-converter/template"
	"strings"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(globalExceptionHandler())
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "invalid route"})
	})

	server := &Server{
		router: router,
	}

	// Set up routes
	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	// Serve the embedded index.html file
	s.router.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(indexHtml))
	})
	s.router.GET("/favicon.ico", func(c *gin.Context) {
		c.Data(http.StatusOK, "image/x-icon", favicon)
	})

	// API routes
	s.router.GET("/api/generate", s.handleGenerate)
	s.router.GET("/api/quickstart/*url", s.handleQuickstart)
}

func (s *Server) Run() error {
	port, b := os.LookupEnv("SERVER_PORT")
	if !b {
		port = "5000"
	}
	slog.Info("正在启动服务器", "port", port)
	return s.router.Run(":" + port)
}

func (s *Server) handleGenerate(c *gin.Context) {
	templateFile := c.Query("file")
	if templateFile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing template file parameter"})
		return
	}

	// 获取配置模板
	configs, err := template.GetConfigTemplate(templateFile)
	if err != nil {
		slog.Error("加载模板失败", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to load template"})
		return
	}

	cfg := config.GetConfig()

	// 处理订阅
	handleGenerateConfigForSubscription(c, configs, cfg.Subscribes)
}

func handleGenerateConfigForSubscription(c *gin.Context, configs map[string]any, subscribes []config.Subscription) {
	nodes, err := converter.ProcessSubscribes(subscribes)
	if err != nil {
		slog.Error("处理订阅失败", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process subscribes"})
		return
	}

	// 节点信息添加到模板
	finalConfig, err := template.MergeToConfig(configs, nodes)
	if err != nil {
		slog.Error("合并配置失败", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to merge config"})
		return
	}
	c.JSON(http.StatusOK, finalConfig)
}

func (s *Server) handleQuickstart(c *gin.Context) {
	fullPath := c.Param("url")
	if fullPath == "" || fullPath == "/" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing subscription URL"})
		return
	}

	subURL := fullPath[1:]

	if strings.HasPrefix(subURL, "file://") {
		filePath := strings.TrimPrefix(subURL, "file://")
		safeDir := os.Getenv("SAFE_DIR")
		if safeDir != "" {
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file path"})
				return
			}

			safeDirAbs, err := filepath.Abs(safeDir)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid SAFE_DIR configuration"})
				return
			}

			if !strings.HasPrefix(absPath, safeDirAbs) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: file is outside of safe directory"})
				return
			}
		}
	}

	templateFile := c.Query("file")
	if templateFile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing template file parameter"})
		return
	}

	// 获取配置模板
	configs, err := template.GetConfigTemplate(templateFile)
	if err != nil {
		slog.Error("加载模板失败", "error", err)
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
