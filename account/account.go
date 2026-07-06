package account

import "time"

type Response struct {
	ID         uint      `json:"id"`
	ProviderID string    `json:"provider_id"`
	Login      string    `json:"login"`
	Email      string    `json:"email"`
	Avatar     string    `json:"avatar"`
	CreatedAt  time.Time `json:"created_at"`
}
