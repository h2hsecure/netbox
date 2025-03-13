package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"git.h2hsecure.com/ddos/waf/internal/repository/token"
	"github.com/altcha-org/altcha-lib-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const COOKIE_NAME = "ddos-cookei"
const COOKIE_DURATION = 3600

var altchaHMACKey = os.Getenv("ALTCHA_HMAC_KEY")
var workingDomain = os.Getenv("DOMAIN")

func CreateHumanServer(mux *gin.Engine) error {
	contextPath := os.Getenv("CONTEXT_PATH")

	mux.GET("/"+contextPath+"/challenge", challengeHandler)
	mux.POST("/"+contextPath+"/accept", acceptHandler)

	return nil
}

func challengeHandler(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		http.Error(c.Writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	challenge, err := altcha.CreateChallenge(altcha.ChallengeOptions{
		HMACKey:   altchaHMACKey,
		MaxNumber: 50000,
	})
	if err != nil {
		http.Error(c.Writer, fmt.Sprintf("Failed to create challenge: %s", err), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, challenge)
}

func acceptHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		http.Error(c.Writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	formData := c.Request.FormValue("altcha")
	if formData == "" {
		http.Error(c.Writer, "Altcha payload missing", http.StatusBadRequest)
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

	verified, err := altcha.VerifySolution(formData, altchaHMACKey, true)
	if err != nil || !verified {
		http.Error(c.Writer, "Invalid Altcha payload", http.StatusBadRequest)
		return
	}

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

	c.SetCookie(COOKIE_NAME, token, COOKIE_DURATION, "/", workingDomain, true, false)
	c.Status(http.StatusOK)

	// For demo purposes, echo back the form data
	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                   // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Allow methods
		w.Header().Set("Access-Control-Allow-Headers", "*")                  // Allow headers

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
