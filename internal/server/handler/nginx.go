package handler

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"git.h2hsecure.com/ddos/waf/internal/core/ports"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

//go:embed ui
var uiContent embed.FS

type nginxHandler struct {
	service                 ports.Service
	contextPath             string
	workingDomain           string
	disableProcessing       bool
	enableSearchEngineBoots bool
	cookieName              string
	basePath                string
	locationHeaderKey       string
	locationHeaderIsSet     bool
}

func CreateNginxAdapter(g *gin.Engine, service ports.Service, cfg domain.ConfigParams) error {

	handler := nginxHandler{
		service:     service,
		contextPath: cfg.Nginx.ContextPath,

		enableSearchEngineBoots: cfg.User.SearchEngineBots,
		workingDomain:           cfg.Nginx.BackendHost,
		cookieName:              cfg.User.CookieName,
		locationHeaderKey:       cfg.User.CountryHeader,
		locationHeaderIsSet:     cfg.User.CountryHeader != "",
		basePath: fmt.Sprintf("%s://%s/%s/app",
			cfg.Nginx.DomainProto,
			cfg.Nginx.Domain,
			cfg.Nginx.ContextPath),
	}

	g.GET("/"+cfg.Nginx.ContextPath+"/auth", handler.authzHandler)
	g.GET("/"+cfg.Nginx.ContextPath+"/health", handler.healthHandler)
	g.GET("/"+cfg.Nginx.ContextPath+"/waidih", handler.WaidihHandler)

	r := g.Use(loggimgMiddleware())

	dist, err := fs.Sub(uiContent, "ui")
	if err != nil {
		log.Err(err).Msg("sub error")
	}

	r.StaticFS("/"+cfg.Nginx.ContextPath+"/app/", http.FS(dist))
	return nil
}

func (n *nginxHandler) authzHandler(c *gin.Context) {
	requestUri := c.Request.Header.Get("X-Original-Uri")

	for _, t := range domain.StaticTypes {
		if strings.HasSuffix(requestUri, t) {
			c.Status(http.StatusOK)
			return
		}
	}

	// allow search engine bots
	if n.enableSearchEngineBoots {
		agent := c.Request.Header.Get("user-agent")

		if IsItSearchEngine(agent) {
			c.Status(http.StatusOK)
			return
		}
	}

	now := time.Now()
	log.Info().
		Interface("header", c.Request.Header).
		Time("when", now).
		Str("path", c.Request.URL.Path).
		Send()

	ip := c.Request.Header.Get("X-Real-IP")

	if ip == "" {
		ip = c.RemoteIP()
	}

	event := domain.AttemptRequest{
		UserIpTime: domain.UserIpTime{
			Ip:        ip,
			Path:      requestUri,
			Timestamp: int64(now.Unix()),
		},
		Location: nil,
	}

	if n.locationHeaderIsSet {
		location := c.Request.Header.Get(n.locationHeaderKey)

		if location != "" {
			event.Location = &location
		} else {
			event.Location = lo.ToPtr(string(domain.CountryPolicyAll))
			log.Warn().
				Str("header key", n.locationHeaderKey).
				Msg("location header is not found in request")
		}
	}

	v, err := c.Cookie(n.cookieName)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			v = ""
		default:
			log.Err(err).Send()
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}

	op := n.service.AccessAtempt(c, v, event)

	log.Info().Msgf("attempt result: %d", op)
	switch op {
	case domain.AttemptUserAllow:
		c.Status(http.StatusOK)
	case domain.AttemptDenyUserByCountry:
		redirect(c, "forbiden.html", requestUri, http.StatusForbidden)
	case domain.AttemptDenyUserByIp:
		redirect(c, "forbiden.html", requestUri, http.StatusForbidden)
	case domain.AttemptValidate:
		redirect(c, "", requestUri, http.StatusForbidden)
	}
}

func (n *nginxHandler) healthHandler(c *gin.Context) {
	err := n.service.Health(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func (n *nginxHandler) WaidihHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}

func redirect(c *gin.Context, location, referer string, statusCode int) {
	c.Writer.Header().Add("X-Location", location)
	c.Writer.Header().Add("X-Referer", referer)
	c.AbortWithStatus(statusCode)
}

func (n *nginxHandler) path(suffix string) string {
	return fmt.Sprintf("%s/%s",
		n.basePath,
		suffix)
}
