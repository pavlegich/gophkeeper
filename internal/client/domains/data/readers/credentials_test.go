package readers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	"github.com/pavlegich/gophkeeper/internal/client/domains/user"
)

func TestNewCredentialsReader(t *testing.T) {
	ctx := context.Background()
	type args struct {
		rw rwmanager.RWService
	}
	tests := []struct {
		name string
		args args
		want *CredentialsReader
	}{
		{
			name: "ok",
			args: args{
				rw: nil,
			},
			want: &CredentialsReader{
				data: &user.User{},
				rw:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCredentialsReader(ctx, tt.args.rw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCredentialsReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredentialsReader_Read(t *testing.T) {
	ctx := context.Background()
	type args struct {
		creds *user.User
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
				creds: &user.User{
					Login:    "login",
					Password: "password",
				},
			},
			wantErr: false,
		},
		{
			name: "empty_login",
			args: args{
				creds: &user.User{
					Login:    "",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty_password",
			args: args{
				creds: &user.User{
					Login:    "login",
					Password: "",
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
			in.Write([]byte(fmt.Sprintf("%s\n%s\n", tt.args.creds.Login, tt.args.creds.Password)))

			r := &CredentialsReader{
				data: &user.User{},
				rw:   rw,
			}
			got, err := r.Read(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CredentialsReader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				tt.want, _ = json.MarshalIndent(tt.args.creds, "", "   ")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CredentialsReader.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
