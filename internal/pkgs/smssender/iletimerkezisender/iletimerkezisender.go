// Package iletimerkezisender sends sms through iletimetkezi.
package iletimerkezisender

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const iletiMerkeziURL = "https://api.iletimerkezi.com/v1/send-sms"

type IletiMerkezi struct {
	auth   iletiUserAuth
	sender string
	client Doer
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// New create sms sender for iletimerkezi.com
func New(username string,
	password string,
	sender string, client Doer) IletiMerkezi {

	if client == nil {
		c := *http.DefaultClient
		c.Timeout = 2 * time.Second
		client = &c
	}

	return IletiMerkezi{
		auth: iletiUserAuth{
			Username: username,
			Password: password,
		},
		sender: sender,
		client: client,
	}
}

func (i IletiMerkezi) Send(ctx context.Context, phone, msg string) error {
	obj := iletiRequest{Auth: i.auth,
		Order: iletiOrder{Sender: i.sender,
			Messages: []iletiMessage{
				{
					Text:       iletiTXT{Text: msg},
					Receipents: []string{phone},
				},
			},
		},
	}
	b, err := xml.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal iletimerkezi request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, iletiMerkeziURL, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to create iletimerkezi request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "text/xml")

	htmlResp, err := i.client.Do(req)
	if err != nil {
		return err
	}
	defer htmlResp.Body.Close()

	var resp iletiResponse
	err = xml.NewDecoder(htmlResp.Body).Decode(&resp)
	if err != nil {
		return fmt.Errorf("failed to decode response from iletimerkezi: %w", err)
	}

	if resp.Status.StatusCode == 200 {
		return nil
	}
	return errors.New(resp.Status.Message)
}

type iletiUserAuth struct {
	XMLName  xml.Name `xml:"authentication"`
	Username string   `xml:"username"`
	Password string   `xml:"password"`
}

type iletiTXT struct {
	Text string `xml:",cdata"`
}

type iletiMessage struct {
	XMLName    xml.Name `xml:"message"`
	Text       iletiTXT `xml:"text"`
	Receipents []string `xml:"receipents>number"`
}

type iletiOrder struct {
	XMLName      xml.Name `xml:"order"`
	Sender       string   `xml:"sender"`
	SendDateTime string   `xml:"sendDateTime,omitempty"`
	Messages     []iletiMessage
}

type iletiRequest struct {
	XMLName xml.Name `xml:"request"`
	Auth    iletiUserAuth
	Order   iletiOrder
}

type iletiStatus struct {
	StatusCode int    `xml:"code"`
	Message    string `xml:"message"`
}

type iletiResponse struct {
	XMLName xml.Name    `xml:"response"`
	Status  iletiStatus `xml:"status"`
	OrderID int         `xml:"order>id"`
}
