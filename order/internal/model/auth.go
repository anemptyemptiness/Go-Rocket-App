package model

import "time"

type UserInfo struct {
	Login string
}

type User struct {
	UUID      string
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Session struct {
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}
