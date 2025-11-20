package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/get_clubs"
	"sport_platform/internal/mapper"
	"sport_platform/internal/service_wrapper"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetClubByIdHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid club ID",
			},
		)
		return
	}

	club, err := wrapper.Db.Queries.GetClubById(ctx, id)
	if err != nil {
		fmt.Printf("Error while getting club: %s\n", err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Something unusual happened",
			},
		)
		return
	}

	var response get_clubs.Club
	mappingError := mapper.Mapper{}.Map(&response, club)
	if mappingError != nil {
		fmt.Printf("Error while mapping club: %s\n", mappingError)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Something unusual happened",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		response,
	)
}
