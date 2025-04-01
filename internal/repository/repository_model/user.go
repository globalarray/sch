package repository_model

type User struct {
	TelegramID int64
	Name       string
	Surname    string
	Patronymic string
	Role       string
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
