package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/get_join_request"
	"sport_platform/internal/mapper"
	"sport_platform/internal/service_wrapper"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetJoinRequestsByIdHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	club_idStr := ctx.Param("id")
	id, err := strconv.ParseInt(club_idStr, 10, 64)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid ID",
			},
		)
		return
	}
	join_request, dbError := wrapper.Db.Queries.GetJoinRequestById(ctx, id)
	if dbError != nil {
		fmt.Printf("Error while getting join requests: %s\n", dbError)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Something unusual happened",
			},
		)
		return
	}

	var response get_join_request.GetJoinRequestResponse
	mappingError := mapper.Mapper{}.Map(&response, join_request)
	if mappingError != nil {
		fmt.Printf("GetJoinRequestByIdHandler: mapping error: %s\n", mappingError)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Mapping error",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		response,
	)
}
