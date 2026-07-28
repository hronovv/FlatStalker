package models

type Listing struct {
	ID     int64
	UserID int64
	URL    string
	Paused bool
}
