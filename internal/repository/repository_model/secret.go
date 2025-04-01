package repository_model

import (
	"benzo/internal/user/role"
	"time"
)

type Secret struct {
	Key        string
	Role       string
	Name       string
	Patronymic string
	Surname    string
	Expiration time.Time
	Creation   time.Time
	CreatedBy  int64
}

func NewSecretDefault(key string) Secret {
	return Secret{
		Key:        key,
		Role:       role.Student{}.Name(),
		Name:       "nil",
		Patronymic: "nil",
		Surname:    "nil",
		Expiration: time.Time{},
		Creation:   time.Time{},
		CreatedBy:  0,
	}
}
