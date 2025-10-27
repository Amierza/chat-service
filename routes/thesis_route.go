package routes

import (
	"github.com/Amierza/chat-service/handler"
	"github.com/Amierza/chat-service/jwt"
	"github.com/Amierza/chat-service/middleware"
	"github.com/gin-gonic/gin"
)

func Thesis(route *gin.Engine, thesisHandler handler.IThesisHandler, jwt jwt.IJWT) {
	routes := route.Group("/api/v1/theses").Use(middleware.Authentication(jwt))
	{
		routes.GET("/:id", thesisHandler.GetDetail)
		routes.PUT("/:id", thesisHandler.Update)

		routes.GET("/lecturer/:lecturer_id", thesisHandler.GetAllByLecturer)
	}
}
