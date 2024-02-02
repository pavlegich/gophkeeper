package data

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/pavlegich/gophkeeper/internal/client/domains/data/models"
	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
	"github.com/pavlegich/gophkeeper/internal/client/utils"
	"github.com/pavlegich/gophkeeper/internal/common/infra/config"
)

// DataService contains objects for data service.
type DataService struct {
	rw  rwmanager.RWService
	cfg *config.ClientConfig
}

// NewDataService creates and returns new data service.
func NewDataService(ctx context.Context, rw rwmanager.RWService, cfg *config.ClientConfig) *DataService {
	return &DataService{
		rw:  rw,
		cfg: cfg,
	}
}

// CreateOrUpdate reads information about data from the input,
// sends request to the server with requested action.
func (s *DataService) CreateOrUpdate(ctx context.Context, act string) error {
	d, err := readDataTypeAndName(ctx, s.rw)
	if err != nil {
		return fmt.Errorf("CreateOrUpdate: couldn't read data type and name %w", err)
	}

	// Read data and put it into multipart
	var buf bytes.Buffer
	multipartWriter := multipart.NewWriter(&buf)
	defer multipartWriter.Close()

	// Data
	dataPart, err := multipartWriter.CreateFormField("data")
	if err != nil {
		fmt.Errorf("CreateOrUpdate: create multipart data form failed %w", err)
	}
	var part []byte
	switch d.Type {
	case "credentials":
		part, err = models.ReadCredentials(ctx, s.rw)
		if err != nil {
			return fmt.Errorf("CreateOrUpdate: couldn't read credentials %w", err)
		}
	case "card":
		part, err = models.ReadCardDetails(ctx, s.rw)
		if err != nil {
			return fmt.Errorf("CreateOrUpdate: couldn't read card details %w", err)
		}
	case "text":
		part, err = models.ReadText(ctx, s.rw)
		if err != nil {
			return fmt.Errorf("CreateOrUpdate: couldn't read text %w", err)
		}
	}
	dataPart.Write(part)

	// Metadata
	metaPart, err := multipartWriter.CreateFormField("metadata")
	if err != nil {
		return fmt.Errorf("CreateOrUpdate: create multipart metadata form failed %w", err)
	}
	meta, err := models.ReadMetadata(ctx, s.rw)
	if err != nil {
		return fmt.Errorf("CreateOrUpdate: read metadata failed %w", err)
	}
	metaPart.Write(meta)

	// Prepare request
	target := "http://" + s.cfg.Address + "/api/user/data/" + d.Type + "/" + d.Name

	ctxReq, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var method string
	switch act {
	case "create":
		method = http.MethodPost
	case "update":
		method = http.MethodPut
	}

	req, err := http.NewRequestWithContext(ctxReq, method, target, &buf)
	if err != nil {
		return fmt.Errorf("CreateOrUpdate: new request failed %w", err)
	}

	if s.cfg.Cookie != nil {
		req.AddCookie(s.cfg.Cookie)
	}
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// Send request
	resp, err := utils.DoRequestWithRetry(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateOrUpdate: send request failed %w", err)
	}
	defer resp.Body.Close()

	// Check response
	err = utils.CheckStatusCode(resp.StatusCode)
	if err != nil {
		return fmt.Errorf("CreateOrUpdate: create data failed %w", err)
	}

	s.rw.WriteString(ctx, utils.Success)

	return nil
}

func (s *DataService) GetValue(ctx context.Context) error {
	d, err := readDataTypeAndName(ctx, s.rw)
	if err != nil {
		return fmt.Errorf("GetValue: couldn't read data type and name %w", err)
	}

	// Prepare request
	target := "http://" + s.cfg.Address + "/api/user/data/" + d.Type + "/" + d.Name

	ctxReq, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxReq, http.MethodGet, target, nil)
	if err != nil {
		return fmt.Errorf("GetValue: new request failed %w", err)
	}

	if s.cfg.Cookie != nil {
		req.AddCookie(s.cfg.Cookie)
	}

	// Send request
	resp, err := utils.DoRequestWithRetry(ctx, req)
	if err != nil {
		return fmt.Errorf("GetValue: send request failed %w", err)
	}
	defer resp.Body.Close()

	// Check response
	err = utils.CheckStatusCode(resp.StatusCode)
	if err != nil {
		return fmt.Errorf("GetValue: create data failed %w", err)
	}

	var buf bytes.Buffer

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("GetValue: read data from body failed %w", err)
	}
	s.rw.WriteString(ctx, buf.String())

	return nil
}

// Delete reads information about data from the input,
// sends request to the server to delete requested data.
func (s *DataService) Delete(ctx context.Context) error {
	d, err := readDataTypeAndName(ctx, s.rw)
	if err != nil {
		return fmt.Errorf("Delete: couldn't read data type and name %w", err)
	}

	// Prepare request
	target := "http://" + s.cfg.Address + "/api/user/data/" + d.Type + "/" + d.Name

	ctxReq, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxReq, http.MethodDelete, target, nil)
	if err != nil {
		return fmt.Errorf("Delete: new request failed %w", err)
	}

	if s.cfg.Cookie != nil {
		req.AddCookie(s.cfg.Cookie)
	}

	// Send request
	resp, err := utils.DoRequestWithRetry(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateOrUpdate: send request failed %w", err)
	}
	defer resp.Body.Close()

	// Check response
	err = utils.CheckStatusCode(resp.StatusCode)
	if err != nil {
		return fmt.Errorf("CreateOrUpdate: create data failed %w", err)
	}

	s.rw.WriteString(ctx, utils.Success)

	return nil
}

// readDataTypeAndName reads from the input and returns data type and data name.
func readDataTypeAndName(ctx context.Context, rw rwmanager.RWService) (*Data, error) {
	d := &Data{}
	var err error

	// Read data type
	rw.Write(ctx, "Data type (credentials/card/text/binary): ")
	d.Type, err = rw.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("readDataTypeAndName: couldn't read data type %w", err)
	}
	d.Type = strings.ToLower(d.Type)
	if !utils.IsValidDataType(d.Type) {
		return nil, fmt.Errorf("readDataTypeAndName: %w", errs.ErrInvalidDataType)
	}

	// Read data name
	rw.Write(ctx, "Data name: ")
	d.Name, err = rw.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("readDataTypeAndName: couldn't read data name %w", err)
	}

	return d, nil
}
