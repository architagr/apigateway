package router

import (
	"apigateway/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, handler *handler.Handler) {
	r.GET("/hello", handler.SayHelloHandler)
}
