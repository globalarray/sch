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

	if err := new(datasource.DataSource).QuerySQL(repo.db.Queryx(query, id)).Scan(row); err != nil {
		return q, err
	}

	return q, nil
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
