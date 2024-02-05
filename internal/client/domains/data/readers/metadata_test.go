package readers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
)

func TestNewMetadataReader(t *testing.T) {
	ctx := context.Background()
	type args struct {
		rw rwmanager.RWService
	}
	tests := []struct {
		name string
		args args
		want *MetadataReader
	}{
		{
			name: "ok",
			args: args{
				rw: nil,
			},
			want: &MetadataReader{
				data: map[string]string{},
				rw:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMetadataReader(ctx, tt.args.rw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetadataReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadataReader_Read(t *testing.T) {
	ctx := context.Background()
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				in: "meta : data\nclose",
			},
			wantErr: false,
		},
		{
			name: "empty_input",
			args: args{
				in: "",
			},
			wantErr: false,
		},
		{
			name: "without_value",
			args: args{
				in: "key : ",
			},
			wantErr: true,
		},
		{
			name: "many_dots",
			args: args{
				in: "f : f : f",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var in bytes.Buffer
			var out bytes.Buffer
			rw := rwmanager.NewRWManager(context.Background(), &in, &out)
			in.Write([]byte(fmt.Sprintf("%s\n", tt.args.in)))

			r := &MetadataReader{
				data: map[string]string{},
				rw:   rw,
			}
			got, err := r.Read(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetadataReader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.want, _ = json.MarshalIndent(r.data, "", "   ")
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetadataReader.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
