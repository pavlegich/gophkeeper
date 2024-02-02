package data

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

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

// Create reads information about data from the input,
// sends request to the server with requested action.
func (s *DataService) Create(ctx context.Context) error {
	d := &Data{}
	var err error

	s.rw.Write(ctx, "Data type (credentials/card/text/binary): ")
	d.Type, err = s.rw.Read(ctx)
	if err != nil {
		return fmt.Errorf("Create: couldn't read data type %w", err)
	}
	d.Type = strings.ToLower(d.Type)
	if !utils.IsValidDataType(d.Type) {
		return fmt.Errorf("Create: %w", errs.ErrInvalidDataType)
	}

	s.rw.Write(ctx, "Data name: ")
	d.Name, err = s.rw.Read(ctx)
	if err != nil {
		return fmt.Errorf("Create: couldn't read password %w", err)
	}
	if d.Name == "" {
		return fmt.Errorf("Create: %w", errs.ErrEmptyInput)
	}

	var buf bytes.Buffer
	multipartWriter := multipart.NewWriter(&buf)
	defer multipartWriter.Close()

	dataPart, _ := multipartWriter.CreateFormField("data")
	var part []byte
	switch d.Type {
	case "credentials":
		part, err = ReadCredentials(ctx, s.rw)
		if err != nil {
			return fmt.Errorf("Create: couldn't read credentials %w", err)
		}
	case "card":
		part, err = ReadCardDetails(ctx, s.rw)
		if err != nil {
			return fmt.Errorf("Create: couldn't read card details %w", err)
		}
	}
	dataPart.Write(part)

	target := "http://" + s.cfg.Address + "/api/user/data/" + d.Type + "/" + d.Name

	ctxReq, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctxReq, http.MethodPost, target, &buf)
	if err != nil {
		return fmt.Errorf("Create: new request failed %w", err)
	}

	if s.cfg.Cookie != nil {
		req.AddCookie(s.cfg.Cookie)
	}
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	resp, err := utils.GetRequestWithRetry(ctx, req)
	if err != nil {
		return fmt.Errorf("Create: send request failed %w", err)
	}
	resp.Body.Close()

	err = utils.CheckStatusCode(resp.StatusCode)
	if err != nil {
		return fmt.Errorf("Create: create data failed %w", err)
	}

	s.rw.WriteString(ctx, utils.Success)

	return nil
}
func (s *DataService) Update(ctx context.Context) error {
	return nil
}
func (s *DataService) GetValue(ctx context.Context) error {
	return nil
}
func (s *DataService) Delete(ctx context.Context) error {
	return nil
}
