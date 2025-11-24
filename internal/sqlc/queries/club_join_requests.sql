-- name: UpdateClubJoinRequestStatus :one
UPDATE club_join_requests
SET
    status = @status,
    updated_at = now()
WHERE id = @id and is_deleted = FALSE
RETURNING club_join_requests.*;