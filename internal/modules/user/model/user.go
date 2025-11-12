package model

type User = struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Mail     string `json:"gmail,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
	Role     string `json:"role,omitempty"`
}

type ResponseAuthentication = struct {
	Id           int
	HeshPassword string
	Role         string
}
