package role

import "benzo/internal/lang"

type Student struct{}

func (Student) Name() string {
	return "student"
}

func (Student) Translation() string {
	return lang.Student
}
