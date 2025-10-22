-- name: GetSportClubs :many
SELECT 
    sc.id,
    sc.name,
    sc.description,
    sc.sport_type,
    sc.total_slots,
    sc.reserved_slots,
    sc.location,
    sc.min_hours_per_week,
    sc.skill_level,
    sc.created_by,
    sc.created_at,
    sc.updated_at,
    sc.is_deleted,
    u.full_name as creator_name
FROM sport_clubs sc
LEFT JOIN users u ON sc.created_by = u.id AND u.is_deleted = false
WHERE sc.is_deleted = false
    AND (sqlc.narg('sport_type')::text IS NULL OR sc.sport_type = sqlc.narg('sport_type'))
    AND (sqlc.narg('skill_level')::text IS NULL OR sc.skill_level = sqlc.narg('skill_level'))
    AND (sqlc.narg('location')::text IS NULL OR sc.location ILIKE '%' || sqlc.narg('location') || '%')
    AND (sqlc.narg('search')::text IS NULL OR sc.name ILIKE '%' || sqlc.narg('search') || '%' OR sc.description ILIKE '%' || sqlc.narg('search') || '%')
ORDER BY sc.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetSportClubById :one
SELECT 
    sc.id,
    sc.name,
    sc.description,
    sc.sport_type,
    sc.total_slots,
    sc.reserved_slots,
    sc.location,
    sc.min_hours_per_week,
    sc.skill_level,
    sc.created_by,
    sc.created_at,
    sc.updated_at,
    sc.is_deleted,
    u.full_name as creator_name
FROM sport_clubs sc
LEFT JOIN users u ON sc.created_by = u.id AND u.is_deleted = false
WHERE sc.id = sqlc.arg('id') AND sc.is_deleted = false;

-- name: GetClubSchedules :many
SELECT 
    id,
    club_id,
    day_of_week,
    start_time,
    end_time,
    location,
    created_at,
    updated_at,
    is_deleted
FROM club_schedules
WHERE club_id = sqlc.arg('club_id') AND is_deleted = false
ORDER BY 
    CASE day_of_week
        WHEN 'monday' THEN 1
        WHEN 'tuesday' THEN 2
        WHEN 'wednesday' THEN 3
        WHEN 'thursday' THEN 4
        WHEN 'friday' THEN 5
        WHEN 'saturday' THEN 6
        WHEN 'sunday' THEN 7
    END,
    start_time;

-- name: GetClubKeyMetrics :many
SELECT 
    id,
    club_id,
    name,
    unit,
    description,
    created_at,
    updated_at,
    is_deleted
FROM club_key_metrics
WHERE club_id = sqlc.arg('club_id') AND is_deleted = false;

-- name: UpdateSportClub :one
UPDATE sport_clubs 
SET 
    name = sqlc.arg('name'),
    description = sqlc.arg('description'),
    sport_type = sqlc.arg('sport_type'),
    total_slots = sqlc.arg('total_slots'),
    reserved_slots = sqlc.arg('reserved_slots'),
    location = sqlc.arg('location'),
    min_hours_per_week = sqlc.arg('min_hours_per_week'),
    skill_level = sqlc.arg('skill_level'),
    updated_at = now()
WHERE id = sqlc.arg('id') AND is_deleted = false
RETURNING *;

-- name: UpdateClubSchedule :one
UPDATE club_schedules 
SET 
    day_of_week = sqlc.arg('day_of_week'),
    start_time = sqlc.arg('start_time'),
    end_time = sqlc.arg('end_time'),
    location = sqlc.arg('location'),
    updated_at = now()
WHERE id = sqlc.arg('id') AND club_id = sqlc.arg('club_id') AND is_deleted = false
RETURNING *;

-- name: CreateClubSchedule :one
INSERT INTO club_schedules (club_id, day_of_week, start_time, end_time, location)
VALUES (sqlc.arg('club_id'), sqlc.arg('day_of_week'), sqlc.arg('start_time'), sqlc.arg('end_time'), sqlc.arg('location'))
RETURNING *;

-- name: DeleteClubSchedules :exec
UPDATE club_schedules 
SET is_deleted = true, updated_at = now()
WHERE club_id = sqlc.arg('club_id') AND is_deleted = false;

-- name: UpdateClubKeyMetric :one
UPDATE club_key_metrics 
SET 
    name = sqlc.arg('name'),
    unit = sqlc.arg('unit'),
    description = sqlc.arg('description'),
    updated_at = now()
WHERE id = sqlc.arg('id') AND club_id = sqlc.arg('club_id') AND is_deleted = false
RETURNING *;

-- name: CreateClubKeyMetric :one
INSERT INTO club_key_metrics (club_id, name, unit, description)
VALUES (sqlc.arg('club_id'), sqlc.arg('name'), sqlc.arg('unit'), sqlc.arg('description'))
RETURNING *;

-- name: DeleteClubKeyMetrics :exec
UPDATE club_key_metrics 
SET is_deleted = true, updated_at = now()
WHERE club_id = sqlc.arg('club_id') AND is_deleted = false;
