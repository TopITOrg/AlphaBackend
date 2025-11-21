package get_join_request

type GetJoinRequestRequest struct {
	ID     *int64 `json:"id"`
	ClubID *int64 `json:"club_id"`
	UserID *int64 `json:"user_id"`
	Limit  *int64 `json:"limit_"`
	Offset *int64 `json:"offset_"`
}
