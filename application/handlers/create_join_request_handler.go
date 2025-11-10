package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/create_join_request"
	"sport_platform/application/models/shared"
	"sport_platform/internal/mapper"
	"sport_platform/internal/service_wrapper"
	"sport_platform/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

func CreateJoinRequestHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	var request create_join_request.CreateJoinRequestRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("can't parse query as error happend: %s", err),
		})
		return
	}

	var createParams db_queries.CreateJoinRequestParams

	paramsMappingError := mapper.Mapper{}.Map(
		&createParams,
		request,
		struct {
			Status string
		}{
			Status: shared.NotAccepted,
		},
	)

	if paramsMappingError != nil {
		fmt.Println(paramsMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}
	join_request, err := wrapper.Db.Queries.CreateJoinRequest(ctx, createParams)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	var response create_join_request.CreateJoinRequestResponse

	responseMappingError := mapper.Mapper{}.Map(
		&response,
		join_request,
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
