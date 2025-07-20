package models

type Location struct {
	UserID    string  `json:"user_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

type NearbyFriend struct {
	UserID   string  `json:"user_id"`
	Name     string  `json:"name"`
	Distance float64 `json:"distance"`
}
