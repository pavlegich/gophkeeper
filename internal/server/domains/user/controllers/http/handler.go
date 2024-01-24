// Package http contains object of user handler,
// functions for activating the user handler in controller
// and user handlers.
package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pavlegich/gophkeeper/internal/infra/config"
	"github.com/pavlegich/gophkeeper/internal/infra/logger"
	"github.com/pavlegich/gophkeeper/internal/server/domains/user"
	repo "github.com/pavlegich/gophkeeper/internal/server/domains/user/repository"
	errs "github.com/pavlegich/gophkeeper/internal/server/errors"
	"go.uber.org/zap"
)

// UserHandler contatins objects for work
// with user handlers.
type UserHandler struct {
	Config  *config.ServerConfig
	Service user.Service
}

// Activate activates handler for user.
func Activate(r *chi.Mux, cfg *config.ServerConfig, db *sql.DB) {
	s := user.NewUserService(repo.NewUserRepository(db))
	newHandler(r, cfg, s)
}

// newHandler initializes handler for user.
func newHandler(r *chi.Mux, cfg *config.ServerConfig, s user.Service) {
	h := &UserHandler{
		Config:  cfg,
		Service: s,
	}
	r.Post("/api/user/register", h.HandleRegister)
	r.Post("/api/user/login", h.HandleLogin)
	r.Post("/api/user/logout", h.HandleLogout)
}

// HandleRegister registers new user.
func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req user.User
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.Error("HandleRegister: read request body failed",
			zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		logger.Log.Error("HandleRegister: request unmarshal failed",
			zap.String("body", buf.String()),
			zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.Service.Register(ctx, &req)
	if err != nil {
		if errors.Is(err, errs.ErrLoginBusy) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		logger.Log.Error("HandleRegister: user register failed",
			zap.Error(err))
		return
	}

	cookie := &http.Cookie{
		Name:  "auth",
		Value: "secret",
		Path:  "/api/user/",
		// Secure: true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

// HandleLogin handles user login or authorization with received data.
func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req user.User
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.Error("HandleLogin: read request body failed",
			zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		logger.Log.Error("HandleLogin: request unmarshal failed",
			zap.String("body", buf.String()),
			zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.Service.Login(ctx, &req)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
		} else if errors.Is(err, errs.ErrPasswordNotMatch) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		logger.Log.Error("HandleLogin: user login failed",
			zap.Error(err))
		return
	}

	cookie := http.Cookie{
		Name:  "auth",
		Value: "secret",
		Path:  "/api/user/",
		// Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

// HandleLogout log user out from the service.
func (h *UserHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: "auth",
		Path: "/api/user/",
		// Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
}
