package repository_model

import "time"

type Quiz struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedBy int64     `db:"created_by"`
	Creation  time.Time `db:"creation"`
}

func NewQuiz(name string, createdBy int64) Quiz {
	return Quiz{
		Name:      name,
		CreatedBy: createdBy,
	}
}
