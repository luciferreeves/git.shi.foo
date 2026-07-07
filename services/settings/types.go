package settings

import "time"

type TokenView struct {
	ID         uint
	Label      string
	Preview    string
	CreatedAt  time.Time
	LastUsedAt *time.Time
}

type IndexContext struct {
	Title    string
	Tokens   []TokenView
	NewToken string
}
