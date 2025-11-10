-- name: GetWorkoutsByClub :many
SELECT id, club_id, start_date, end_date, cancelled, created_at, updated_at
FROM workouts
WHERE club_id = $1 AND is_deleted = FALSE AND cancelled = FALSE
ORDER BY start_date DESC;

-- name: CheckWorkoutExists :one
SELECT EXISTS(
    SELECT 1 FROM workouts
    WHERE id = $1 AND is_deleted = FALSE
);

-- name: SoftDeleteWorkout :exec
UPDATE workouts
SET is_deleted = TRUE, updated_at = NOW()
WHERE id = $1;
