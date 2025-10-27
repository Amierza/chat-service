package routes

import (
	"github.com/Amierza/chat-service/handler"
	"github.com/Amierza/chat-service/jwt"
	"github.com/Amierza/chat-service/middleware"
	"github.com/gin-gonic/gin"
)

func File(route *gin.Engine, fileHandler handler.IFileHandler, jwt jwt.IJWT) {
	routes := route.Group("/api/v1/uploads").Use(middleware.Authentication(jwt))
	{
		routes.POST("", fileHandler.UploadFiles)
	}
}
