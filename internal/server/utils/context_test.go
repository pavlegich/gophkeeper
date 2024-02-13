// Package utils contains additional methods for server.
package utils

import (
	"context"
	"testing"
)

func TestGetUserIDFromContext(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx: context.WithValue(ctx, ContextIDKey, 1),
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "no_value",
			args: args{
				ctx: ctx,
			},
			want:    -1,
			wantErr: true,
		},
		{
			name: "invalid_value_type",
			args: args{
				ctx: context.WithValue(ctx, ContextIDKey, "1"),
			},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserIDFromContext(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserIDFromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserIDFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
