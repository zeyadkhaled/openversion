// Package bulutfonsender sends sms through iletimetkezi.
package bulutfonsender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const prodBaseURL = "https://api.bulutfon.com"
const messagesPath = "/v2/sms/messages"

type Bulutfon struct {
	messagesURL string
	token       string
	sender      string
	client      Doer
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// New create sms sender for iletimerkezi.com
func New(baseURL, token, sender string, client Doer) Bulutfon {

	if baseURL == "" {
		baseURL = prodBaseURL
	} else {
		baseURL = strings.TrimSuffix(baseURL, "/")
	}

	if client == nil {
		c := *http.DefaultClient
		c.Timeout = 2 * time.Second
		client = &c
	}

	return Bulutfon{
		messagesURL: baseURL + messagesPath,
		token:       token,
		sender:      sender,
		client:      client,
	}
}

func (b Bulutfon) Send(ctx context.Context, phone, msg string) error {

	phone = strings.TrimPrefix(phone, "+")
	body, err := json.Marshal(messageRequest{
		Title:      b.sender,
		Content:    msg,
		Receivers:  []string{phone},
		RejectLink: false,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal sms req: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, b.messagesURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create bulutfon sms request: %v", err)
	}
	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("apikey", b.token)
	req.URL.RawQuery = q.Encode()

	httpResp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do bulutfon sms request: %v", err)
	}
	defer httpResp.Body.Close()

	switch {
	case httpResp.StatusCode/100 == 2:
		return nil
	default:
		var errMsg bulutfonErr
		err := json.NewDecoder(httpResp.Body).Decode(&errMsg)
		if err != nil {
			return fmt.Errorf("failed to parse bulutfon sms %s resp: %v", httpResp.Status, err)
		}
		return errMsg.Error
	}
}

type messageRequest struct {
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	Receivers  []string   `json:"receivers"`
	RejectLink bool       `json:"reject_link"`
	SendDate   *time.Time `json:"send_date"`
}

type bulutfonErr struct {
	Error errResp `json:"error"`
}

type errResp struct {
	Code    int                 `json:"code"`
	Title   string              `json:"title"`
	Message string              `json:"message"`
	Details map[string][]string `json:"details"`
}

func (err errResp) Error() string {
	return fmt.Sprintf("failed: code %d, title %q, message %q, details: %v", err.Code, err.Title, err.Message, err.Details)
}
