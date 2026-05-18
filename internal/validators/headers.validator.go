package validators

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func ValidateMessageBodyLength(req *http.Request) error {
	hasChunked := false
	for _, te := range req.Header["Transfer-Encoding"] {
		if strings.Contains(strings.ToLower(te), "chunked") {
			hasChunked = true
			break
		}
	}

	if hasChunked {
		req.Header.Del("Content-Length")
		return nil
	}

	clHeaders := req.Header["Content-Length"]
	if len(clHeaders) > 1 {
		firstCL := clHeaders[0]
		for _, cl := range clHeaders {
			if cl != firstCL {
				return errors.New("Conflicting Content-Length Headers")
			}
		}
	}

	if len(clHeaders) == 1 {
		if val, err := strconv.Atoi(clHeaders[0]); err != nil || val < 0 {
			return errors.New("Invalid Content-Length Header Value")
		}
	}

	for key := range req.Header {
		if strings.ContainsAny(key, " \t\r\n:") {
			return errors.New("Invalid Header Field Name")
		}
	}

	return nil
}

func StripHopByHopHeaders(req *http.Request) {
	var standardHopHeaders = []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailer",
		"Transfer-Encoding",
		"Upgrade",
	}

	if connectionHeader := req.Header.Get("Connection"); connectionHeader != "" {
		for token := range strings.SplitSeq(connectionHeader, ",") {
			trimmedToken := strings.TrimSpace(token)
			if trimmedToken != "" {
				req.Header.Del(trimmedToken)
			}
		}
	}

	for _, hopHeader := range standardHopHeaders {
		req.Header.Del(hopHeader)
	}
}
