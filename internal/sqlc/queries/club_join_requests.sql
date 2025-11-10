-- name: GetJoinRequestsByClub :many
SELECT *
FROM club_join_requests
WHERE
    club_id = @club_id and
    is_deleted = false;

-- name: GetJoinRequestById :one
SELECT *
FROM club_join_requests
WHERE
    id = @id and
    is_deleted = false
LIMIT 1;

-- name: GetAllJoinRequests :many
SELECT *
FROM club_join_requests;

-- name: CreateJoinRequest :one
INSERT INTO club_join_requests 
(club_id, user_id, status, created_at, updated_at)
VALUES 
(@club_id, @user_id, @status, @created_at, @updated_at)
RETURNING *;