package readers

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
)

func TestNewTextReader(t *testing.T) {
	ctx := context.Background()
	type args struct {
		rw rwmanager.RWService
	}
	tests := []struct {
		name string
		args args
		want *TextReader
	}{
		{
			name: "ok",
			args: args{
				rw: nil,
			},
			want: &TextReader{
				text: "",
				rw:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTextReader(ctx, tt.args.rw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTextReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextReader_Read(t *testing.T) {
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
				in: "text text text",
			},
			want:    []byte("text text text"),
			wantErr: false,
		},
		{
			name: "empty_input",
			args: args{
				in: "",
			},
			want:    []byte(""),
			wantErr: false,
		},
		{
			name: "several_lines",
			args: args{
				in: "line\nline\nline line 3",
			},
			want:    []byte("linelineline line 3"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var in bytes.Buffer
			var out bytes.Buffer
			rw := rwmanager.NewRWManager(context.Background(), &in, &out)
			in.Write([]byte(fmt.Sprintf("%s\nclose\n", tt.args.in)))

			r := &TextReader{
				text: "",
				rw:   rw,
			}
			got, err := r.Read(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TextReader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TextReader.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
