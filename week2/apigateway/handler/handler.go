package handler

import (
	"net/http"

	"apigateway/pkg/forwarder"
	"apigateway/pkg/logger"
	service1 "apigateway/pkg/proto/service1"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type 	Handler struct {
	forwarder *forwarder.Forwarder
	logger    *logger.ZapLogger
}

func NewHandler(f *forwarder.Forwarder, l *logger.ZapLogger) *Handler {
	return &Handler{
		forwarder: f,
		logger: l,
	}
}

func (h *Handler) SayHelloHandler(c *gin.Context) {

	//get a request-scoped logger with trace ID
	reqLogger := h.logger.WithTrace(c.Request.Context())

	name := c.Query("name")

	request := &service1.HelloRequest{Name: name}

	resp, err := h.forwarder.Service1Client.SayHello(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		reqLogger.Error("request failed", zap.Error(err))
		return
	}

	//return response
	c.JSON(http.StatusOK, gin.H{"reply": resp})
}
