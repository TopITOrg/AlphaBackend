package update_club

type UpdateClubRequest struct {
	Name                   string `json:"name"`
	Description            string `json:"description"`
	SportTypeID            int64  `json:"sport_type_id"`
	TeacherID              int64  `json:"teacher_id"`
	TotalPlaces            *int32 `json:"total_places,omitempty"`
	Place                  string `json:"place"`
	EducationLevelID       int64  `json:"education_level_id"`
	RequiredWorkoutPerWeek int32  `json:"required_workout_per_week"`
}
