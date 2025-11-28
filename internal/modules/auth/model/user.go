package model

import "time"

type User = struct {
	Id        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Nickname  string    `json:"nickname,omitempty"`
	Mail      string    `json:"gmail,omitempty"`
	Password  string    `json:"password,omitempty"`
	Token     string    `json:"token,omitempty"`
	Role      string    `json:"role,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

type ResponseAuthentication = struct {
	Id           int
	HeshPassword string
	Role         string
	CreatedAt    time.Time
}
