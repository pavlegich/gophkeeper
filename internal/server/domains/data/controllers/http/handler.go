// Package http contains object of data handler,
// functions for activating the data handler in controller
// and data handlers.
package http

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pavlegich/gophkeeper/internal/common/infra/config"
	"github.com/pavlegich/gophkeeper/internal/common/infra/logger"
	"github.com/pavlegich/gophkeeper/internal/server/domains/data"
	repo "github.com/pavlegich/gophkeeper/internal/server/domains/data/repository"
	errs "github.com/pavlegich/gophkeeper/internal/server/errors"
	"github.com/pavlegich/gophkeeper/internal/server/utils"
	"go.uber.org/zap"
)

// DataHandler contains objects for work
// with data handlers.
type DataHandler struct {
	Config  *config.ServerConfig
	Service data.Service
}

// Activate activates handler for data object.
func Activate(ctx context.Context, r *chi.Mux, cfg *config.ServerConfig, db *sql.DB) {
	s := data.NewDataService(ctx, repo.NewDataRepository(ctx, db))
	newHandler(ctx, r, cfg, s)
}

// newHandler initializes handler for data object.
func newHandler(ctx context.Context, r *chi.Mux, cfg *config.ServerConfig, s data.Service) {
	h := &DataHandler{
		Config:  cfg,
		Service: s,
	}
	r.Post("/api/user/data/{dataType}/{dataName}", h.HandleDataUpload)
	r.Get("/api/user/data/{dataType}/{dataName}", h.HandleDataValue)
	r.Put("/api/user/data/{dataType}/{dataName}", h.HandleDataUpdate)
	r.Delete("/api/user/data/{dataType}/{dataName}", h.HandleDataDelete)
}

// HandleDataUpload uploads new data into the storage.
func (h *DataHandler) HandleDataUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := utils.GetUserIDFromContext(ctx)
	idString := strconv.Itoa(userID)
	if err != nil {
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpload: get user id from context failed",
			zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := &data.Data{
		UserID: userID,
		Type:   chi.URLParam(r, "dataType"),
		Name:   chi.URLParam(r, "dataName"),
	}

	req, err = utils.GetMultipartDataFromRequest(ctx, r, req)

	if err != nil {
		if errors.Is(err, mime.ErrInvalidMediaParameter) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpload: read data from the request failed",
			zap.Error(err))
		return
	}

	err = h.Service.Create(ctx, req)
	if err != nil {
		if errors.Is(err, errs.ErrDataAlreadyUpload) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpload: upload new data failed",
			zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleDataValue writes requested data into response body
// if this data found in storage successfuly.
func (h *DataHandler) HandleDataValue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dType := chi.URLParam(r, "dataType")
	dName := chi.URLParam(r, "dataName")

	userID, err := utils.GetUserIDFromContext(ctx)
	idString := strconv.Itoa(userID)
	if err != nil {
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataValue: get user id from context failed",
			zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	storedData, err := h.Service.Unload(ctx, dType, dName)
	if err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataValue: unload requested data failed",
			zap.Error(err))
		return
	}

	if dType == "binary" {
		var buf bytes.Buffer
		multipartWriter := multipart.NewWriter(&buf)
		defer multipartWriter.Close()
		dataPart, err := multipartWriter.CreateFormField("file")
		if err != nil {
			logger.Log.With(zap.String("user_id", idString)).Error("HandleDataValue: create form field failed",
				zap.Error(err))
		}
		dataPart.Write(storedData.Data)
		w.Header().Set("Content-Type", multipartWriter.FormDataContentType())
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(storedData.Data))
}

// HandleDataUpdate updates the requested data in storage.
func (h *DataHandler) HandleDataUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := utils.GetUserIDFromContext(ctx)
	idString := strconv.Itoa(userID)
	if err != nil {
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpdate: get user id from context failed",
			zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := &data.Data{
		UserID: userID,
		Type:   chi.URLParam(r, "dataType"),
		Name:   chi.URLParam(r, "dataName"),
	}

	req, err = utils.GetMultipartDataFromRequest(ctx, r, req)

	if err != nil {
		if errors.Is(err, mime.ErrInvalidMediaParameter) {
			w.WriteHeader(http.StatusBadRequest)

		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpdate: read data from the request failed",
			zap.Error(err))
		return
	}

	err = h.Service.Edit(ctx, req)
	if err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpdate: update user's data failed",
			zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleDataDelete deletes requested data from the storage.
func (h *DataHandler) HandleDataDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dType := chi.URLParam(r, "dataType")
	dName := chi.URLParam(r, "dataName")

	userID, err := utils.GetUserIDFromContext(ctx)
	idString := strconv.Itoa(userID)
	if err != nil {
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpdate: get user id from context failed",
			zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.Service.Delete(ctx, dType, dName)
	if err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataDelete: delete requested data failed",
			zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
