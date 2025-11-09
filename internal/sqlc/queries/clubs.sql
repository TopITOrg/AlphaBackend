-- name: GetClubById :one
SELECT
    clubs.id,
    clubs.name,
    clubs.description,
    clubs.sport_type_id,
    clubs.teacher_id,
    clubs.total_places,
    clubs.place,
    clubs.education_level_id,
    clubs.required_workout_per_week,
    clubs.created_at,
    clubs.updated_at,
    clubs.is_deleted,
    sport_types.name as sport_type_name,
    education_levels.name as education_level_name,
    users.full_name as teacher_name
FROM clubs
         LEFT JOIN sport_types ON clubs.sport_type_id = sport_types.id AND sport_types.is_deleted = false
         LEFT JOIN education_levels ON clubs.education_level_id = education_levels.id AND education_levels.is_deleted = false
         LEFT JOIN users ON clubs.teacher_id = users.id AND users.is_deleted = false
WHERE clubs.id = $1 AND clubs.is_deleted = false;

-- name: GetAllClubs :many
SELECT
    clubs.id,
    clubs.name,
    clubs.description,
    clubs.sport_type_id,
    clubs.teacher_id,
    clubs.total_places,
    clubs.place,
    clubs.education_level_id,
    clubs.required_workout_per_week,
    clubs.created_at,
    clubs.updated_at,
    clubs.is_deleted,
    sport_types.name as sport_type_name,
    education_levels.name as education_level_name,
    users.full_name as teacher_name
FROM clubs
         LEFT JOIN sport_types ON clubs.sport_type_id = sport_types.id AND sport_types.is_deleted = false
         LEFT JOIN education_levels ON clubs.education_level_id = education_levels.id AND education_levels.is_deleted = false
         LEFT JOIN users ON clubs.teacher_id = users.id AND users.is_deleted = false
WHERE clubs.is_deleted = false;

-- name: UpdateClub :one
UPDATE clubs
SET
    name = $2,
    description = $3,
    sport_type_id = $4,
    teacher_id = $5,
    total_places = $6,
    place = $7,
    education_level_id = $8,
    required_workout_per_week = $9,
    updated_at = now()
WHERE id = $1 AND is_deleted = false
RETURNING id, name, description, sport_type_id, teacher_id, total_places, place, education_level_id, required_workout_per_week, created_at, updated_at, is_deleted;