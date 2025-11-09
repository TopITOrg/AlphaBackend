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

	if request.TotalPlaces < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "TotalPlaces must be greater than or equal to 0.",
		})
		return
	}

	if request.RequiredWorkoutPerWeek <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "RequiredWorkoutPerWeek must be greater than 0.",
		})
		return
	}

	existsSport, err := wrapper.Db.Queries.CheckSportTypeExists(ctx, request.SportTypeID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database error checking sport type"})
		return
	}
	if !existsSport {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid SportTypeID provided"})
		return
	}

	existsEducation, err := wrapper.Db.Queries.CheckEducationLevelExists(ctx, request.EducationLevelID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database error checking education level"})
		return
	}
	if !existsEducation {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid EducationLevelID provided"})
		return
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
			"message": "Unknown error",
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
