package controllers

import (
	"sport_platform/application/handlers"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func ClubJoinRequestController(engine *gin.Engine, wrapper *service_wrapper.Wrapper) {
	routerGroup := engine.Group("/club_join_requests")

	routerGroup.PUT(("/update"), func(context *gin.Context) {
		handlers.UpdateClubJoinRequestStatusHandler(context, wrapper)
	})
}
