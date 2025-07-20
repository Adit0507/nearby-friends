package models

type User struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Friends []string `json:"friends"`
}

type FriendRequest struct {
	FriendID string `json:"friend_id"`
}
