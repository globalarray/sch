package cmd

import (
	"log/slog"
	"sync"
)

type Manager struct {
	log             *slog.Logger
	commands        map[string]Command
	commandsRegOnce sync.Once
}

var (
	once      sync.Once
	managerMu sync.Mutex
	mgr       *Manager
)

func Mgr() *Manager {
	managerMu.Lock()
	defer managerMu.Unlock()
	return mgr
}

func NewManager(log *slog.Logger) *Manager {
	once.Do(func() {
		mgr = &Manager{log: log.With(slog.String("level", "cmd")), commands: make(map[string]Command)}
		mgr.regCommands()
	})

	return Mgr()
}

func (m *Manager) regCommands() {
	m.commandsRegOnce.Do(func() {
		cmds := []Command{
			&Start{log: m.log.With(slog.String("level", "cmd/start"))},
		}

		for _, cmd := range cmds {
			m.commands[cmd.Endpoint()] = cmd
		}
	})
}

func (m *Manager) Get(endpoint string) (Command, bool) {
	c, ok := m.commands[endpoint]

	return c, ok
}
