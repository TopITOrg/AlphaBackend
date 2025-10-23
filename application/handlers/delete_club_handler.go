package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/claims"
	"sport_platform/application/models/delete_club"
	"sport_platform/application/models/shared"
	"sport_platform/internal/middleware"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func DeleteClubHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
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

	if userClaims.Role != shared.Teacher && userClaims.Role != string(shared.Admin) {
		ctx.JSON(
			http.StatusForbidden,
			gin.H{
				"message": "No permission",
			},
		)
		return
	}

	clubID, err := parseInt64Param(ctx, "id")
	if err != nil || clubID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid club ID",
		})
		return
	}

	accessToken, refreshToken, tokenGenerationError := wrapper.JwtHandler.GenerateJwtPair(userClaims, fmt.Sprint(userClaims.ID))
	if tokenGenerationError != nil {
		fmt.Println(tokenGenerationError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	response := delete_club.DeleteClubResponse{
		IsDeleted:    true,
		Message:      "Club deleted",
		ClubID:       clubID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	ctx.JSON(
		http.StatusOK,
		response,
	)
}

func parseInt64Param(ctx *gin.Context, param string) (int64, error) {
	idStr := ctx.Param(param)
	var id int64
	_, err := fmt.Sscanf(idStr, "%d", &id)
	return id, err
}
