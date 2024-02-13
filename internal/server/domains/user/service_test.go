package user

import (
	"context"
	"reflect"
	"testing"
)

func TestNewUserService(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx  context.Context
		repo Repository
	}
	tests := []struct {
		name string
		args args
		want *UserService
	}{
		{
			name: "ok",
			args: args{
				ctx:  ctx,
				repo: nil,
			},
			want: &UserService{
				repo: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(tt.args.ctx, tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}
