-- name: UpdateWorkout :one
UPDATE workouts
SET
    cancelled = COALESCE(sqlc.narg(cancelled), cancelled),
    start_date = COALESCE(sqlc.narg(start_date), start_date),
    end_date = COALESCE(sqlc.narg(end_date), end_date),
    updated_at = now()
WHERE id = @id
RETURNING workouts.*;