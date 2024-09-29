package main

import (
	"errors"

	"github.com/amengdv/gator/internal/config"
	"github.com/amengdv/gator/internal/database"
)

type state struct {
    db *database.Queries
    cfg *config.Config
}

type command struct {
    name string
    args []string
}

type commands struct {
    handlers map[string]func(*state, command) error
}

func (cmds *commands) register(name string, f func(*state, command) error) {
    cmds.handlers[name] = f
}

func (cmds *commands) run(s *state, cmd command) error {
    f, ok := cmds.handlers[cmd.name]
    if !ok {
        return errors.New("Command does not exist")
    }

    err := f(s, cmd)
    if err != nil {
        return err
    }

    return nil
}
