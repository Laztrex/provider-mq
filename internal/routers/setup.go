package routers

import (
	"github.com/gin-gonic/gin"
)

// SetupRouters Function to setup routers and router groups
func SetupRouters(WebApp *gin.Engine) {
	WebApp.GET("/health", Health)
}
