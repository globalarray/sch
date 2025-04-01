package role

import "benzo/internal/lang"

type Admin struct{}

func (Admin) Name() string {
	return "admin"
}

func (Admin) Translation() string {
	return lang.Admin
}
