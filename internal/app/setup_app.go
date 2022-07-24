//go:build exclude

package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"provider_mq/internal/middlewares"
	"provider_mq/internal/routers"
)

// SetupApp Function to setup the app object
func SetupApp() *gin.Engine {
	log.Info().Msg("Initializing service")

	// Create barebone engine
	app := gin.New()
	// Add default recovery middleware
	app.Use(gin.Recovery())

	// disabling the trusted proxy feature
	app.SetTrustedProxies(nil)

	// Add cors, request ID and request logging middleware
	log.Info().Msg("Adding cors, request id and request logging middleware")
	app.Use(middlewares.RequestID(), middlewares.RequestLogger())

	// Setup routers
	log.Info().Msg("Setting up routers")
	routers.SetupRouters(app)

	return app
}
