package get_workouts

type GetWorkoutsRequest struct {
	ClubID int64 `uri:"club_id" binding:"required"`
}
