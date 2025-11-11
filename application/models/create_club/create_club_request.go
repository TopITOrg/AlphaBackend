package create_club

type CreateClubRequest struct {
	Name                   string `json:"name" binding:"required"`
	Description            string `json:"description"`
	SportTypeID            int64  `json:"sport_type_id" binding:"required"`
	TeacherID              int64  `json:"teacher_id" binding:"required"`
	TotalPlaces            int    `json:"total_places"`
	Place                  string `json:"place" binding:"required"`
	EducationLevelName     string `json:"education_level_name" binding:"required"`
	RequiredWorkoutPerWeek int    `json:"required_workout_per_week"`
}
