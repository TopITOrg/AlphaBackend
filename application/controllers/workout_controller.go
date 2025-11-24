package controllers

import (
	"sport_platform/application/handlers"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func WorkoutController(engine *gin.Engine, wrapper *service_wrapper.Wrapper) {
	routerGroup := engine.Group("/workouts")
	routerGroup.POST("/create", func(context *gin.Context) {
		handlers.CreateWorkoutHandler(context, wrapper)
	routerGroup.PUT("/update", func(context *gin.Context) {
		handlers.UpdateWorkoutHandler(context, wrapper)
	})
}
