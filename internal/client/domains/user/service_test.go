package user

import (
	"context"
	"reflect"
	"testing"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
	"github.com/pavlegich/gophkeeper/internal/common/infra/config"
)

func TestNewUserService(t *testing.T) {
	ctx := context.Background()
	type args struct {
		rw  rwmanager.RWService
		cfg *config.ClientConfig
	}
	tests := []struct {
		name string
		args args
		want *UserService
	}{
		{
			name: "ok",
			args: args{
				rw:  nil,
				cfg: nil,
			},
			want: &UserService{
				rw:  nil,
				cfg: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(ctx, tt.args.rw, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}
