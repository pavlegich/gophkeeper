// Package repository contains repository object
// and methods for interaction between service and storage.
package repository

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pavlegich/gophkeeper/internal/server/domains/data"
	"github.com/pavlegich/gophkeeper/internal/server/mocks"
)

func TestNewDataRepository(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
		db  *sql.DB
	}
	tests := []struct {
		name string
		args args
		want *Repository
	}{
		{
			name: "ok",
			args: args{
				ctx: ctx,
				db:  nil,
			},
			want: &Repository{
				db: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDataRepository(tt.args.ctx, tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDataRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetDataByName(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockDataRepository(ctrl)

	gomock.InOrder(
		mock.EXPECT().GetDataByName(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&data.Data{Name: "myCreds", Type: "credentials"}, nil),
	)

	type args struct {
		ctx   context.Context
		dType string
		name  string
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
				ctx:   ctx,
				dType: "credentials",
				name:  "myCreds",
			},
			want: &data.Data{
				Name: "myCreds",
				Type: "credentials",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mock.GetDataByName(tt.args.ctx, tt.args.dType, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetDataByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetDataByName() = %v, want %v", got, tt.want)
			}
		})
	}
}
