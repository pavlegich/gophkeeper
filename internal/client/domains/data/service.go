package data

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pavlegich/gophkeeper/internal/client/domains/data/readers"
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
func (s *DataService) CreateOrUpdate(ctx context.Context) error {
	d, err := readDataTypeAndName(ctx, s.rw)
	if err != nil {
		return fmt.Errorf("CreateOrUpdate: couldn't read data type and name %w", err)
	}

	// Read data and put it into multipart
	var buf bytes.Buffer
	multipartWriter := multipart.NewWriter(&buf)
	defer multipartWriter.Close()
	err = createMultipartData(ctx, s.rw, multipartWriter, d)
	if err != nil {
		return fmt.Errorf("CreateOrUpdate: create multipart data failed %w", err)
	}

	// Prepare request
	target := s.cfg.Address + "/api/user/data/" + d.Type + "/" + d.Name

	act, err := utils.GetActionFromContext(ctx)
	if err != nil {
		return fmt.Errorf("CreateOrUpdate: get action from context failed %w", err)
	}

	var method string
	switch act {
	case "create":
		method = http.MethodPost
	case "update":
		method = http.MethodPut
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, target, &buf)
	if err != nil {
		return fmt.Errorf("CreateOrUpdate: new request failed %w", err)
	}

	// Add values into the header
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

	s.rw.Writeln(ctx, utils.Success)

	return nil
}

// GetValue sends request to the server to get value that was requested by the client.
func (s *DataService) GetValue(ctx context.Context) error {
	d, err := readDataTypeAndName(ctx, s.rw)
	if err != nil {
		return fmt.Errorf("GetValue: couldn't read data type and name %w", err)
	}

	// Prepare request
	target := s.cfg.Address + "/api/user/data/" + d.Type + "/" + d.Name

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
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

	if d.Type == "binary" {
		s.rw.Write(ctx, "Type path for save file: ")
		path, err := s.rw.Read(ctx)
		if err != nil {
			return fmt.Errorf("GetValue: read file path failed %w", err)
		}
		err = utils.SaveFromMultipartToFile(ctx, resp, path)
		if err != nil {
			return fmt.Errorf("GetValue: save file failed %w", err)
		}

		s.rw.Writeln(ctx, utils.Success)

		return nil
	}

	var buf bytes.Buffer

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("GetValue: read data from body failed %w", err)
	}
	s.rw.Writeln(ctx, buf.String())

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
	target := s.cfg.Address + "/api/user/data/" + d.Type + "/" + d.Name

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, target, nil)
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

	s.rw.Writeln(ctx, utils.Success)

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

// createMultipartData put data parts into multipart fields.
func createMultipartData(ctx context.Context, rw rwmanager.RWService, mpwriter *multipart.Writer, d *Data) error {
	// Data
	var part []byte
	var err error
	var dataReader DataReader
	switch d.Type {
	case "credentials":
		dataReader = readers.NewCredentialsReader(ctx, rw)
	case "card":
		dataReader = readers.NewCardReader(ctx, rw)
	case "text":
		dataReader = readers.NewTextReader(ctx, rw)
	}

	part, err = dataReader.Read(ctx)
	if err != nil {
		return fmt.Errorf("createMultipartData: couldn't read %s %w", d.Type, err)
	}

	// Put data into the dataPart
	var dataPart io.Writer
	if d.Type == "binary" {
		rw.Write(ctx, "Type absolute path to file: ")
		path, err := rw.Read(ctx)
		if err != nil {
			return fmt.Errorf("createMultipartData: couldn't read file path %w", err)
		}
		dataPart, err = mpwriter.CreateFormFile("file", path)
		if err != nil {
			return fmt.Errorf("createMultipartData: couldn't create form file %w", err)
		}
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("createMultipartData: %w", errs.ErrInvalidFilePath)
		}
		defer file.Close()
		_, err = io.Copy(dataPart, file)
		if err != nil {
			return fmt.Errorf("createMultipartData: copy file failed %w", err)
		}
	} else {
		dataPart, err = mpwriter.CreateFormField("data")
		if err != nil {
			return fmt.Errorf("createMultipartData: create multipart data form failed %w", err)
		}
		dataPart.Write(part)
	}

	// Metadata
	metaPart, err := mpwriter.CreateFormField("metadata")
	if err != nil {
		return fmt.Errorf("createMultipartData: create multipart metadata form failed %w", err)
	}
	metaReader := readers.NewMetadataReader(ctx, rw)
	metaBytes, err := metaReader.Read(ctx)
	if err != nil {
		return fmt.Errorf("createMultipartData: read metadata failed %w", err)
	}
	metaPart.Write(metaBytes)

	return nil
}
