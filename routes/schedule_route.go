package routes

import (
	"github.com/Amierza/chat-service/handler"
	"github.com/Amierza/chat-service/jwt"
	"github.com/Amierza/chat-service/middleware"
	"github.com/gin-gonic/gin"
)

func Schedule(route *gin.Engine, scheduleHandler handler.IScheduleHandler, jwt jwt.IJWT) {
	routes := route.Group("/api/v1/schedules").Use(middleware.Authentication(jwt))
	{
		routes.POST("", scheduleHandler.Create)
		routes.GET("", scheduleHandler.GetAll)
		routes.GET("/:id", scheduleHandler.GetDetail)
		routes.PUT("/:id", scheduleHandler.Update)
		routes.POST("/:id/approval", scheduleHandler.Approval)
		routes.DELETE("/:id", scheduleHandler.Delete)
	}
}
