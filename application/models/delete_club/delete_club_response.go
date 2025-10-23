package delete_club

type DeleteClubResponse struct {
	ClubID       int64  `json:"club_id"`
	IsDeleted    bool   `json:"is_deleted"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
