package role

import (
	"errors"
	"slices"
)

type Role interface {
	Name() string
	Translation() string
}

var (
	ErrRoleNotFound = errors.New("role not found")
)

var (
	roles = []string{Student{}.Name(), Teacher{}.Name(), Admin{}.Name()}
)

func RightsLevel(r Role) int {
	return slices.Index(roles, r.Name())
}

func FromName(name string) (r Role, err error) {
	switch name {
	case "admin":
		return Admin{}, nil
	case "teacher":
		return Teacher{}, nil
	case "student":
		return Student{}, nil
	}

	return r, ErrRoleNotFound
}
