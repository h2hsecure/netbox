package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(g *gin.Context) {
		g.Writer.Header().Set("Access-Control-Allow-Origin", "*")                   // Allow all origins
		g.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Allow methods
		g.Writer.Header().Set("Access-Control-Allow-Headers", "*")                  // Allow headers

		if g.Request.Method == http.MethodOptions {
			g.Writer.WriteHeader(http.StatusOK)
			return
		}
	})
}
