package utils

import (
	"testing"
)

func TestIsValidDataType(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok",
			args: args{
				t: "text",
			},
			want: true,
		},
		{
			name: "invalid_value",
			args: args{
				t: "not_text",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidDataType(tt.args.t); got != tt.want {
				t.Errorf("IsValidDataType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidCardNumber(t *testing.T) {
	type args struct {
		number int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok",
			args: args{
				number: 5536913798031973,
			},
			want: true,
		},
		{
			name: "short_number",
			args: args{
				number: 5536913798031,
			},
			want: false,
		},
		{
			name: "invalid_number",
			args: args{
				number: 1234567812345678,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidCardNumber(tt.args.number); got != tt.want {
				t.Errorf("IsValidCardNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checksum(t *testing.T) {
	type args struct {
		number int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "check",
			args: args{
				number: 553691379803197,
			},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checksum(tt.args.number); got != tt.want {
				t.Errorf("checksum() = %v, want %v", got, tt.want)
			}
		})
	}
}
