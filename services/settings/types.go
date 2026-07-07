package settings

import "time"

type TokenView struct {
	ID         uint
	Label      string
	Preview    string
	CreatedAt  time.Time
	LastUsedAt *time.Time
}

type KeyView struct {
	ID          uint
	Title       string
	KeyType     string
	Fingerprint string
	Source      string
	CreatedAt   time.Time
}

type IndexContext struct {
	Title    string
	Tokens   []TokenView
	Keys     []KeyView
	NewToken string
}
