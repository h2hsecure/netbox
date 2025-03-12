package handler

import (
	"fmt"
	"os"

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

func NewProbeHandler(c *gin.Engine) error {
	contextPath := os.Getenv("CONTEXT_PATH")

	c.GET("/"+contextPath+"/probe.js", probeHandler)

	return nil
}

func probeHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/javascript")
	c.Writer.WriteString(fmt.Sprintf(js_func, "/what_am_i_doing_in_here", 60000))
}
