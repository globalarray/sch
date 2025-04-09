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

	//go:embed quiz/delete.sql
	DeleteQuiz string

	//go:embed quiz/select_by_created_by.sql
	SelectQuizzesCreatedByUserID string

	//go:embed question/insert.sql
	InsertQuestion string

	//go:embed question/select.sql
	SelectQuestion string

	//go:embed question/delete.sql
	DeleteQuestion string

	//go:embed question/delete_by_quiz_id.sql
	DeleteQuestionsByQuizID string

	//go:embed question/select_by_quiz_id.sql
	SelectQuestionsByQuizID string

	//go:embed question/update_answers.sql
	UpdateQuestionAnswers string

	//go:embed quiz_progress/insert.sql
	InsertQuizProgress string

	//go:embed quiz_progress/select_user.sql
	SelectUserQuizProgress string

	//go:embed quiz_progress/delete_by_quiz_id.sql
	DeleteProgressesByQuizID string

	//go:embed quiz_progress/select_quiz.sql
	SelectAllUsersQuizProgress string

	//go:embed quiz_result/insert.sql
	InsertQuizResult string

	//go:embed quiz_result/select.sql
	SelectQuizResult string

	//go:embed quiz_result/select_by_quiz_id.sql
	SelectQuizResultByQuizID string

	//go:embed quiz_result/delete_by_quiz_id.sql
	DeleteResultsByQuizID string

	//go:embed user/select_by_surname.sql
	SelectUsersBySurname string
)
