package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/update_club"
	"sport_platform/internal/mapper"
	"sport_platform/internal/service_wrapper"
	"sport_platform/internal/sqlc/db_queries"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdateClubHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Printf("UpdateClubHandler: Invalid club ID: %s\n", idStr)
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Invalid club ID",
			},
		)
		return
	}

	fmt.Printf("UpdateClubHandler: Updating club ID %d\n", id)

	var request update_club.UpdateClubRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		fmt.Printf("UpdateClubHandler: JSON bind error: %s\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	fmt.Printf("UpdateClubHandler: Request data - Name: %s, Description: %s\n", request.Name, request.Description)

	// Создаем структуру параметров для обновления
	updateParams := db_queries.UpdateClubParams{
		ID:                     id, // Устанавливаем ID из URL параметра
		Name:                   request.Name,
		Description:            request.Description,
		SportTypeID:            request.SportTypeID,
		TeacherID:              request.TeacherID,
		TotalPlaces:            request.TotalPlaces,
		Place:                  request.Place,
		EducationLevelID:       request.EducationLevelID,
		RequiredWorkoutPerWeek: request.RequiredWorkoutPerWeek,
	}

	fmt.Printf("UpdateClubHandler: Update params - ID: %d, Name: %s\n", updateParams.ID, updateParams.Name)

	// Обновляем клуб в БД
	updatedClub, err := wrapper.Db.Queries.UpdateClub(ctx, updateParams)
	if err != nil {
		fmt.Printf("UpdateClubHandler: Database update error: %s\n", err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Something unusual happened",
			},
		)
		return
	}

	fmt.Printf("UpdateClubHandler: Successfully updated club ID %d\n", id)

	// Получаем обновленный клуб с JOIN'ами для полной информации
	fullClubInfo, err := wrapper.Db.Queries.GetClubById(ctx, id)
	if err != nil {
		fmt.Printf("UpdateClubHandler: Error getting updated club with joins: %s\n", err)
		// Если не удалось получить с JOIN'ами, вернем хотя бы базовые данные
		var response update_club.UpdateClubResponse
		responseMappingError := mapper.Mapper{}.Map(&response, updatedClub)
		if responseMappingError != nil {
			fmt.Printf("UpdateClubHandler: Response mapping error: %s\n", responseMappingError)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "Something unusual happened",
				},
			)
			return
		}
		fmt.Printf("UpdateClubHandler: Success - returning basic club info\n")
		ctx.JSON(
			http.StatusOK,
			response,
		)
		return
	}

	// Преобразуем обновленный клуб с JOIN'ами в ответ
	var response update_club.UpdateClubResponse
	responseMappingError := mapper.Mapper{}.Map(&response, fullClubInfo)
	if responseMappingError != nil {
		fmt.Printf("UpdateClubHandler: Response mapping error: %s\n", responseMappingError)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Something unusual happened",
			},
		)
		return
	}

	fmt.Printf("UpdateClubHandler: Success - returning updated club with full info\n")
	ctx.JSON(
		http.StatusOK,
		response,
	)
}
