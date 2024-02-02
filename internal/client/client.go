// Package client contains Client object and its methods.
package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/pavlegich/gophkeeper/internal/client/controllers"
	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
	"github.com/pavlegich/gophkeeper/internal/client/utils"
	"github.com/pavlegich/gophkeeper/internal/common/infra/config"
)

// Client contains client attributes.
type Client struct {
	rw         rwmanager.RWService
	controller *controllers.Controller
	config     *config.ClientConfig
}

// NewClient initializes controller and router, returns new client object.
func NewClient(ctx context.Context, ctrl *controllers.Controller, rw rwmanager.RWService, cfg *config.ClientConfig) (*Client, error) {
	return &Client{
		controller: ctrl,
		rw:         rw,
		config:     cfg,
	}, nil
}

// Serve starts listening and catching the commands from standart input.
func (c *Client) Serve(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := utils.DoWithRetryIfUnknown(ctx, c.controller.HandleCommand)
			if err != nil {
				got := utils.GetKnownErr(err)
				if errors.Is(err, errs.ErrUnknownCommand) {
					c.rw.WriteString(ctx, errs.ErrUnknownCommand.Error())
				}
				if got == nil {
					return fmt.Errorf("Serve: handle command failed %w", err)
				}
				c.rw.WriteString(ctx, got.Error())
			}
		}
	}
}
