package handlers

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/proxy/debug"
	"github.com/kpn/pion/pkg/pion/proxy/utils"
	"github.com/labstack/echo"
)

func HandleRequest(c echo.Context) error {
	var printBody = false
	if glog.V(4) {
		printBody = true
	}

	traceOutput := new(bytes.Buffer)
	traceOutput.Grow(1024)

	originalRequest := c.Request()
	fmt.Fprintln(traceOutput, "---------START-ORIGINAL-HTTP---------")
	fmt.Fprintln(traceOutput, debug.FormatRequest(originalRequest, printBody))
	fmt.Fprintln(traceOutput, "---------END-ORIGINAL-HTTP---------")
	glog.V(3).Info(traceOutput.String())

	glog.V(3).Info("Resigning with master key")

	populatedRequest, err := PopulateRequest(originalRequest)
	if err != nil {
		glog.Errorf("failed to populate request: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	traceOutput.Reset()
	fmt.Fprintln(traceOutput, "---------START-UPSTREAM-HTTP---------")
	fmt.Fprintln(traceOutput, debug.FormatRequestOut(populatedRequest, printBody))

	glog.V(2).Info("Forwarding to upstream")
	upstreamResp, err := sendRequest(populatedRequest)
	if err != nil {
		glog.Errorf("Failed to forward to upstream: %v", err)
		glog.V(3).Info(traceOutput.String())
		return c.NoContent(http.StatusBadGateway)
	}
	glog.V(2).Infof("Response from upstream:\n%s", debug.FormatResponse(upstreamResp, printBody))

	body := upstreamResp.Body
	defer func() {
		err = body.Close()
		if err != nil {
			glog.Errorf("[proxy error]: %v", err)
		}
	}()

	response := c.Response()
	utils.CloneHeaders(response.Header(), upstreamResp.Header)

	err = c.Stream(upstreamResp.StatusCode, upstreamResp.Header.Get(echo.HeaderContentType), body)
	if err != nil {
		glog.Errorf("Failed to copy response: %v", err)
		return c.NoContent(http.StatusBadGateway)
	}

	// log the request and response code
	glog.V(1).Infof(`"%s" "%s %s" %d "%s"`,
		originalRequest.Host,
		originalRequest.Method,
		originalRequest.RequestURI,
		upstreamResp.StatusCode,
		originalRequest.UserAgent())
	return nil
}
