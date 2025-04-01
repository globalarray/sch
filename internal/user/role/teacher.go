package role

import "benzo/internal/lang"

type Teacher struct{}

func (Teacher) Name() string {
	return "teacher"
}

func (Teacher) Translation() string {
	return lang.Teacher
}
