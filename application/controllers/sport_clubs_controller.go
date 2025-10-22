package controllers

import (
    "sport_platform/application/handlers"
    "sport_platform/internal/service_wrapper"

    "github.com/gin-gonic/gin"
)

func SportClubsController(engine *gin.Engine, wrapper *service_wrapper.Wrapper) {
    routerGroup := engine.Group("/sport-clubs")
    
    // Публичные endpoints - получение секций
    routerGroup.GET("", func(context *gin.Context) {
        handlers.GetSportClubsHandler(context, wrapper)
    })
    routerGroup.GET("/:id", func(context *gin.Context) {
        handlers.GetSportClubHandler(context, wrapper)
    })

    // Защищенные endpoints - обновление секций
    routerGroup.PUT("/:id", func(context *gin.Context) {
        handlers.UpdateSportClubHandler(context, wrapper)
    })
}
