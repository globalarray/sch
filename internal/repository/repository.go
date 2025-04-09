package repository

import (
	"benzo/internal/repository/repository_model"
	"benzo/internal/repository/repository_query"
	"benzo/pkg/datasource"
	"benzo/pkg/utils"
	"github.com/jmoiron/sqlx"
	"strings"
	"sync"
)

type Repository struct {
	db *sqlx.DB
}

var (
	repoMu sync.Mutex
	repo   *Repository
)

func Repo() *Repository {
	repoMu.Lock()
	defer repoMu.Unlock()
	return repo
}

func New(login, pwd, addr, db string, port uint16) (*Repository, error) {
	dat, err := datasource.NewDatabase(login, pwd, addr, db, port)

	if err != nil {
		return nil, err
	}

	repo = &Repository{
		db: dat,
	}

	return repo, nil
}

func (repo *Repository) GetSecretByKey(key string) (sec repository_model.Secret, err error) {
	query := repository_query.SelectSecretInformation

	row := func(idx int) utils.Array {
		return utils.Array{
			&sec.Key,
			&sec.Name,
			&sec.Patronymic,
			&sec.Surname,
			&sec.Creation,
			&sec.Expiration,
			&sec.Role,
			&sec.CreatedBy,
		}
	}

	err = new(datasource.DataSource).QuerySQL(repo.db.Queryx(query, key)).Scan(row)

	return sec, err
}

func (repo *Repository) GetUserByTelegramID(id int64) (u repository_model.User, err error) {
	query := repository_query.SelectUser

	row := func(idx int) utils.Array {
		return utils.Array{
			&u.TelegramID,
			&u.Name,
			&u.Patronymic,
			&u.Surname,
			&u.Role,
		}
	}

	err = new(datasource.DataSource).QuerySQL(repo.db.Queryx(query, id)).Scan(row)

	return u, err
}

func (repo *Repository) SaveNewUser(u repository_model.User) error {
	query := repository_query.InsertUser

	args := utils.Array{
		u.TelegramID,
		u.Name,
		u.Patronymic,
		u.Surname,
		u.Role,
	}

	return new(datasource.DataSource).ExecSQL(repo.db.Exec(query, args...)).Scan(nil, nil)
}

func (repo *Repository) SaveNewSecret(secret repository_model.Secret) error {
	query := repository_query.InsertSecret

	args := utils.Array{
		secret.Key,
		secret.Name,
		secret.Patronymic,
		secret.Surname,
		secret.Role,
		secret.CreatedBy,
	}

	return new(datasource.DataSource).ExecSQL(repo.db.Exec(query, args...)).Scan(nil, nil)
}

func (repo *Repository) UpdateSecretPersonalData(key, surname, name, patronymic string) error {
	query := repository_query.UpdateSecretPersonalData

	args := utils.Array{
		surname,
		name,
		patronymic,
		key,
	}

	return new(datasource.DataSource).ExecSQL(repo.db.Exec(query, args...)).Scan(nil, nil)
}

func (repo *Repository) UpdateSecretRole(key, roleName string) error {
	query := repository_query.UpdateSecretRole

	args := utils.Array{
		roleName,
		key,
	}

	return new(datasource.DataSource).ExecSQL(repo.db.Exec(query, args...)).Scan(nil, nil)
}

func (repo *Repository) RemoveSecretByKey(key string) error {
	query := repository_query.DeleteSecret

	return new(datasource.DataSource).ExecSQL(repo.db.Exec(query, key)).Scan(nil, nil)
}

func (repo *Repository) SaveNewQuiz(quiz repository_model.Quiz) (id int64, err error) {
	query := repository_query.InsertQuiz

	args := utils.Array{
		quiz.Name,
		quiz.CreatedBy,
	}

	err = new(datasource.DataSource).ExecSQL(repo.db.Exec(query, args...)).Scan(nil, &id)

	return id, err
}

func (repo *Repository) SaveNewQuestion(question repository_model.Question) (id int64, err error) {
	query := repository_query.InsertQuestion

	args := utils.Array{
		question.QuizID,
		question.Question,
		question.Answers,
	}

	err = new(datasource.DataSource).ExecSQL(repo.db.Exec(query, args...)).Scan(nil, &id)

	return id, err
}

func (repo *Repository) GetQuizByID(id int64) (q repository_model.Quiz, err error) {
	query := repository_query.SelectQuiz

	row := func(idx int) utils.Array {
		return utils.Array{
			&q.ID,
			&q.Name,
			&q.Creation,
			&q.CreatedBy,
		}
	}

	err = new(datasource.DataSource).QuerySQL(repo.db.Queryx(query, id)).Scan(row)

	return q, err
}

func (repo *Repository) GetQuestionByID(id int64) (q repository_model.Question, err error) {
	query := repository_query.SelectQuestion

	row := func(idx int) utils.Array {
		return utils.Array{
			&q.ID,
			&q.QuizID,
			&q.Question,
			&q.Answers,
		}
	}

	err = new(datasource.DataSource).QuerySQL(repo.db.Queryx(query, id)).Scan(row)

	return q, err
}

