// Package twiliosender implements smssender.Sender using twilio programmable sms service.
package twiliosender

import (
	"context"
	"net/http"
	"os"
	"testing"
)

func TestTwilioSMSSender_Send(t *testing.T) {
	account := os.Getenv("TWILIO_ACCOUNT")
	auth := os.Getenv("TWILIO_AUTHTOKEN")
	sender := os.Getenv("TWILIO_SENDER")

	if account == "" || auth == "" || sender == "" {
		t.SkipNow()
	}

	type fields struct {
		account   string
		authToken string
		url       string
		sender    string
		client    Doer
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
				account:   account,
				authToken: auth,
				url:       twilioURL + account + "/Messages.json",
				sender:    sender,
				client:    http.DefaultClient,
			},
			args: args{
				ctx:   context.Background(),
				phone: "+905364474476",
				msg:   "test mesaji",
			},
		},
		{
			name: "fail",
			fields: fields{
				account:   account,
				authToken: auth,
				url:       twilioURL + account + "/Messages.json",
				sender:    sender,
				client:    http.DefaultClient,
			},
			args: args{
				ctx:   context.Background(),
				phone: "+90536447447",
				msg:   "test mesaji",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TwilioSMSSender{
				account:   tt.fields.account,
				authToken: tt.fields.authToken,
				url:       tt.fields.url,
				sender:    tt.fields.sender,
				client:    tt.fields.client,
			}
			if err := s.Send(tt.args.ctx, tt.args.phone, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("TwilioSMSSender.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
