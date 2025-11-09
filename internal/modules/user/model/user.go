package model

type User = struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Mail     string `json:"gmail,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

type HashPasswordAndId = struct {
	Id           int
	HeshPassword string
}
