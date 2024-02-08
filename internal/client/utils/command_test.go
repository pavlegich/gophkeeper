package utils

import (
	"bytes"
	"context"
	"errors"
	"syscall"
	"testing"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
)

func TestDoWithRetryIfEmpty(t *testing.T) {
	var in bytes.Buffer
	var out bytes.Buffer
	rw := rwmanager.NewRWManager(context.Background(), &in, &out)

	ctx := context.Background()
	type args struct {
		rw rwmanager.RWService
		f  func(ctx context.Context) error
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
				rw: rw,
				f: func(ctx context.Context) error {
					return nil
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "empty_input_error",
			args: args{
				rw: rw,
				f: func(ctx context.Context) error {
					return errs.ErrEmptyInput
				},
			},
			want:    errs.ErrEmptyInput,
			wantErr: true,
		},
		{
			name: "not_empty_input_and_nil_error",
			args: args{
				rw: rw,
				f: func(ctx context.Context) error {
					return errs.ErrBadRequest
				},
			},
			want:    errs.ErrBadRequest,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DoWithRetryIfEmpty(ctx, tt.args.rw, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoWithRetryIfEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !errors.Is(err, tt.want) {
				t.Errorf("DoWithRetryIfEmpty() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestGetKnownErr(t *testing.T) {
	type args struct {
		err error
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
				err: nil,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "bad_request_error",
			args: args{
				err: errs.ErrBadRequest,
			},
			want:    errs.ErrBadRequest,
			wantErr: true,
		},
		{
			name: "already_exists_error",
			args: args{
				err: errs.ErrAlreadyExists,
			},
			want:    errs.ErrAlreadyExists,
			wantErr: true,
		},
		{
			name: "server_internal_error",
			args: args{
				err: errs.ErrServerInternal,
			},
			want:    errs.ErrServerInternal,
			wantErr: true,
		},
		{
			name: "not_exist_error",
			args: args{
				err: errs.ErrNotExist,
			},
			want:    errs.ErrNotExist,
			wantErr: true,
		},
		{
			name: "empty_input_error",
			args: args{
				err: errs.ErrEmptyInput,
			},
			want:    errs.ErrEmptyInput,
			wantErr: true,
		},
		{
			name: "unauthorized_error",
			args: args{
				err: errs.ErrUnauthorized,
			},
			want:    errs.ErrUnauthorized,
			wantErr: true,
		},
		{
			name: "connection_refused_error",
			args: args{
				err: syscall.ECONNREFUSED,
			},
			want:    errs.ErrConnectionRefused,
			wantErr: true,
		},
		{
			name: "invalid_card_number_error",
			args: args{
				err: errs.ErrInvalidCardNumber,
			},
			want:    errs.ErrInvalidCardNumber,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GetKnownErr(tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKnownErr() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !errors.Is(err, tt.want) {
				t.Errorf("GetKnownErr() = %v, wantErr %v", err, tt.want)
			}
		})
	}
}
