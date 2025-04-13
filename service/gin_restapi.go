package service

import (
	"github.com/achmadnr21/bankapi/internal/middleware"
	"github.com/gin-gonic/gin"
)

type RESTapi struct {
	Router *gin.Engine
}

func (r *RESTapi) Init() *gin.RouterGroup {
	r.Router = gin.Default()
	r.Router.Use(middleware.RateLimiter)
	r.Router.NoRoute(middleware.NoRouteExists)
	return r.Router.Group("/api/v1")

}
