// Package controllers contains client controller object and
// methods for client work.
package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	"github.com/pavlegich/gophkeeper/internal/client/domains/user"
	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
	"github.com/pavlegich/gophkeeper/internal/infra/config"
)

// Controller contains configuration for building the client app.
type Controller struct {
	rw   rwmanager.RWService
	cfg  *config.ClientConfig
	user user.Service
}

// NewController creates and returns new client controller.
func NewController(ctx context.Context, rw rwmanager.RWService, cfg *config.ClientConfig) *Controller {
	userService := user.NewUserService(ctx, rw, cfg)

	return &Controller{
		rw:   rw,
		cfg:  cfg,
		user: userService,
	}
}

// HandleCommand handles commands from the input, selects and does the requested action.
func (c *Controller) HandleCommand(ctx context.Context) error {
	c.rw.WriteString(ctx, "Type the command (or exit): ")
	act, err := c.rw.Read(ctx)
	if err != nil {
		return fmt.Errorf("HandleCommand: read command failed %w", err)
	}
	commands := strings.Split(act, " ")

	switch strings.ToLower(commands[0]) {
	case "register":
		err := c.user.Register(ctx)
		if err != nil {
			return fmt.Errorf("HandleCommand: register user failed %w", err)
		}
	case "exit":
		return fmt.Errorf("NewClient: %w", errs.ErrExit)
	default:
		return fmt.Errorf("NewClient: unknown command")
	}
	return nil
}
