package iletimerkezisender

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
)

var user = os.Getenv("ILETIMERKEZI_USER")
var pass = os.Getenv("ILETIMERKEZI_PASSWORD")
var sender = os.Getenv("ILETIMERKEZI_SENDER")

func TestIletiMerkezi_Send(t *testing.T) {
	if user == "" || pass == "" || sender == "" {
		t.SkipNow()
	}

	type fields struct {
		auth   iletiUserAuth
		sender string
		client Doer
	}
	type args struct {
		ctx   context.Context
		phone string
		msg   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "SMS test",
			fields: fields{
				auth: iletiUserAuth{
					Username: user,
					Password: pass,
				},
				sender: sender,
				client: http.DefaultClient,
			},
			args: args{
				ctx:   context.Background(),
				phone: "+905364474476",
				// Random int is used to prevent "Tekrar eden sipariş"
				msg: fmt.Sprintf("Test code: %d", rand.Int()), //nolint:gosec
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := IletiMerkezi{
				auth:   tt.fields.auth,
				sender: tt.fields.sender,
				client: tt.fields.client,
			}
			if err := i.Send(tt.args.ctx, tt.args.phone, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("IletiMerkezi.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
