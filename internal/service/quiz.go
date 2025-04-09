package service

import (
	"benzo/internal/lang"
	"benzo/internal/repository"
	"benzo/internal/repository/repository_model"
	"benzo/pkg/i18n"
	"errors"
	"fmt"
	tele "gopkg.in/telebot.v4"
	"iter"
	"math/rand"
	"slices"
	"strings"
	"time"
)

type QuizService struct {
}

var (
	ErrQuizWithoutQuestions = errors.New("quiz without questions")
)

var (
	quizService *QuizService
)

func Quiz() *QuizService {
	return quizService
}

func (qs *QuizService) ProcessQuiz(ctx tele.Context, quizID, userID int64, languageCode string) error {
	q, err := repository.Repo().GetQuizByID(quizID)

	if err != nil {
		return err
	}

	questions, err := repository.Repo().GetQuestionsByQuizID(quizID)

	if err != nil {
		return err
	}

	if len(questions) == 0 {
		return ErrQuizWithoutQuestions
	}

	qResult, err := repository.Repo().GetQuizResultByUserID(quizID, userID)

	if err != nil {
		return err
	}

	if qResult.UserID == userID {
		return ctx.Send(i18n.Translatef(lang.QuizAlreadyCompleted, languageCode, q.Name))
	}

	progressedQuestions, err := repository.Repo().GetQuizProgressByUserID(quizID, userID)

	if err != nil {
		return err
	}

	correctQuestionsCnt := len(slices.Collect(qs.getCorrectQuestions(progressedQuestions)))

	if len(progressedQuestions) >= len(questions) {
		if err := repository.Repo().SaveNewQuizResult(repository_model.NewQuizResult(quizID, userID, correctQuestionsCnt)); err != nil {
			return err
		}

		return ctx.Send(i18n.Translatef(lang.QuizCompletedMessage, languageCode, q.Name, correctQuestionsCnt, len(questions)))
	}
	var passedQuestions []int64

	for _, q := range progressedQuestions {
		passedQuestions = append(passedQuestions, q.QuestionID)
	}

	var nextQuestionID int64
	var nextQuestionNum int

	for idx, question := range questions {
		if !slices.Contains(passedQuestions, question.ID) {
			nextQuestionID = question.ID
			nextQuestionNum = idx + 1
			break
		}
	}

	question, err := repository.Repo().GetQuestionByID(nextQuestionID)

	if err != nil {
		return err
	}

	var response []string

	response = append(response, i18n.Translatef(lang.QuestionTitle, languageCode, nextQuestionNum, len(questions)))
	response = append(response, "")
	response = append(response, question.Question)

	selector := &tele.ReplyMarkup{}

	var rows []tele.Row

	answers := strings.Split(question.Answers, ";")

	randAnswers := make([]string, len(answers))

	copy(randAnswers, answers)

	r := rand.New(rand.NewSource(time.Now().Unix()))

	for n := len(randAnswers); n > 0; n-- {
		randIndex := r.Intn(n)
		// We swap the value at index n-1 and the random index
		// to move our randomly chosen value to the end of the
		// slice, and to move the value that was at n-1 into our
		// unshuffled portion of the slice.
		randAnswers[n-1], randAnswers[randIndex] = randAnswers[randIndex], randAnswers[n-1]
	}

	for _, ans := range randAnswers {
		realIdx := slices.Index(answers, ans)

		ansBtn := selector.Data(ans, fmt.Sprintf("question_answer-%d-%d", question.ID, realIdx))

		if len(rows) == 0 {
			rows = append(rows, selector.Row(ansBtn))
			continue
		}

		row := rows[len(rows)-1]

		if len(row) == 2 {
			rows = append(rows, selector.Row(ansBtn))
			continue
		}
		rows[len(rows)-1] = append(row, ansBtn)
	}

	selector.Inline(rows...)

	return ctx.Send(strings.Join(response, "\n"), selector)
}

func (qs *QuizService) getCorrectQuestions(progressedQuestions []repository_model.QuizProgress) iter.Seq[repository_model.QuizProgress] {
	return func(yield func(repository_model.QuizProgress) bool) {
		for _, q := range progressedQuestions {
			if q.Correct && !yield(q) {
				return
			}
		}
	}
}

/* (qs *QuizService) GetQuizPassedUsers(quizID int64) (map[int64]int, error) {
	questions, err := repository.Repo().GetQuestionsByQuizID(quizID)

	if err != nil {
		return nil, err
	}

	if len(questions) == 0 {
		return nil, ErrQuizWithoutQuestions
	}
}*/
