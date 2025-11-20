package controllers

import (
	"sport_platform/application/handlers"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func ClubController(engine *gin.Engine, wrapper *service_wrapper.Wrapper) {
	routerGroup := engine.Group("/clubs")
	routerGroup.GET("/", func(context *gin.Context) {
		handlers.GetClubsHandler(context, wrapper)
	})
	routerGroup.GET("/:id", func(context *gin.Context) {
		handlers.GetClubByIdHandler(context, wrapper)
	})
	routerGroup.PUT("/:id", func(context *gin.Context) {
		handlers.UpdateClubHandler(context, wrapper)
	})
}
