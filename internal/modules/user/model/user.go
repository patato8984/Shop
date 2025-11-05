package model

type User = struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Mail     string `json:"gmail"`
	Password string `json:"password"`
}
