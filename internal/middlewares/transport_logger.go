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
	var requestId, resStatus, exception string

	requestId = req.Header.Get(consts.RequestIdHttpHeaderName)
	if requestId == "" {
		requestId = uuid.New().String()
		req.Header.Set(consts.RequestIdHttpHeaderName, requestId)
	}

	path := req.URL.Path
	t := time.Now()

	fmt.Printf("Sending request to %v\n", req.URL)

	res, e = lrt.Proxy.RoundTrip(req)

	latency := float32(time.Since(t).Seconds())

	if res == nil {
		log.Error().Msg("Service <Model Application> unavailable")
		return
	}

	if e != nil || res.StatusCode != 200 {
		resStatus = "FAILURE"
		exception = "error"
	} else {
		resStatus = "SUCCESS"
	}

	stdlog.Printf("%s %s: '%s' %f - [%s]", resStatus, path, res.Status, latency, requestId)
	logToFile(requestId, latency, exception, resStatus)

	return
}

func logToFile(requestId string, latency float32, exception string, resStatus string) {
	tempFile, err := os.OpenFile(consts.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(err).Msg("there was an error creating a temporary file four our log")
	}

	fileLogger := zerolog.New(tempFile)
	fileLogger.Log().
		Str("request_id", requestId).
		Str("status", resStatus).
		Float32("provider_mq_response_seconds", latency).
		Str("exception", exception).
		Msg(time.Now().Format(time.RFC3339))
}
