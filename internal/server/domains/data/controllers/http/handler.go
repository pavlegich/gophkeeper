// Package http contains object of data handler,
// functions for activating the data handler in controller
// and data handlers.
package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
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
	r.Post("/api/user/data/{dataType}/{dataName}", h.HandleDataUpload)
	r.Get("/api/user/data/{dataType}/{dataName}", h.HandleDataValue)
	r.Put("/api/user/data", h.HandleDataUpdate)
	r.Delete("/api/user/data/{dataType}/{dataName}", h.HandleDataDelete)
}

// HandleDataUpload uploads new data into the storage.
func (h *DataHandler) HandleDataUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := utils.GetUserIDFromContext(ctx)
	idString := strconv.Itoa(userID)
	if err != nil {
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpload: get user id from context failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := data.Data{
		UserID: userID,
		Type:   chi.URLParam(r, "dataType"),
		Name:   chi.URLParam(r, "dataName"),
	}

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpload: read request media type failed",
			zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Type {
	case "text", "binary":
		if strings.HasPrefix(mediaType, "multipart/") {
			multipartReader := multipart.NewReader(r.Body, params["boundary"])
			defer r.Body.Close()

			for {
				field, err := multipartReader.NextPart()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpload: get next multi part failed",
						zap.Error(err))
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				defer field.Close()

				// pay attention to read large file
				switch field.FormName() {
				case "data":
					req.Data, err = io.ReadAll(field)
					if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
						logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpload: couldn't read data from the field data",
							zap.Error(err))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				case "metadata":
					req.Metadata, err = io.ReadAll(field)
					if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
						logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpload: couldn't read data from the field metadata",
							zap.Error(err))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	err = h.Service.Create(ctx, &req)
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

// HandleDataValue writes requested data in response body
// if this data found in storage successfuly.
func (h *DataHandler) HandleDataValue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dType := chi.URLParam(r, "dataType")
	dName := chi.URLParam(r, "dataName")

	userID, err := utils.GetUserIDFromContext(ctx)
	idString := strconv.Itoa(userID)
	if err != nil {
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataValue: get user id from context failed")
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

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(storedData.Data))
}

// HandleDataUpdate updates the requested data in storage.
func (h *DataHandler) HandleDataUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req data.Data
	var buf bytes.Buffer

	userID, err := utils.GetUserIDFromContext(ctx)
	idString := strconv.Itoa(userID)
	if err != nil {
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpdate: get user id from context failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpdate: read request body failed",
			zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpdate: unmarshal data failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.UserID = userID

	err = h.Service.Edit(ctx, &req)
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
		logger.Log.With(zap.String("user_id", idString)).Error("HandleDataUpdate: get user id from context failed")
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
