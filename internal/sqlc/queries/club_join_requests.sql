-- name: GetJoinRequests :many
SELECT *
FROM club_join_requests
WHERE 
    (@id::bigint IS NULL OR id = @id::bigint)
    AND (@club_id::bigint IS NULL OR club_id = @club_id::bigint)
    AND (@user_id::bigint IS NULL OR user_id = @user_id::bigint)
    AND (is_deleted = false)
ORDER BY id
LIMIT CASE WHEN @limit_::bigint IS NOT NULL THEN @limit_::bigint END
OFFSET CASE WHEN @offset_::bigint IS NOT NULL THEN @offset_::bigint ELSE 0 END;

-- name: CreateJoinRequest :one
INSERT INTO club_join_requests 
(club_id, user_id, status)
VALUES 
(@club_id, @user_id, @status)
RETURNING *;