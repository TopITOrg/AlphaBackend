-- name: CreateClub :one
INSERT INTO clubs (
    name,
    description,
    sport_type_id,
    teacher_id,
    total_places,
    place,
    education_level_id,
    required_workout_per_week
) VALUES (
             sqlc.arg(name),
             sqlc.arg(description),
             sqlc.arg(sport_type_id),
             sqlc.arg(teacher_id),
             sqlc.arg(total_places),
             sqlc.arg(place),
             (SELECT id FROM education_levels WHERE education_levels.name = sqlc.arg(education_level_name)),
             sqlc.arg(required_workout_per_week)
         )
    RETURNING *;

-- name: GetClubByID :one
SELECT * FROM clubs WHERE id = $1 AND is_deleted = false;

-- name: SoftDeleteClub :exec
UPDATE clubs
SET is_deleted = true, updated_at = NOW()
WHERE id = $1;

-- name: GetClubsByTeacher :many
SELECT * FROM clubs
WHERE teacher_id = $1 AND is_deleted = false
ORDER BY created_at DESC;

-- name: ListActiveClubs :many
SELECT * FROM clubs
WHERE is_deleted = false
ORDER BY name;

-- name: CheckClubOwnership :one
SELECT EXISTS(
    SELECT 1 FROM clubs
    WHERE id = $1 AND teacher_id = $2 AND is_deleted = false
);
-- name: CheckClubExists :one
SELECT EXISTS(
    SELECT 1 FROM clubs
    WHERE id = $1 AND is_deleted = FALSE
);

-- name: CheckSportTypeExists :one
SELECT EXISTS(
    SELECT 1 FROM sport_types
    WHERE id = $1
);

-- name: CheckEducationLevelExistsByName :one
SELECT EXISTS(
    SELECT 1 FROM education_levels
    WHERE name = $1
);