// Package http contains object of data handler,
// functions for activating the data handler in controller
// and data handlers.
package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pavlegich/gophkeeper/internal/infra/config"
	"github.com/pavlegich/gophkeeper/internal/infra/logger"
	"github.com/pavlegich/gophkeeper/internal/server/domains/data"
	repo "github.com/pavlegich/gophkeeper/internal/server/domains/data/repository"
	errs "github.com/pavlegich/gophkeeper/internal/server/errors"
	"github.com/pavlegich/gophkeeper/internal/utils"
	"go.uber.org/zap"
)

// DataHandler contains objects for work
// with data handlers.
type DataHandler struct {
	Config  *config.ServerConfig
	Service data.Service
}

// Activate activates handler for data object.
func Activate(r *chi.Mux, cfg *config.ServerConfig, db *sql.DB) {
	s := data.NewDataService(repo.NewDataRepository(db))
	newHandler(r, cfg, s)
}

// newHandler initializes handler for data object.
func newHandler(r *chi.Mux, cfg *config.ServerConfig, s data.Service) {
	h := &DataHandler{
		Config:  cfg,
		Service: s,
	}
	r.Post("/api/user/data/new", h.HandleDataUpload)
}

// HandleDataUpload uploads new data into the storage.
func (h *DataHandler) HandleDataUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req data.Data
	var buf bytes.Buffer

	userID, err := utils.GetUserIDFromContext(ctx)
	idString := strconv.Itoa(userID)
	if err != nil {
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpload: get user id from context failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpload: read request body failed",
			zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpload: unmarshal data failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.UserID = userID
	// ???
	req.Type = strings.ToUpper(req.Type)

	err = h.Service.Create(ctx, &req)
	if err != nil {
		if errors.Is(err, errs.ErrDataAlreadyUpload) {
			w.WriteHeader(http.StatusConflict)
		} else if errors.Is(err, errs.ErrDataTypeIncorrect) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpload: upload new data failed",
			zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
