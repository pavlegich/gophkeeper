package utils

import (
	"errors"
	"net/http"
	"testing"

	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
)

func TestCheckStatusCode(t *testing.T) {
	type args struct {
		code int
	}
	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				code: http.StatusOK,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "bad_request_status",
			args: args{
				code: http.StatusBadRequest,
			},
			want:    errs.ErrBadRequest,
			wantErr: true,
		},
		{
			name: "unauthorized_status",
			args: args{
				code: http.StatusUnauthorized,
			},
			want:    errs.ErrUnauthorized,
			wantErr: true,
		},
		{
			name: "conflict_status",
			args: args{
				code: http.StatusConflict,
			},
			want:    errs.ErrAlreadyExists,
			wantErr: true,
		},
		{
			name: "internal_server_status",
			args: args{
				code: http.StatusInternalServerError,
			},
			want:    errs.ErrServerInternal,
			wantErr: true,
		},
		{
			name: "no_content_status",
			args: args{
				code: http.StatusNoContent,
			},
			want:    errs.ErrNotExist,
			wantErr: true,
		},
		{
			name: "unknown_status",
			args: args{
				code: http.StatusNotFound,
			},
			want:    errs.ErrUnknownStatusCode,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckStatusCode(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckStatusCode() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !errors.Is(err, tt.want) {
				t.Errorf("CheckStatusCode() = %v, want %v", err, tt.want)
			}
		})
	}
}
