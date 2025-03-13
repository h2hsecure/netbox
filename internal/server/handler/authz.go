package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/core/ports"
	"github.com/altcha-org/altcha-lib-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type humanServer struct {
	token          ports.TokenService
	cookieName     string
	cookieDuration time.Duration
	hmacKey        string
	domain         string
}

func CreateHumanServer(mux *gin.Engine, token ports.TokenService) error {
	contextPath := os.Getenv("CONTEXT_PATH")

	cookieDuration, err := time.ParseDuration(os.Getenv("COOKIE_DURATION"))

	if err != nil {
		return fmt.Errorf("cookie duration parse: %w", err)
	}

	hs := humanServer{
		token:          token,
		hmacKey:        os.Getenv("CHALLENGE_HMAC_KEY"),
		domain:         os.Getenv("DOMAIN"),
		cookieName:     os.Getenv("COOKIE_NAME"),
		cookieDuration: cookieDuration,
	}

	mux.Use(corsMiddleware())

	mux.GET("/"+contextPath+"/challenge", hs.challengeHandler)
	mux.POST("/"+contextPath+"/accept", hs.acceptHandler)

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

	// Decode the Base64 encoded payload
	// decodedPayload, err := base64.StdEncoding.DecodeString(formData)
	// if err != nil {
	// 	http.Error(c.Writer, "Failed to decode Altcha payload", http.StatusBadRequest)
	// 	return
	// }

	// // Unmarshal the JSON payload
	// var payload map[string]interface{}
	// if err := json.Unmarshal(decodedPayload, &payload); err != nil {
	// 	http.Error(c.Writer, "Failed to parse Altcha payload", http.StatusBadRequest)
	// 	return
	//}

	verified, err := altcha.VerifySolution(formData, h.hmacKey, true)
	if err != nil || !verified {
		http.Error(c.Writer, "Invalid challenge payload", http.StatusBadRequest)
		return
	}

	log.Info().
		Interface("header", c.Request.Header).
		Str("path", c.Request.URL.Path).
		Send()

	id, _ := uuid.NewRandom()
	ip := c.Request.Header.Get("X-Real-Ip")

	token, err := h.token.CreateToken(id.String(), ip, time.Duration(0))

	if err != nil {
		log.Err(err).Msg("create token")
	}

	c.SetCookie(h.cookieName, token, int(h.cookieDuration), "/", h.domain, true, false)
	c.Status(http.StatusOK)

	// For demo purposes, echo back the form data
	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
