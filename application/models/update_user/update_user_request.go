package update_user

import "time"

type UpdateUserRequest struct {
	FullName          *string   `json:"full_name"`
	SocialNetworkLink *string   `json:"social_network_link"`
	PhoneNumber       *string   `json:"phone_number"`
	Email             *string   `json:"email"`
	BirthDate         time.Time `json:"birth_date"`
	Password          *string   `json:"password" mapper:"exclude"`
	GroupID           *int64    `json:"group_id"`
}
