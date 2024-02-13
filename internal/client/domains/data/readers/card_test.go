package readers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/pavlegich/gophkeeper/internal/client/domains/rwmanager"
)

func TestNewCardReader(t *testing.T) {
	ctx := context.Background()
	type args struct {
		rw rwmanager.RWService
	}
	tests := []struct {
		name string
		args args
		want *CardReader
	}{
		{
			name: "ok",
			args: args{
				rw: nil,
			},
			want: &CardReader{
				details: &CardDetails{},
				rw:      nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCardReader(ctx, tt.args.rw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCardReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardReader_Read(t *testing.T) {
	ctx := context.Background()
	type args struct {
		card *CardDetails
		exp  string
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
				card: &CardDetails{
					Number: 5536913798031973,
					Owner:  "Card Holder",
					CV:     123,
				},
				exp: "08/34",
			},
			wantErr: false,
		},
		{
			name: "invalid_card_number",
			args: args{
				card: &CardDetails{
					Number: 8476442824861248,
					Owner:  "Card Holder",
					CV:     123,
				},
				exp: "08/34",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid_card_exp_date",
			args: args{
				card: &CardDetails{
					Number: 5536913798031973,
					Owner:  "Card Holder",
					CV:     123,
				},
				exp: "32/34",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "card_owner_empty",
			args: args{
				card: &CardDetails{
					Number: 5536913798031973,
					Owner:  "",
					CV:     123,
				},
				exp: "02/34",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid_card_cv",
			args: args{
				card: &CardDetails{
					Number: 5536913798031973,
					Owner:  "Card Holder",
					CV:     1232,
				},
				exp: "02/34",
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
			in.Write([]byte(fmt.Sprintf("%d\n%s\n%s\n%d\n", tt.args.card.Number, tt.args.exp,
				tt.args.card.Owner, tt.args.card.CV)))

			r := &CardReader{
				details: &CardDetails{},
				rw:      rw,
			}
			got, err := r.Read(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CardReader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				tt.args.card.Expires, _ = time.Parse("01/06", tt.args.exp)
				tt.want, _ = json.MarshalIndent(tt.args.card, "", "   ")
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CardReader.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
