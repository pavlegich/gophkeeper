package data

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"reflect"
	"testing"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	"github.com/pavlegich/gophkeeper/internal/common/infra/config"
)

func TestNewDataService(t *testing.T) {
	ctx := context.Background()
	type args struct {
		rw  rwmanager.RWService
		cfg *config.ClientConfig
	}
	tests := []struct {
		name string
		args args
		want *DataService
	}{
		{
			name: "ok",
			args: args{
				rw:  nil,
				cfg: nil,
			},
			want: &DataService{
				rw:  nil,
				cfg: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDataService(ctx, tt.args.rw, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDataService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readDataTypeAndName(t *testing.T) {
	ctx := context.Background()

	type args struct {
		in *Data
	}
	tests := []struct {
		name    string
		args    args
		want    *Data
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				in: &Data{
					Type: "credentials",
					Name: "myCreds",
				},
			},
			want: &Data{
				Type: "credentials",
				Name: "myCreds",
			},
			wantErr: false,
		},
		{
			name: "ivalid_type",
			args: args{
				in: &Data{
					Type: "wrong_type",
					Name: "anyName",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty_name",
			args: args{
				in: &Data{
					Type: "credentials",
					Name: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty_type",
			args: args{
				in: &Data{
					Type: "",
					Name: "anyName",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var in bytes.Buffer
			var out bytes.Buffer
			rw := rwmanager.NewRWManager(context.Background(), &in, &out)

			in.Write([]byte(fmt.Sprintf("%s\n%s\n", tt.args.in.Type, tt.args.in.Name)))
			got, err := readDataTypeAndName(ctx, rw)
			if (err != nil) != tt.wantErr {
				t.Errorf("readDataTypeAndName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readDataTypeAndName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createMultipartData(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	defer mw.Close()

	type args struct {
		mpwriter *multipart.Writer
		d        *Data
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "text_ok",
			args: args{
				mpwriter: mw,
				d: &Data{
					Name:     "myText",
					Type:     "text",
					Data:     []byte(`text data`),
					Metadata: []byte(`{"meta": "data"}`),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var in bytes.Buffer
			var out bytes.Buffer
			rw := rwmanager.NewRWManager(context.Background(), &in, &out)

			in.Write([]byte(fmt.Sprintf("%s\n%s\n%s\nclose\n", tt.args.d.Type, tt.args.d.Name, tt.args.d.Data)))
			if err := createMultipartData(ctx, rw, tt.args.mpwriter, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("createMultipartData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
