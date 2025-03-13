package handler

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"git.h2hsecure.com/ddos/waf/cmd"
	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"git.h2hsecure.com/ddos/waf/internal/core/ports"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:embed ui
var uiContent embed.FS

type nginxHandler struct {
	cache                   ports.Cache
	mq                      ports.MessageQueue
	token                   ports.TokenService
	contextPath             string
	workingDomain           string
	dispatcher              *domain.Dispatcher
	disableProcessing       bool
	enableSearchEngineBoots bool
	cookieName              string
	basePath                string
}

type innerJob struct {
	event domain.UserIpTime
	mq    ports.MessageQueue
}

func (i *innerJob) Send(ctx context.Context) error {
	return i.mq.Sent(ctx, i.event)
}

func (n *nginxHandler) push(event domain.UserIpTime) {
	n.dispatcher.Push(&innerJob{
		event: event,
		mq:    n.mq,
	})
}

func CreateNginxAdapter(memcache ports.Cache, messageQueue ports.MessageQueue, tokenService ports.TokenService, cfg cmd.ConfigParams) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	mux := gin.New()
	mux.Use(gin.Recovery())
	dispatcher := domain.NewDispatcher(10, 100)

	handler := nginxHandler{
		cache:                   memcache,
		mq:                      messageQueue,
		token:                   tokenService,
		contextPath:             cfg.Nginx.ContextPath,
		dispatcher:              dispatcher,
		disableProcessing:       cfg.User.DisableProcessing,
		enableSearchEngineBoots: cfg.User.SearchEngineBots,
		workingDomain:           cfg.Nginx.BackendHost,
		cookieName:              cfg.User.CookieName,
		basePath: fmt.Sprintf("%s://%s/%s/app",
			cfg.Nginx.DomainProto,
			cfg.Nginx.Domain,
			cfg.Nginx.ContextPath),
	}

	dispatcher.Run()

	mux.GET("/"+cfg.Nginx.ContextPath+"/auth", handler.authzHandler)
	mux.GET("/"+cfg.Nginx.ContextPath+"/health", handler.healthHandler)
	mux.GET("/"+cfg.Nginx.ContextPath+"/waidih", handler.WaidihHandler)

	dist, err := fs.Sub(uiContent, "ui")
	if err != nil {
		log.Err(err).Msg("sub error")
	}

	mux.StaticFS("/"+cfg.Nginx.ContextPath+"/app/", http.FS(dist))

	return mux
}

func (n *nginxHandler) authzHandler(c *gin.Context) {
	requestUri := c.Request.Header.Get("X-Original-Uri")

	for _, t := range domain.StaticTypes {
		if strings.HasSuffix(requestUri, t) {
			c.Status(http.StatusOK)
			return
		}
	}

	// in some case, we don't want to process our enforcer etc.
	if n.disableProcessing {
		c.Status(http.StatusOK)
		return
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

	event := domain.UserIpTime{
		Ip:        ip,
		Path:      requestUri,
		Timestamp: int32(now.Unix()),
	}

	defer n.push(event)

	v, err := c.Cookie(n.cookieName)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			redirect(c, n.path(""), requestUri, http.StatusUnauthorized)
		default:
			log.Err(err).Send()
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	t, err := n.token.VerifyToken(v)

	if err != nil {
		log.Err(err).Send()
		redirect(c, n.path(""), requestUri, http.StatusUnauthorized)
		return
	}

	event.User = t.UserId
	event.Ip = t.Ip

	if n.disableProcessing {
		c.Status(http.StatusOK)
		return
	}

	last, err := n.cache.Get(c, t.UserId)

	if err != nil {
		log.Err(err).Msg("cache get sub")
	}

	if last != "" {
		log.Warn().Str("sub", last).Msg("user found in cache")
		// ask again for verfying as a human
		redirect(c, n.path("/index.html"), requestUri, http.StatusForbidden)
		return
	}

	last, err = n.cache.Get(c, t.Ip)

	if err != nil {
		log.Err(err).Msg("cache get ip")
	}

	if last != "" {
		log.Warn().Str("ip", last).Msg("ip found in cache")
		redirect(c, n.path("forbiden.html"), requestUri, http.StatusForbidden)
		return
	}

	c.Status(http.StatusOK)
}

func (n *nginxHandler) healthHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (n *nginxHandler) WaidihHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}

func redirect(c *gin.Context, location, referer string, statusCode int) {
	c.Writer.Header().Add("Location", location)
	c.Writer.Header().Add("Referer", referer)
	c.AbortWithStatus(statusCode)
}

func (n *nginxHandler) path(suffix string) string {
	return fmt.Sprintf("%s/%s",
		n.basePath,
		suffix)
}
