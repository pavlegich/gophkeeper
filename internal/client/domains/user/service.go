package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	"github.com/pavlegich/gophkeeper/internal/client/utils"
	"github.com/pavlegich/gophkeeper/internal/common/infra/config"
)

// UserService contains objects for user service.
type UserService struct {
	rw  rwmanager.RWService
	cfg *config.ClientConfig
}

// NewUserService creates and returns new user service.
func NewUserService(ctx context.Context, rw rwmanager.RWService, cfg *config.ClientConfig) *UserService {
	return &UserService{
		rw:  rw,
		cfg: cfg,
	}
}

// Register requests the server for user registration.
func (s *UserService) Register(ctx context.Context) error {
	u := &User{}
	var err error

	s.rw.Write(ctx, "Login: ")
	u.Login, err = s.rw.Read(ctx)
	if err != nil {
		return fmt.Errorf("Register: couldn't read login %w", err)
	}

	s.rw.Write(ctx, "Password: ")
	u.Password, err = s.rw.Read(ctx)
	if err != nil {
		return fmt.Errorf("Register: couldn't read password %w", err)
	}

	body, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("Register: marshal user failed %w", err)
	}

	target := "http://" + s.cfg.Address + "/api/user/register"
	ctxReq, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctxReq, http.MethodPost, target, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Register: new request failed %w", err)
	}

	resp, err := utils.DoRequestWithRetry(ctx, req)
	if err != nil {
		return fmt.Errorf("Register: send request failed %w", err)
	}
	defer resp.Body.Close()

	err = utils.CheckStatusCode(resp.StatusCode)
	if err != nil {
		return fmt.Errorf("Register: user register failed %w", err)
	}

	for _, c := range resp.Cookies() {
		if c.Name == "auth" {
			s.cfg.Cookie = c
		}
	}
	if s.cfg == nil {
		return fmt.Errorf("Register: cookie not found")
	}

	s.rw.WriteString(ctx, utils.Success)

	return nil
}

// Login requests the server for user login.
func (s *UserService) Login(ctx context.Context) error {
	u := &User{}
	var err error

	s.rw.Write(ctx, "Login: ")
	u.Login, err = s.rw.Read(ctx)
	if err != nil {
		return fmt.Errorf("Login: couldn't read login %w", err)
	}

	s.rw.Write(ctx, "Password: ")
	u.Password, err = s.rw.Read(ctx)
	if err != nil {
		return fmt.Errorf("Login: couldn't read password %w", err)
	}

	body, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("Login: marshal user failed %w", err)
	}

	target := "http://" + s.cfg.Address + "/api/user/login"
	ctxReq, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctxReq, http.MethodPost, target, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Login: new request failed %w", err)
	}

	resp, err := utils.DoRequestWithRetry(ctx, req)
	if err != nil {
		return fmt.Errorf("Login: send request failed %w", err)
	}
	defer resp.Body.Close()

	err = utils.CheckStatusCode(resp.StatusCode)
	if err != nil {
		return fmt.Errorf("Login: user register failed %w", err)
	}

	for _, c := range resp.Cookies() {
		if c.Name == "auth" {
			s.cfg.Cookie = c
		}
	}
	if s.cfg == nil {
		return fmt.Errorf("Login: cookie not found")
	}

	s.rw.WriteString(ctx, utils.Success)

	return nil
}
