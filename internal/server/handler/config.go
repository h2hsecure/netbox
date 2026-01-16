package handler

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/h2hsecure/netbox/internal/core/domain"
	"golang.org/x/text/language"
)

/*
(function() {
window._ntb_dds = %s
})();
*/

var (
	config_js_func = `
        (function() {window._ntb_dds = %s})();
	`
)

//go:embed l8n/*
var languageContent embed.FS

type translation struct {
	Main struct {
		H1 string `json:"h1"`
		H2 string `json:"h2"`
		At struct {
			Label     string `json:"label"`
			Error     string `json:"error"`
			Expired   string `json:"expired"`
			Footer    string `json:"footer"`
			Verified  string `json:"verified"`
			Verifying string `json:"verifying"`
			WaitAlert string `json:"waitAlert"`
		} `json:"at"`
		Timeout string `json:"timeout"`
		Footer  struct {
			System string `json:"system"`
			Per    string `json:"per"`
		} `json:"footer"`
	} `json:"main"`
}

type clientConfig struct {
	SystemId    string      `json:"id"`
	Hostname    string      `json:"h"`
	ReturnUrl   string      `json:"r"`
	Translation translation `json:"t"`
	Logo        string      `json:"l"`
}

type configHandlerImpl struct {
	clientConfig clientConfig
	tanslations  map[string]translation
	matcher      language.Matcher
}

func NewConfigHandler(c *gin.Engine, cfg domain.ConfigParams) error {
	contextPath := cfg.Nginx.ContextPath

	chi := configHandlerImpl{
		tanslations: make(map[string]translation),
	}

	var tags []language.Tag

	fs.WalkDir(languageContent, ".", func(path string, d fs.DirEntry, err2 error) error {
		if d.IsDir() {
			return nil
		}

		contentFile, err := languageContent.Open(path)

		if err != nil {
			return fmt.Errorf("opening file contents: %w", err)
		}

		var trans translation

		if err := json.NewDecoder(contentFile).Decode(&trans); err != nil {
			return fmt.Errorf("decoding: %w", err)
		}

		locale := strings.Split(strings.TrimPrefix(path, "l8n/"), ".")[0]

		chi.tanslations[locale] = trans
		tags = append(tags, language.MustParse(locale))

		return nil
	})

	slices.SortFunc(tags, func(a language.Tag, b language.Tag) int {
		if a.String() == cfg.User.DefaultLocale {
			return 1
		}

		if b.String() == cfg.User.DefaultLocale {
			return -1
		}

		return 0
	})

	chi.matcher = language.NewMatcher(tags)
	chi.clientConfig = clientConfig{
		Hostname: cfg.Nginx.Domain,
		SystemId: cfg.SystemId,
		Logo:     cfg.User.BackendLogo,
	}

	c.GET("/"+contextPath+"/ntb_dds", chi.configHandler)

	return nil
}

func (chi *configHandlerImpl) configHandler(c *gin.Context) {

	accept := c.Request.Header.Get("Accept-Language")
	referer := c.Request.Header.Get("Referer")
	tag, _ := language.MatchStrings(chi.matcher, accept)

	cc := chi.clientConfig

	cc.ReturnUrl = referer
	cc.Translation = chi.tanslations[tag.String()]

	c.Writer.Header().Set("Content-Type", "application/javascript")

	buf, _ := json.Marshal(cc)

	c.Writer.WriteString(fmt.Sprintf(config_js_func, string(buf)))
}
