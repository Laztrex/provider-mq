package app

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"provider_mq/internal/routers"
)

// SetupApp Function to setup the app object
func SetupApp() *http.ServeMux {
	log.Info().Msg("Initializing service")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", routers.Health)

	err := http.ListenAndServe(":5080", mux)
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("error running http server: %s\n", err)
		}
	}

	return mux

}
