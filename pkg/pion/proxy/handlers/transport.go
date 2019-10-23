package handlers

import (
	"net"
	"net/http"
	"os"
	"time"
)

var UpstreamAddress = os.Getenv("MINIO_SERVICE_URL")

var myTransport *http.Transport = nil

func sendRequest(req *http.Request) (*http.Response, error) {
	return getTransport().RoundTrip(req)
}

func getTransport() *http.Transport {
	if myTransport == nil {
		myTransport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       10 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	}

	return myTransport
}
