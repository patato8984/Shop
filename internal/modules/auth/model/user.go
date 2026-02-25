package model

import (
	"errors"
	"time"
)

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

var (
	ErrCheckPassword           = errors.New("incorrect password")
	ErrShortPasswordOrNickname = errors.New("short password or nickname")
	ErrUserNotFound            = errors.New("user not found")
	ErrNickNameBusy            = errors.New("nickname busy")
	ErrMailBusy                = errors.New("mail busy")
	ErrJson                    = errors.New("error json")
)
