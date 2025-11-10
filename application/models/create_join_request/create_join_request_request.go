package create_join_request

type CreateJoinRequestRequest struct {
	ClubID int64 `json:"club_id"`
	UserID int64 `json:"user_id"`
}
