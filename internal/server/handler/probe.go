package handler

import (
	"fmt"

	"git.h2hsecure.com/ddos/waf/internal/core/domain"
	"github.com/gin-gonic/gin"
)

/*
(function() {
function probeDDoS() {
    var xhr = new XMLHttpRequest();
    xhr.open("GET", "/ddos/probe", true);
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 302) {
            var newLocation = xhr.getResponseHeader("Location");
            if (newLocation) {
                window.location.replace(newLocation);
            }
        }
    };
    xhr.send();
}

    probeDDoS(); // Run immediately
    setInterval(probeDDoS, 60000);

})();
*/

var (
	js_func = `
!function(){function e(){var e=new XMLHttpRequest;e.open("GET","%s",!0),e.onreadystatechange=function(){if(4===e.readyState&&302===e.status){var n=e.getResponseHeader("Location");n&&window.location.replace(n)}},e.send()}e(),setInterval(e,%d)}();
	`
)

func NewProbeHandler(c *gin.Engine, cfg domain.NginxParams) error {
	c.GET("/"+cfg.ContextPath+"/probe.js", probeHandler(cfg))

	return nil
}

func probeHandler(cfg domain.NginxParams) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/javascript")
		c.Writer.WriteString(fmt.Sprintf(js_func, "/what_am_i_doing_in_here", 60000))
	}
}
