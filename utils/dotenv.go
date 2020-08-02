package utils

import (
	"os"

	"github.com/gin-gonic/gin"
)

func LoadEnvVariable() string {
	env := ".env"
	if mode := os.Getenv("APP_MODE"); mode == "production" {
		env += ".production"
		gin.SetMode(gin.ReleaseMode)
		return env
	} else {
		env += ".production"
		return env
	}
}
