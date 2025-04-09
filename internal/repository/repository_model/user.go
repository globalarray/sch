package repository_model

import "strings"

type User struct {
	TelegramID int64  `db:"tg_id"`
	Name       string `db:"name"`
	Surname    string `db:"surname"`
	Patronymic string `db:"patronymic"`
	Role       string `db:"role"`
}

func NewUser(id int64, name, surname, patronymic, role string) User {
	return User{
		TelegramID: id,
		Name:       name,
		Surname:    surname,
		Patronymic: patronymic,
		Role:       role,
	}
}

func (u User) PrettyName() string {
	if u.Patronymic == "" {
		return u.Name
	}

	return u.Name + " " + u.Patronymic
}

func (u User) FullName() string {
	fullName := strings.Builder{}

	for _, d := range []string{u.Surname, u.Name, u.Patronymic} {
		if d != "" {
			fullName.WriteString(d + " ")
		}
	}

	return strings.TrimSpace(fullName.String())
}
