package server

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"git.h2hsecure.com/ddos/waf/internal/core/ports"
	"git.h2hsecure.com/ddos/waf/internal/repository/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const COOKIE_NAME = "ddos-cookei"
const COOKIE_DURATION = 3600

//go:embed ui
var content embed.FS

type nginxHandler struct {
	cache                   ports.Cache
	mq                      ports.MessageQueue
	contextPath             string
	dispatcher              *domain.Dispatcher
	disableProcessing       bool
	enableSearchEngineBoots bool
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

func CreateHttpServer(memcache ports.Cache, messageQueue ports.MessageQueue) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	mux := gin.New()
	mux.Use(gin.Recovery())
	contextPath := os.Getenv("CONTEXT_PATH")
	dispatcher := domain.NewDispatcher(10, 100)

	_, disableProcessing := os.LookupEnv("DISABLE_PROCESSING")
	_, enableSearchEngineBoots := os.LookupEnv("ENABLE_SEARCH_ENGINE_BOTS")

	handler := nginxHandler{
		cache:                   memcache,
		mq:                      messageQueue,
		contextPath:             os.Getenv("CONTEXT_PATH"),
		dispatcher:              dispatcher,
		disableProcessing:       disableProcessing,
		enableSearchEngineBoots: enableSearchEngineBoots,
	}

	dispatcher.Run()

	mux.GET("/"+contextPath+"/auth", handler.authzHandler)
	mux.GET("/"+contextPath+"/check", handler.checkHandler)
	mux.GET("/"+contextPath+"/health", handler.healthHandler)

	dist, err := fs.Sub(content, "ui")
	if err != nil {
		log.Err(err).Msg("sub error")
	}

	mux.StaticFS("/"+contextPath+"/app/", http.FS(dist))

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

	v, err := c.Cookie(COOKIE_NAME)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			redirect(c, "", requestUri, http.StatusUnauthorized)
		default:
			log.Err(err).Send()
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	t, err := token.VerifyToken(v)

	if err != nil {
		log.Err(err).Send()
		redirect(c, "", requestUri, http.StatusUnauthorized)
		return
	}

	event.User = t.UserId
	event.Ip = t.Ip

	if _, has := os.LookupEnv("DISABLE_ENFORCING"); has {
		c.Status(http.StatusOK)
		return
	}

	last, err := n.cache.Get(c, t.UserId)

	if err != nil {
		log.Err(err).Msg("cache get sub")
	}

	if last != "" {
		log.Warn().Str("sub", last).Msg("user found in cache")
		redirect(c, "forbiden.html", requestUri, http.StatusForbidden)
		return
	}

	last, err = n.cache.Get(c, t.Ip)

	if err != nil {
		log.Err(err).Msg("cache get ip")
	}

	if last != "" {
		log.Warn().Str("ip", last).Msg("ip found in cache")
		redirect(c, "forbiden.html", requestUri, http.StatusForbidden)
		return
	}

	c.Status(http.StatusOK)
}

func (n *nginxHandler) checkHandler(c *gin.Context) {

	log.Info().
		Interface("header", c.Request.Header).
		Str("path", c.Request.URL.Path).
		Send()

	id, _ := uuid.NewRandom()
	ip := c.Request.Header.Get("X-Real-Ip")

	token, err := token.CreateToken(id.String(), ip, time.Duration(COOKIE_DURATION))

	if err != nil {
		log.Err(err).Msg("create token")
	}

	c.SetCookie(COOKIE_NAME, token, COOKIE_DURATION, "/", os.Getenv("DOMAIN"), true, false)
	c.Status(http.StatusOK)
}

func (n *nginxHandler) healthHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}

func redirect(c *gin.Context, location, referer string, statusCode int) {
	c.Writer.Header().Add("Location", path(location))
	c.Writer.Header().Add("Referer", referer)
	c.AbortWithStatus(statusCode)
}

func path(suffix string) string {
	return fmt.Sprintf("%s://%s/%s/app/%s",
		os.Getenv("DOMAIN_PROTO"),
		os.Getenv("DOMAIN"),
		os.Getenv("CONTEXT_PATH"),
		suffix)
}
