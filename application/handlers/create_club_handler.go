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
	"sport_platform/internal/sqlc/db_queries"

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

	if userClaims.Role != shared.Teacher && userClaims.Role != shared.Admin {
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

	if request.TotalPlaces <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "TotalPlaces must be greater than zero.",
		})
		return
	}

	if request.RequiredWorkoutPerWeek < 0 || request.RequiredWorkoutPerWeek > 7 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "RequiredWorkoutPerWeek must be between 0 and 7.",
		})
		return
	}

	if userClaims.Role == shared.Teacher {
		request.TeacherID = userClaims.ID
	}

	if userClaims.Role == shared.Teacher {
		request.TeacherID = userClaims.ID
	}

	var createParams db_queries.CreateClubParams

	paramsMappingError := mapper.Mapper{}.Map(&createParams, request)
	if paramsMappingError != nil {
		fmt.Println(paramsMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	club, err := wrapper.Db.Queries.CreateClub(ctx, createParams)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create club",
		})
		return
	}

	var response create_club.CreateClubResponse

	responseMappingError := mapper.Mapper{}.Map(&response, club)
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
