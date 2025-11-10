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

func GetJoinRequestsByClubHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	club_idStr := ctx.Param("club_id")
	club_id, err := strconv.ParseInt(club_idStr, 10, 64)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid club ID",
			},
		)
		return
	}
	join_requests, dbError := wrapper.Db.Queries.GetJoinRequestsByClub(ctx, club_id)
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

	var response get_join_request.GetJoinRequestsResponse
	response.JoinRequests = make([]get_join_request.JoinRequest, 0, len(join_requests))

	for _, join_request := range join_requests {
		var join_requestResponse get_join_request.JoinRequest
		mappingError := mapper.Mapper{}.Map(&join_requestResponse, join_request)
		if mappingError != nil {
			fmt.Printf("GetJoinRequestsByClubHandler: mapping error: %s\n", mappingError)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "Mapping error",
				},
			)
			return
		}
		response.JoinRequests = append(response.JoinRequests, join_requestResponse)
	}

	ctx.JSON(
		http.StatusOK,
		response,
	)
}
