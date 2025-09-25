package handler

import "github.com/gin-gonic/gin"

// RouteDef describes an HTTP route in a declarative style.
type RouteDef struct {
	Method   string
	Path     string
	Handlers []gin.HandlerFunc
}
