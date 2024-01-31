package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/pavlegich/gophkeeper/internal/server/domains/data"
)

// GetMultipartDataFromRequest reads multipart fields from the request and returns
// the data object with the obtained multipart data.
func GetMultipartDataFromRequest(ctx context.Context, r *http.Request, d *data.Data) (*data.Data, error) {
	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("GetMultipartDataFromRequest: couldn't get media type %w", err)
	}

	multiparted := d

	if strings.HasPrefix(mediaType, "multipart/") {
		multipartReader := multipart.NewReader(r.Body, params["boundary"])
		defer r.Body.Close()

		for {
			field, err := multipartReader.NextPart()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return nil, fmt.Errorf("GetMultipartDataFromRequest: get next multi part failed %w", err)
			}
			defer field.Close()

			switch field.FormName() {
			case "data":
				multiparted.Data, err = io.ReadAll(field)
				if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
					return nil, fmt.Errorf("GetMultipartDataFromRequest: couldn't read data from the field data %w", err)
				}
			case "metadata":
				multiparted.Metadata, err = io.ReadAll(field)
				if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
					return nil, fmt.Errorf("GetMultipartDataFromRequest: couldn't read data from the field metadata %w", err)
				}
			}
		}
	}

	return multiparted, nil
}
