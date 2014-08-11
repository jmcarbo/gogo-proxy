package proxy

import (
	"net/http"
	"net/url"
	"strings"
)

// normalizeUrl try's to add a scheme to a url if doesn't any
func normalizeUrl(rawurl string) string {
	// default "://" to "http://"
	if strings.HasPrefix(rawurl, "://") {
		rawurl = strings.Replace(rawurl, "://", "http://", 1)
	}

	parsed, err := url.Parse(rawurl)
	if err != nil {
		return rawurl
	}

	// Cleanup or default scheme to http
	newScheme := httpScheme(parsed.Scheme)
	if newScheme != parsed.Scheme {
		// Use new scheme
		parsed.Scheme = newScheme

		// We need to reparse the URL because now that there is a prefix
		// the "Host" and "Path" fields are most likely going to change
		parsed, err = url.Parse(parsed.String())
		if err != nil {
			return rawurl
		}
	}

	// Default path
	if parsed.Path == "" {
		parsed.Path = "/"
	}

	// Return URL string
	return parsed.String()
}

// websocketScheme picks a suitable websocket scheme
func websocketScheme(scheme string) string {
	switch scheme {
	case "http":
		return "ws"
	case "https":
		return "wss"
	case "ws":
	case "wss":
		return scheme
	}
	// Default
	return "ws"
}

// httpScheme picks a suitable http scheme
func httpScheme(scheme string) string {
	switch scheme {
	case "ws":
		return "http"
	case "wss":
		return "https"
	case "http":
	case "https":
		return scheme
	}
	// Default
	return "http"
}

// isWebsocket checks wether the incoming request is a part of websocket handshake
func isWebsocket(req *http.Request) bool {
	if strings.ToLower(req.Header.Get("Upgrade")) != "websocket" ||
		!strings.Contains(strings.ToLower(req.Header.Get("Connection")), "upgrade") {
		return false
	}
	return true
}
