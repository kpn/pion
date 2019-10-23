package utils

import (
	"net/http"
	"strings"
)

func CloneHeaders(dst, src http.Header) {
	for k, vv := range src {
		dst[k] = vv
	}
}

func CopyWhiteListHeaders(src http.Header) http.Header {
	whiteListPrefixes := []string{
		// "x-amz-",
		"content-", "transfer-encoding", "host", "range",
	}
	dst := make(http.Header)
	for h, val := range src {
		normalizedHeader := strings.ToLower(h)
		matched := false
		for _, prefix := range whiteListPrefixes {
			if strings.HasPrefix(normalizedHeader, prefix) {
				matched = true
				break
			}
		}
		if matched {
			dst[h] = val
		}
	}
	return dst
}
