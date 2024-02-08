// Package controllers contains client controller object and
// methods for client work.
package controllers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pavlegich/gophkeeper/internal/client/domains/data"
	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	"github.com/pavlegich/gophkeeper/internal/client/domains/user"
	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
	"github.com/pavlegich/gophkeeper/internal/client/utils"
	"github.com/pavlegich/gophkeeper/internal/common/infra/config"
)

// Controller contains configuration for building the client app.
type Controller struct {
	rw   rwmanager.RWService
	cfg  *config.ClientConfig
	user user.Service
	data data.Service
}

// NewController creates and returns new client controller.
func NewController(ctx context.Context, rw rwmanager.RWService, cfg *config.ClientConfig) *Controller {
	userService := user.NewUserService(ctx, rw, cfg)
	dataService := data.NewDataService(ctx, rw, cfg)

	return &Controller{
		rw:   rw,
		cfg:  cfg,
		user: userService,
		data: dataService,
	}
}

// HandleCommand handles commands from the input, selects and does the requested action.
func (c *Controller) HandleCommand(ctx context.Context) error {
	c.rw.Write(ctx, "Type the command, or exit: ")
	act, err := c.rw.Read(ctx)
	if err != nil && !errors.Is(err, errs.ErrEmptyInput) {
		return fmt.Errorf("HandleCommand: read command failed %w", err)
	}

	act = strings.ToLower(act)
	ctx = context.WithValue(ctx, utils.ContextActionKey, act)
	var clientAct func(ctx context.Context) error

	switch act {
	case "register":
		err := c.user.Register(ctx)
		if err != nil {
			return fmt.Errorf("HandleCommand: register user failed %w", err)
		}
	case "login":
		err := c.user.Login(ctx)
		if err != nil {
			return fmt.Errorf("HandleCommand: login user failed %w", err)
		}
	case "create", "update":
		clientAct = c.data.CreateOrUpdate
	case "get":
		clientAct = c.data.GetValue
	case "delete":
		clientAct = c.data.Delete
	case "exit":
		return fmt.Errorf("HandleCommand: %w", errs.ErrExit)
	default:
		return fmt.Errorf("HandleCommand: %w", errs.ErrUnknownCommand)
	}

	err = utils.DoWithRetryIfEmpty(ctx, c.rw, clientAct)
	if err != nil {
		return fmt.Errorf("HandleCommand: %s data failed %w", act, err)
	}
	return nil
}
