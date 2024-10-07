package server

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"git.h2hsecure.com/ddos/waf/internal/core/ports"
	"git.h2hsecure.com/ddos/waf/internal/repository/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const COOKIE_NAME = "ddos-cookei"

//go:embed ui
var content embed.FS

type nginxHandler struct {
	cache       ports.Cache
	mq          ports.MessageQueue
	contextPath string
	dispatcher  *domain.Dispatcher
}

type innerJob struct {
	event domain.UserIpTime
	mq    ports.MessageQueue
}

func (i *innerJob) Send() error {
	return i.mq.Sent(context.Background(), i.event)
}

func CreateHttpServer(memcache ports.Cache, messageQueue ports.MessageQueue) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	mux := gin.New()
	mux.Use(gin.Recovery())
	contextPath := os.Getenv("CONTEXT_PATH")
	dispatcher := domain.NewDispatcher(10, 100)

	handler := nginxHandler{
		cache:       memcache,
		mq:          messageQueue,
		contextPath: os.Getenv("CONTEXT_PATH"),
		dispatcher:  dispatcher,
	}

	dispatcher.Run()

	mux.GET("/"+contextPath+"/auth", handler.authzHandler)
	mux.GET("/"+contextPath+"/check", handler.checkHandler)

	dist, err := fs.Sub(content, "ui")
	if err != nil {
		log.Err(err).Msg("sub error")
	}

	mux.StaticFS("/"+contextPath+"/app/", http.FS(dist))

	return mux
}

func path(suffix string) string {
	return fmt.Sprintf("%s://%s/%s/app/%s",
		os.Getenv("DOMAIN_PROTO"),
		os.Getenv("DOMAIN"),
		os.Getenv("CONTEXT_PATH"),
		suffix)
}
func (n *nginxHandler) authzHandler(c *gin.Context) {
	log.Info().
		Interface("header", c.Request.Header).
		Str("path", c.Request.URL.Path).
		Send()

	//contextPath := os.Getenv("CONTEXT_PATH")
	// if strings.HasPrefix(c.Request.URL.Path, "/"+contextPath+"/") {
	// 	c.Status(http.StatusOK)
	// 	return
	// }

	v, err := c.Cookie(COOKIE_NAME)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			c.Writer.Header().Add("Location", path(""))
			c.Writer.Header().Add("Referer", c.Request.Referer())
			c.AbortWithStatus(http.StatusUnauthorized)
		default:
			log.Err(err).Send()
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	t, err := token.VerifyToken(v)

	if err != nil {
		log.Err(err).Send()
		c.Writer.Header().Add("Location", path(""))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sub, _ := t.Claims.GetSubject()

	last, err := n.cache.Get(c, sub)

	if err != nil {
		log.Err(err).Msg("cache get sub")
	}

	if last != "" {
		log.Warn().Str("sub", last).Msg("user found in cache")
		c.Writer.Header().Add("Location", path("forbiden.html"))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	ip := c.Request.Header.Get("X-Real-IP")

	if ip == "" {
		ip = c.RemoteIP()
	}

	last, err = n.cache.Get(c, ip)

	if err != nil {
		log.Err(err).Msg("cache get ip")
	}

	if last != "" {
		log.Warn().Str("ip", last).Msg("ip found in cache")
		c.Writer.Header().Add("Location", path("forbiden.html"))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	n.dispatcher.Push(&innerJob{
		event: domain.UserIpTime{
			Ip:        ip,
			User:      sub,
			Timestamp: int32(time.Now().Unix()),
		},
		mq: n.mq,
	})

	c.Status(http.StatusOK)
}

func (n *nginxHandler) checkHandler(c *gin.Context) {

	log.Info().
		Interface("header", c.Request.Header).
		Str("path", c.Request.URL.Path).
		Send()

	id, _ := uuid.NewRandom()
	ip := c.Request.Header.Get("X-Real-Ip")

	token, _ := token.CreateToken(domain.WithDefaultCliam(id.String(), ip))

	c.SetCookie(COOKIE_NAME, token, 3600, "/", os.Getenv("DOMAIN"), true, false)
	// n.cache.Inc(c, "online-count", 1)
	// n.cache.Inc(c, id.String(), 1)
	c.Status(http.StatusOK)
}
