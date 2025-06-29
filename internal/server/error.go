package server

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

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
