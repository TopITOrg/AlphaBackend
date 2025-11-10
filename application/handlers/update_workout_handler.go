package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/claims"
	"sport_platform/application/models/shared"
	"sport_platform/application/models/update_workout"
	"sport_platform/internal/middleware"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func UpdateWorkoutHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	var request update_workout.UpdateWorkoutRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("can't parse query as error happend: %s", err),
		})
		return
	}

	claimsRaw, exists := ctx.Get(middleware.ClaimsKey)
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "Unauthorized",
			},
		)
		return
	}

	userClaims := claimsRaw.(claims.UserClaims)

	if userClaims.Role == shared.Student {
		ctx.JSON(
			http.StatusForbidden,
			gin.H{
				"message": "No permission",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		request,
	)
}
