package utils

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"reflect"
	"testing"

	"github.com/pavlegich/gophkeeper/internal/server/domains/data"
)

func createRequest(d *data.Data) *http.Request {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	defer mw.Close()

	dataPart, _ := mw.CreateFormField("data")
	dataPart.Write(d.Data)

	metaPart, _ := mw.CreateFormField("metadata")
	metaPart.Write(d.Metadata)

	r, _ := http.NewRequest("", "", &buf)
	r.Header.Set("Content-Type", mw.FormDataContentType())

	return r
}

func TestGetMultipartDataFromRequest(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
		d   *data.Data
	}
	tests := []struct {
		name    string
		args    args
		want    *data.Data
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx: ctx,
				d: &data.Data{
					Data:     []byte(`test`),
					Metadata: []byte(`{"meta": "meta"}`),
				},
			},
			want: &data.Data{
				Data:     []byte(`test`),
				Metadata: []byte(`{"meta": "meta"}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := createRequest(tt.args.d)
			got, err := GetMultipartDataFromRequest(tt.args.ctx, req, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMultipartDataFromRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMultipartDataFromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
