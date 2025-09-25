package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/order-nest/config"
)

// NewServer creates an *http.Server with proper timeouts and context
func NewHTTPServer(ctx context.Context, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", config.GetConfig().Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
	}
}
