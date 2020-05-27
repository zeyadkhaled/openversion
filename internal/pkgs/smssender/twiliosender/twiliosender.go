// Package twiliosender implements smssender.Sender using twilio programmable sms service.
package twiliosender

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const twilioURL = "https://api.twilio.com/2010-04-01/Accounts/"

type TwilioSMSSender struct {
	account   string
	authToken string

	url    string
	sender string
	client Doer
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// New create sms sender for twilio.
func New(accountSID string,
	authToken string,
	sender string, client Doer) TwilioSMSSender {

	if client == nil {
		c := *http.DefaultClient
		c.Timeout = 2 * time.Second
		client = &c
	}

	return TwilioSMSSender{
		account:   accountSID,
		authToken: authToken,

		url:    twilioURL + accountSID + "/Messages.json",
		sender: sender,
		client: client,
	}
}

func (s TwilioSMSSender) Send(ctx context.Context, phone, msg string) error {
	v := url.Values{}
	v.Set("From", s.sender)
	v.Set("To", phone)
	v.Set("Body", msg)
	// fail messages with delivery price higher than 1$
	v.Set("MaxPrice", "1")

	req, err := http.NewRequest(http.MethodPost, s.url, strings.NewReader(v.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create twilio sms request: %v", err)
	}
	req.SetBasicAuth(s.account, s.authToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %v", err)
	}
	defer resp.Body.Close()

	var r createResp
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return fmt.Errorf("failed to parse twilio sms, response status(%s): %v", resp.Status, err)
	}
	if r.ErrorCode != nil || r.ErrorMessage != nil {
		return fmt.Errorf("send sms error: %v, %v", r.ErrorCode, r.ErrorMessage)
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("send sms failed with non OK resp: status code %d", resp.StatusCode)
	}

	return nil
}

type createResp struct {
	Status       string
	ErrorCode    *int
	ErrorMessage *string
}

/* example response
{
    "sid": "SM1e9f29bb48554b13aa92632eb1b0282b",
    "date_created": "Fri, 07 Feb 2020 18:44:48 +0000",
    "date_updated": "Fri, 07 Feb 2020 18:44:48 +0000",
    "date_sent": null,
    "account_sid": "AC679ed094bf1188d719001563f0ee3936",
    "to": "+905364474476",
    "from": "+14159410721",
    "messaging_service_sid": null,
    "body": "Sent from your Twilio trial account - dene asd --=?",
    "status": "queued",
    "num_segments": "1",
    "num_media": "0",
    "direction": "outbound-api",
    "api_version": "2010-04-01",
    "price": null,
    "price_unit": "USD",
    "error_code": null,
    "error_message": null,
    "uri": "/2010-04-01/Accounts/AC679ed094bf1188d719001563f0ee3936/Messages/SM1e9f29bb48554b13aa92632eb1b0282b.json",
    "subresource_uris": {
        "media": "/2010-04-01/Accounts/AC679ed094bf1188d719001563f0ee3936/Messages/SM1e9f29bb48554b13aa92632eb1b0282b/Media.json"
    }
}
*/
