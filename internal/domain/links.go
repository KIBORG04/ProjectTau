package domain

type Link struct {
	ID   uint
	Link string
	Date string
}

type RoundInDate struct {
	Round      string `json:"name"`
	IDontKnown string `json:"type"`
	CreatedAt  string `json:"mtime"`
}
