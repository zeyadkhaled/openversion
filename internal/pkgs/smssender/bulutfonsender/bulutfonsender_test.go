package bulutfonsender

import (
	"context"
	"net/http"
	"os"
	"testing"
)

var token = os.Getenv("BULUTFON_SMS_TOKEN")

func TestBulutfon_Send(t *testing.T) {
	if token == "" {
		t.Skip("SMS_TOKEN not set")
	}

	type fields struct {
		messagesURL string
		token       string
		sender      string
		client      Doer
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
			name: "ok",
			fields: fields{
				messagesURL: prodBaseURL + messagesPath,
				token:       token,
				sender:      "HOP.TECH",
				client:      http.DefaultClient,
			},
			args: args{
				ctx:   context.Background(),
				phone: "+905364474476",
				msg:   "Test code: ASDEF",
			},
			wantErr: false,
		},
		{
			name: "sanity",
			fields: fields{
				messagesURL: prodBaseURL + messagesPath,
				token:       "ASD",
				sender:      "HOP.TECH",
				client:      http.DefaultClient,
			},
			args: args{
				ctx:   context.Background(),
				phone: "+905364474476",
				msg:   "Test code: 12345",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bulutfon{
				messagesURL: tt.fields.messagesURL,
				token:       tt.fields.token,
				sender:      tt.fields.sender,
				client:      tt.fields.client,
			}
			if err := b.Send(tt.args.ctx, tt.args.phone, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Bulutfon.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
