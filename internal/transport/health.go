package transport

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(`{"status":"OK"}`))
	w.WriteHeader(http.StatusOK)
	if err != nil {
		log.Error().Err(err).Msg("can't send a reply in Health method")
	}
}
