package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/get_clubs"
	"sport_platform/internal/mapper"
	"sport_platform/internal/middleware"
	"sport_platform/internal/service_wrapper"
	"sport_platform/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

func GetClubsHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	var request get_clubs.GetClubRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("can't parse query as error happend: %s", err),
		})
		return
	}

	_, exists := ctx.Get(middleware.ClaimsKey)
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "Unauthorized",
			},
		)
		return
	}

	clubs, dbError := wrapper.Db.Queries.GetAllClubs(ctx)
	if dbError != nil {
		fmt.Printf("Error while getting clubs: %s\n", dbError)

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Something unusual happened",
			},
		)
		return
	}

	var response get_clubs.GetAllClubsResponse

	response.Clubs = make([]get_clubs.GetClubResponse, len(clubs))

	type SourceWrapper struct {
		Items []db_queries.GetAllClubsRow
	}

	type DestWrapper struct {
		Items []get_clubs.GetClubResponse
	}

	sourceWrapper := SourceWrapper{Items: clubs}
	destWrapper := DestWrapper{Items: response.Clubs}

	mappingError := mapper.Mapper{}.Map(&destWrapper, sourceWrapper)
	if mappingError != nil {
		fmt.Printf("Clubs mapping error: %s\n", mappingError)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Unknown error",
			},
		)
		return
	}
	response.Clubs = destWrapper.Items

	ctx.JSON(
		http.StatusOK,
		response,
	)
}
