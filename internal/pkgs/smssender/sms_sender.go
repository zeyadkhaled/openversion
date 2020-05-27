// Package smssender provides multiple sms sender for different providers that
// implement Interface in subpackages.
package smssender

import (
	"context"
	"fmt"
	"strings"
)

type Sender interface {
	Send(ctx context.Context, phone, msg string) error
}

type SubSender struct {
	Match  func(phone string) bool
	Sender Sender
}

// CombinedSender takes prefix and sender, if phone starts with prefix use this
// sender, if failed to send continue to next sender. The '+' sign is prefix for
// every phone number use it as accept all.
type CombinedSender struct {
	senders []SubSender
}

func NewCombinedSender(senders ...SubSender) CombinedSender {
	return CombinedSender{senders: senders}
}

func (sender CombinedSender) Send(ctx context.Context, phone, msg string) error {
	var errAccum error
	for _, s := range sender.senders {
		if !s.Match(phone) {
			continue
		}

		err := s.Sender.Send(ctx, phone, msg)
		if err != nil {
			if errAccum != nil {
				errAccum = fmt.Errorf("sender failed(%v)", err)
			} else {
				errAccum = fmt.Errorf("sender failed(%v), %v", err, errAccum)
			}
			continue
		}

		return nil
	}

	if errAccum != nil {
		return errAccum
	}

	return fmt.Errorf("failed to send to phone %q: no matching sender", phone)
}

func PrefixMatcher(prefix string) func(phone string) bool {
	return func(phone string) bool {
		return strings.HasPrefix(phone, prefix)
	}
}
