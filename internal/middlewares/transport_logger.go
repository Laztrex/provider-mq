package middlewares

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	stdlog "log"
	"net/http"
	"os"
	"provider_mq/internal/consts"
	"time"
)

// LoggerTransportRoundTripper This type implements the http.RoundTripper interface
type LoggerTransportRoundTripper struct {
	Proxy http.RoundTripper
}

func (lrt LoggerTransportRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	var requestId, msgRes string

	requestId = req.Header.Get(consts.RequestIdHttpHeaderName)
	if requestId == "" {
		requestId = uuid.New().String()
		req.Header.Set(consts.RequestIdHttpHeaderName, requestId)
	}

	method := req.Method
	path := req.URL.Path
	t := time.Now()

	fmt.Printf("Sending request to %v\n", req.URL)

	res, e = lrt.Proxy.RoundTrip(req)

	latency := float32(time.Since(t).Seconds())

	if e != nil || res.StatusCode != 200 {
		msgRes = "FAILURE"
	} else {
		msgRes = "SUCCESS"
	}

	stdlog.Printf("%s %s: '%s' %f - [%s]", msgRes, path, res.Status, latency, requestId)
	logToFile(res.StatusCode, requestId, method, path, latency, msgRes)

	return
}

func logToFile(status int, requestId string, method string, path string, latency float32, msg string) {
	tempFile, err := os.OpenFile(consts.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(err).Msg("there was an error creating a temporary file four our log")
	}

	fileLogger := zerolog.New(tempFile)
	fileLogger.Info().
		Int("status", status).
		Str("request_id", requestId).
		Str("method", method).
		Str("path", path).
		Float32("latency", latency).
		Msg(msg)
}
