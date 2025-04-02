package repository_query

import _ "embed"

var (
	//go:embed secret/select.sql
	SelectSecretInformation string

	//go:embed secret/update_personal_data.sql
	UpdateSecretPersonalData string

	//go:embed secret/update_role.sql
	UpdateSecretRole string

	//go:embed secret/insert.sql
	InsertSecret string

	//go:embed secret/delete.sql
	DeleteSecret string

	//go:embed user/select.sql
	SelectUser string

	//go:embed user/insert.sql
	InsertUser string

	//go:embed quiz/insert.sql
	InsertQuiz string

	//go:embed quiz/select.sql
	SelectQuiz string

	//go:embed question/insert.sql
	InsertQuestion string

	//go:embed question/select_by_quiz_id.sql
	SelectQuestionsByQuizID string

	//go:embed question/update_answers.sql
	UpdateQuestionAnswers string
)
