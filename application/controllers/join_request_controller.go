package controllers

import (
	"sport_platform/application/handlers"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func JoinRequestController(engine *gin.Engine, wrapper *service_wrapper.Wrapper) {
	routerGroup := engine.Group("/club_join_requests")
	routerGroup.POST("/create", func(context *gin.Context) {
		handlers.CreateJoinRequestHandler(context, wrapper)
	})
	routerGroup.GET("/club/:club_id", func(context *gin.Context) {
		handlers.GetJoinRequestsByClubHandler(context, wrapper)
	})
	routerGroup.GET("/", func(context *gin.Context) {
		handlers.GetAllJoinRequestsHandler(context, wrapper)
	})
	routerGroup.GET("/:id", func(context *gin.Context) {
		handlers.GetJoinRequestsByIdHandler(context, wrapper)
	})
}
