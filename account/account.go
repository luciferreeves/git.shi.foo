package account

import "time"

type Response struct {
	ID         uint      `json:"id"`
	ProviderID string    `json:"provider_id"`
	Login      string    `json:"login"`
	Email      string    `json:"email"`
	Avatar     string    `json:"avatar"`
	Admin      bool      `json:"admin"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
}
