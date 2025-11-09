package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/get_clubs"
	"sport_platform/internal/mapper"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func GetClubsHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	fmt.Println("GetClubsHandler: starting...")

	clubs, err := wrapper.Db.Queries.GetAllClubs(ctx)
	if err != nil {
		fmt.Printf("GetClubsHandler: database error: %s\n", err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Database error",
			},
		)
		return
	}

	fmt.Printf("GetClubsHandler: found %d clubs\n", len(clubs))

	var response get_clubs.GetClubsResponse
	response.Clubs = make([]get_clubs.Club, 0, len(clubs))

	for _, club := range clubs {
		var clubResponse get_clubs.Club
		mappingError := mapper.Mapper{}.Map(&clubResponse, club)
		if mappingError != nil {
			fmt.Printf("GetClubsHandler: mapping error: %s\n", mappingError)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "Mapping error",
				},
			)
			return
		}
		response.Clubs = append(response.Clubs, clubResponse)
	}

	fmt.Println("GetClubsHandler: success")
	ctx.JSON(
		http.StatusOK,
		response,
	)
}
