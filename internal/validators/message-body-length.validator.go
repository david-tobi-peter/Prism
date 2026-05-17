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

	return nil
}
