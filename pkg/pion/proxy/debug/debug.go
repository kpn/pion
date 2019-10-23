package debug

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/golang/glog"
)

// FormatRequest converts incoming request to string
func FormatRequest(r *http.Request, body bool) string {
	bytes, err := httputil.DumpRequest(r, body)
	if err != nil {
		return fmt.Sprintf("Cannot dump request: %v", err)
	}
	return string(bytes)
}

// FormatRequestOut converts out-going request to string
func FormatRequestOut(r *http.Request, body bool) string {
	bytes, err := httputil.DumpRequestOut(r, body)
	if err != nil {
		return fmt.Sprintf("Cannot dump request: %v", err)
	}
	return string(bytes)
}

func FormatResponse(resp *http.Response, body bool) string {
	bytes, err := httputil.DumpResponse(resp, body)
	if err != nil {
		glog.Errorf("Cannot dump response: %v", err)
		return ""
	}
	return string(bytes)
}
