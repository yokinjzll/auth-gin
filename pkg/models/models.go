package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type User struct {
	ID       int
	Username string
	Password string
}

type UserDetails struct {
	UserID     int
	First_name string
	Last_name  string
	Dob        string
	Created_at time.Time
}
