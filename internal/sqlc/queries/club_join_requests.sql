-- name: GetJoinRequests :many
SELECT cjr.*
FROM club_join_requests cjr
INNER JOIN clubs c ON cjr.club_id = c.id
WHERE 
    (@id::bigint IS NULL OR cjr.id = @id::bigint)
    AND (@club_id::bigint IS NULL OR cjr.club_id = @club_id::bigint)
    AND (@user_id::bigint IS NULL OR cjr.user_id = @user_id::bigint)
    AND cjr.is_deleted = false
    AND (
        cjr.user_id = @current_user_id::bigint 
        OR c.teacher_id = @current_user_id::bigint
    )
ORDER BY cjr.created_at DESC
LIMIT @limit_::bigint
OFFSET @offset_::bigint;

-- name: CreateJoinRequest :one
INSERT INTO club_join_requests 
(club_id, user_id, status)
VALUES 
(@club_id, @user_id, @status)
RETURNING *;

-- name: UpdateClubJoinRequestStatus :one
UPDATE club_join_requests
SET
    status = @status,
    updated_at = now()
WHERE id = @id and is_deleted = FALSE
RETURNING club_join_requests.*;
