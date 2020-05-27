// Package stdoutsender provides smssender which prints sms to stdout.
package stdoutsender

import (
	"context"
	"fmt"
)

type Client struct{}

func (c Client) Send(ctx context.Context, phone, msg string) error {
	fmt.Printf("SMS to %q: %q\n", phone, msg)
	return nil
}
