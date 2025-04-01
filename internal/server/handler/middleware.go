package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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

func loggimgMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(g *gin.Context) {
		log.Info().
			Interface("uri", g.Request.URL).
			Interface("headers", g.Request.Header).Send()
		g.Next()

	})
}
