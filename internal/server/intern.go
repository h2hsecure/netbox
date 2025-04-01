package server

import (
	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"github.com/gin-gonic/gin"
)

func CreateInternalServer(cfg domain.ConfigParams) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	mux := gin.New()
	mux.Use(gin.Recovery())

	return mux
}
