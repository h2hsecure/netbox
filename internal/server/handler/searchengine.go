package handler

import "strings"

var searchEngineUserAgents = []string{
	"Googlebot",
}

func IsItSearchEngine(userAgent string) bool {

	for _, seua := range searchEngineUserAgents {
		if strings.Contains(userAgent, seua) {
			return true
		}
	}

	return false
}