func (repo *Repository) RemoveQuestionByID(id int64) error {
	query := repository_query.DeleteQuestion

	return new(datasource.DataSource).ExecSQL(repo.db.Exec(query, id)).Scan(nil, nil)
}

func (repo *Repository) GetQuestionsByQuizID(quizID int64) (questions []repository_model.Question, err error) {
	query := repository_query.SelectQuestionsByQuizID

	if err := repo.db.Select(&questions, query, quizID); err != nil {
		return questions, err
	}

	return questions, nil
}

func (repo *Repository) UpdateQuestionAnswers(questionID int64, answers []string) error {
	query := repository_query.UpdateQuestionAnswers

	args := utils.Array{
		strings.Join(answers, ";"),
		questionID,
	}

	return new(datasource.DataSource).ExecSQL(repo.db.Exec(query, args...)).Scan(nil, nil)
}

func (repo *Repository) SelectUsersQuizProgressByQuizID(quizID int64) (progresses []repository_model.QuizProgress, err error) {
	query := repository_query.SelectAllUsersQuizProgress

	if err := repo.db.Select(&progresses, query, quizID); err != nil {
		return progresses, err
	}

	return progresses, nil
}

func (repo *Repository) GetQuizResultsByQuizID(quizID int64) (results []repository_model.QuizResult, err error) {
	query := repository_query.SelectQuizResultByQuizID

	if err := repo.db.Select(&results, query, quizID); err != nil {
		return results, err
	}

	return results, nil
}

func (repo *Repository) GetQuizzesCreatedByUserID(userID int64) (quizzes []repository_model.Quiz, err error) {
	query := repository_query.SelectQuizzesCreatedByUserID

	if err := repo.db.Select(&quizzes, query, userID); err != nil {
		return quizzes, err
	}

	return quizzes, nil
}

func (repo *Repository) RemoveQuizByID(quizID int64) error {
	query := repository_query.DeleteQuiz

	return new(datasource.DataSource).ExecSQL(repo.db.Exec(query, quizID)).Scan(nil, nil)
}

func (repo *Repository) RemoveQuestionsByQuizID(quizID int64) error {
	query := repository_query.DeleteQuestionsByQuizID

	return new(datasource.DataSource).ExecSQL(repo.db.Exec(query, quizID)).Scan(nil, nil)
}

func (repo *Repository) RemoveProgressesByQuizID(quizID int64) error {
	query := repository_query.DeleteProgressesByQuizID

	return new(datasource.DataSource).ExecSQL(repo.db.Exec(query, quizID)).Scan(nil, nil)
}

func (repo *Repository) RemoveResultsByQuizID(quizID int64) error {
	query := repository_query.DeleteResultsByQuizID

	return new(datasource.DataSource).ExecSQL(repo.db.Exec(query, quizID)).Scan(nil, nil)
}

func (repo *Repository) GetQuizResultByUserID(quizID, userID int64) (qr repository_model.QuizResult, err error) {
	query := repository_query.SelectQuizResult

	row := func(idx int) utils.Array {
		return utils.Array{
			&qr.UserID,
			&qr.QuizID,
			&qr.Score,
			&qr.CompletedIn,
		}
	}

	err = new(datasource.DataSource).QuerySQL(repo.db.Queryx(query, quizID, userID)).Scan(row)

	return qr, err
}

func (repo *Repository) GetQuizProgressByUserID(quizID, userID int64) (qr []repository_model.QuizProgress, err error) {
	query := repository_query.SelectUserQuizProgress

	if err := repo.db.Select(&qr, query, quizID, userID); err != nil {
		return qr, err
	}

	return qr, nil
}

func (repo *Repository) SaveNewQuizResult(qr repository_model.QuizResult) error {
	query := repository_query.InsertQuizResult

	args := utils.Array{
		qr.UserID,
		qr.QuizID,
		qr.Score,
	}

	return new(datasource.DataSource).ExecSQL(repo.db.Exec(query, args...)).Scan(nil, nil)
}

func (repo *Repository) SaveNewQuizProgress(qp repository_model.QuizProgress) error {
	query := repository_query.InsertQuizProgress

	args := utils.Array{
		qp.UserID,
		qp.QuizID,
		qp.QuestionID,
		qp.Answer,
		qp.Correct,
	}

	return new(datasource.DataSource).ExecSQL(repo.db.Exec(query, args...)).Scan(nil, nil)
}

func (repo *Repository) GetUsersByPersonalData(surname, name, patronymic string) (users []repository_model.User, err error) {
	query := repository_query.SelectUsersBySurname

	args := utils.Array{
		surname,
	}

	if name != "" {
		query += " and u.name = ?"
		args = append(args, name)
	}

	if patronymic != "" {
		query += " and u.patronymic = ?"
		args = append(args, patronymic)
	}

	if err := repo.db.Select(&users, query, args...); err != nil {
		return users, err
	}
	return users, nil
}
