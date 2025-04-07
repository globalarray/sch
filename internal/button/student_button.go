package button

import "benzo/internal/user/role"

type StudentButton struct {
}

var (
	studentRightsLevel = role.RightsLevel(role.Student{})
)

func (StudentButton) NeedRightsLevel() int {
	return studentRightsLevel
}
