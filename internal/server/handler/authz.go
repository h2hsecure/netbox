package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/altcha-org/altcha-lib-go"
	"github.com/gin-gonic/gin"
	"github.com/h2hsecure/netbox/internal/core/domain"
	"github.com/h2hsecure/netbox/internal/core/ports"
	"github.com/rs/zerolog/log"
)

type humanServer struct {
	service        ports.Service
	cookieName     string
	cookieDuration time.Duration
	hmacKey        string
	domain         string
}

func CreateHumanServer(mux *gin.Engine, serivce ports.Service, cfg domain.ConfigParams) error {

	hs := humanServer{
		service:        serivce,
		hmacKey:        cfg.User.ChallengeHmac,
		domain:         cfg.Nginx.Domain,
		cookieName:     cfg.User.CookieName,
		cookieDuration: cfg.User.CookieDuration,
	}

	mux.Use(corsMiddleware())

	mux.GET("/"+cfg.Nginx.ContextPath+"/challenge", hs.challengeHandler)
	mux.POST("/"+cfg.Nginx.ContextPath+"/accept", hs.acceptHandler)

	return nil
}

func (h *humanServer) challengeHandler(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		http.Error(c.Writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	challenge, err := altcha.CreateChallenge(altcha.ChallengeOptions{
		HMACKey:   h.hmacKey,
		MaxNumber: 50000,
	})
	if err != nil {
		http.Error(c.Writer, fmt.Sprintf("Failed to create challenge: %s", err), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, challenge)
}

func (h *humanServer) acceptHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		http.Error(c.Writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	formData := c.Request.FormValue("challenge")
	if formData == "" {
		http.Error(c.Writer, "challenge payload missing", http.StatusBadRequest)
		return
	}

	verified, err := altcha.VerifySolution(formData, h.hmacKey, true)
	if err != nil || !verified {
		http.Error(c.Writer, "Invalid challenge payload", http.StatusBadRequest)
		return
	}

	log.Info().
		Interface("header", c.Request.Header).
		Str("path", c.Request.URL.Path).
		Send()

	token, err := h.service.OpenSession(c, domain.UserIpTime{
		Ip:   c.Request.Header.Get("X-Real-Ip"),
		Path: c.Request.URL.String(),
	})

	if err != nil {
		log.Err(err).Msg("open session")
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"err":     err.Error(),
		})
		return
	}

	c.SetCookie(h.cookieName, token, int(h.cookieDuration), "/", h.domain, true, false)

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
