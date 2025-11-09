package update_club

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateClubResponse struct {
	ID                     int64            `json:"id"`
	Name                   string           `json:"name"`
	Description            string           `json:"description"`
	SportTypeID            int64            `json:"sport_type_id"`
	SportTypeName          *string          `json:"sport_type_name,omitempty"`
	TeacherID              int64            `json:"teacher_id"`
	TeacherName            *string          `json:"teacher_name,omitempty"`
	TotalPlaces            *int32           `json:"total_places,omitempty"`
	Place                  string           `json:"place"`
	EducationLevelID       int64            `json:"education_level_id"`
	EducationLevelName     *string          `json:"education_level_name,omitempty"`
	RequiredWorkoutPerWeek int32            `json:"required_workout_per_week"`
	CreatedAt              pgtype.Timestamp `json:"created_at"`
	UpdatedAt              pgtype.Timestamp `json:"updated_at"`
}
