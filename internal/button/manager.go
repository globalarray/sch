package button

import (
	"log/slog"
	"sync"
)

type Manager struct {
	log            *slog.Logger
	buttonsRegOnce sync.Once
	buttons        map[string]Button
}

var (
	once      sync.Once
	managerMu sync.Mutex
	manager   *Manager
)

func Mgr() *Manager {
	managerMu.Lock()
	defer managerMu.Unlock()
	return manager
}

func NewManager(log *slog.Logger) *Manager {
	once.Do(func() {
		manager = &Manager{log: log.With(slog.String("level", "button")), buttons: map[string]Button{}}
		manager.regButtons()
	})

	return Mgr()
}

func (m *Manager) Get(endpoint string) (Button, bool) {
	b, ok := m.buttons[endpoint]

	return b, ok
}

func (m *Manager) regButtons() {
	m.buttonsRegOnce.Do(func() {
		bttns := []Button{
			&CreateInvitationKey{log: m.log.With(slog.String("btn", "create_invitation_key"))},
			&FillPersonalData{log: m.log.With(slog.String("btn", "fill_personal_data"))},
			&SetRoleKey{log: m.log.With(slog.String("btn", "set_role_key"))},
			&RemoveInvitationKey{log: m.log.With(slog.String("btn", "remove_invitation_key"))},
			&CreateQuiz{log: m.log.With(slog.String("btn", "create_quiz"))},
			&AddQuestionQuiz{log: m.log.With(slog.String("btn", "add_question_quiz"))},
			&RemoveQuestionQuiz{log: m.log.With(slog.String("btn", "remove_question_quiz"))},
			&GetInfoQuiz{log: m.log.With(slog.String("btn", "get_info_quiz"))},
			&ListQuiz{log: m.log.With(slog.String("btn", "list_quiz"))},
			&RemoveQuiz{log: m.log.With(slog.String("btn", "remove_quiz"))},
			&AnswerQuestion{log: m.log.With(slog.String("btn", "answer_question"))},
		}

		for _, b := range bttns {
			m.buttons[b.Endpoint()] = b
		}
	})
}
