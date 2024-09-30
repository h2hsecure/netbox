package server

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"git.h2hsecure.com/ddos/waf/internal/core/ports"
	"git.h2hsecure.com/ddos/waf/internal/repository/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const COOKIE_NAME = "ddos-cookei"

var cache ports.Cache

//go:embed ui
var content embed.FS

func staticWeb() http.FileSystem {
	dist, err := fs.Sub(content, "ui")
	if err != nil {
		log.Err(err).Msg("sub error")
	}

	return http.FS(dist)
}

func CreateHttpServer(port string, memcache ports.Cache) *gin.Engine {
	cache = memcache
	gin.SetMode(gin.ReleaseMode)
	mux := gin.Default()

	mux.GET("/ddos/auth", authzHandler)
	mux.GET("/ddos/check", checkHandler)
	mux.StaticFS("/ddos/app/", staticWeb())

	fmt.Printf("Server is running on port %s\n", port)
	if err := mux.Run(":" + port); err != nil {
		log.Err(err).Send()
	}

	return mux
}

func authzHandler(c *gin.Context) {
	log.Info().
		Interface("header", c.Request.Header).
		Str("path", c.Request.URL.Path).
		Msgf("auth request")

	v, err := c.Cookie(COOKIE_NAME)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			c.Writer.Header().Add("Location", "/ddos/app/")
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
		c.Writer.Header().Add("Location", "/ddos/app/")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sub, _ := t.Claims.GetSubject()

	last, err := cache.Inc(c, sub, 1)

	if err != nil {
		log.Warn().Interface("error", err).Send()

	}

	if last > 100 {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// user, err := cache.Get(context.Background(), v.Value)

	// if user == "" {
	// 	log.Printf("wrong user: %s %s", user, r.RemoteAddr)
	// 	w.Header().Add("Location", "/ddos/")
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	c.Status(http.StatusOK)
}

func checkHandler(c *gin.Context) {

	log.Info().
		Interface("header", c.Request.Header).
		Str("path", c.Request.URL.Path).
		Msgf("check request")

	id, _ := uuid.NewRandom()

	token, _ := token.CreateToken(domain.SessionCliam{
		UserId: id.String(),
	})

	// cookie := http.Cookie{
	// 	Name:     COOKIE_NAME,
	// 	Value:    token,
	// 	Path:     "/",
	// 	MaxAge:   3600,
	// 	HttpOnly: true,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteLaxMode,
	// }

	//err := cache.Set(context.Background(), id.String(), strings.Split(r.RemoteAddr, ":")[0])

	// if err != nil {
	// 	log.Err(err).Send()
	// 	http.Error(w, "server error", http.StatusInternalServerError)
	// 	return
	// }

	c.SetCookie(COOKIE_NAME, token, 3600, "/", "", true, true)
	cache.Inc(c, "online-count", 1)
	c.Status(http.StatusOK)
}
