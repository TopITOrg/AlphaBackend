package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/claims"
	"sport_platform/application/models/create_club"
	"sport_platform/application/models/shared"
	"sport_platform/internal/mapper"
	"sport_platform/internal/middleware"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func CreateClubHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
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

	var request create_club.CreateClubRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("can't parse query as error happend: %s", err),
		})
		return
	}

	if userClaims.Role == shared.Teacher {
		request.TeacherID = userClaims.ID
	}

	accessToken, refreshToken, tokenGenerationError := wrapper.JwtHandler.GenerateJwtPair(userClaims, fmt.Sprint(userClaims.ID))
	if tokenGenerationError != nil {
		fmt.Println(tokenGenerationError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	var response create_club.CreateClubResponse

	responseMappingError := mapper.Mapper{}.Map(
		&response,
		request,
		struct {
			AccessToken  string
			RefreshToken string
		}{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	)
	if responseMappingError != nil {
		fmt.Println(responseMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	ctx.JSON(
		http.StatusCreated,
		response,
	)
}
