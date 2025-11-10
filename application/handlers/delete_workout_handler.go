package handlers

import (
	"fmt"
	"net/http"
	"sport_platform/application/models/delete_workout"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func DeleteWorkoutHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
	var request delete_workout.DeleteWorkoutRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid workout ID format in URL"})
		return
	}

	var err error
	exists, err := wrapper.Db.Queries.CheckWorkoutExists(ctx, request.WorkoutID)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database error while checking workout existence"})
		return
	}

	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Workout not found or already deleted"})
		return
	}

	err = wrapper.Db.Queries.SoftDeleteWorkout(ctx, request.WorkoutID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"workout_id": request.WorkoutID,
		"message":    "Workout deleted successfully",
	})
}
