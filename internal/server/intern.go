package server

import (
	"github.com/gin-gonic/gin"
	"github.com/h2hsecure/netbox/internal/core/domain"
)

func CreateInternalServer(cfg domain.ConfigParams) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	mux := gin.New()
	mux.Use(gin.Recovery())

	return mux
}
