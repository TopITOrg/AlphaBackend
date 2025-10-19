package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/claims"
	"sport_platform/application/models/create_user"
	"sport_platform/application/models/shared"
	"sport_platform/internal/mapper"
	"sport_platform/internal/service_wrapper"
	"sport_platform/internal/sqlc/db_queries"

	"github.com/gin-gonic/gin"
)

func CreateUserHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	var request create_user.CreateUserRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("can't parse query as error happend: %s", err),
		})
		return
	}

	var createParams db_queries.CreateUserParams

	paramsMappingError := mapper.Mapper{}.Map(
		&createParams,
		request,
		struct {
			Password []byte
			Role     string
		}{
			Password: wrapper.PasswordHandler.HashPassword(request.Password),
			Role:     shared.Student,
		},
	)

	if paramsMappingError != nil {
		fmt.Println(paramsMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	user, err := wrapper.Db.Queries.CreateUser(ctx, createParams)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	var userClaims claims.UserClaims

	claimsMappingError := mapper.Mapper{}.Map(&userClaims, user)
	if claimsMappingError != nil {
		fmt.Println(claimsMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	accessToken, refreshToken, tokenGenerationError := wrapper.JwtHandler.GenerateJwtPair(userClaims, fmt.Sprintf("%d", user.ID))
	if tokenGenerationError != nil {
		fmt.Println(claimsMappingError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	var groupName string
	groupName, groupNameConstructionError := wrapper.Db.Queries.GetGroupByID(ctx, user.ID)
	if groupNameConstructionError != nil {
		fmt.Println(groupNameConstructionError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unknown error",
		})
		return
	}

	var response create_user.CreateUserResponse

	responseMappingError := mapper.Mapper{}.Map(
		&response,
		user,
		struct {
			GroupName    string
			AccessToken  string
			RefreshToken string
		}{
			GroupName:    groupName,
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
