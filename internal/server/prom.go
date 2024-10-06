package server

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func CreatePromServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	mux := gin.New()
	mux.Use(gin.Recovery())

	mux.GET("/metrics", func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})

	return mux
}
