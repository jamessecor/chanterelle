package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FrontendHandler serves the React frontend and handles client-side routing
func FrontendHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First try to serve the requested file
		c.FileFromFS(c.Request.URL.Path, http.Dir("/app/frontend/dist"))
		
		// If file not found, serve index.html for client-side routing
		if c.Writer.Status() == http.StatusNotFound {
			c.FileFromFS("/app/frontend/dist/index.html", http.Dir("/app/frontend/dist"))
		}
	}
}
