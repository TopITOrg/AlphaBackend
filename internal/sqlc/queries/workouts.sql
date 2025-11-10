-- name: CreateWorkout :one
INSERT INTO workouts
(club_id, start_date, end_date)
VALUES
(@club_id, @start_date, @end_date)
RETURNING workouts.*;