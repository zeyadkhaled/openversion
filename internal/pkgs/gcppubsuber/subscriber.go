// Package gcppubsuber provides publisher and subscriber that sends/recieves
// messages from gcp.
package gcppubsuber

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Subscription struct {
	sub *pubsub.Subscription
}

func NewSubscription(serviceAccountJSON []byte, projectID, topicName, subscriptionName string) (*Subscription, error) {

	ctx := context.Background()
	c, err := client(ctx, serviceAccountJSON, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %v", err)
	}

	topic, err := c.CreateTopic(ctx, topicName)
	if err != nil {
		if status.Code(err) != codes.AlreadyExists {
			return nil, err
		}

		topic = c.Topic(topicName)
		ok, err := topic.Exists(ctx)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New("topic not found")
		}
	}

	sub, err := c.CreateSubscription(context.Background(), subscriptionName, pubsub.SubscriptionConfig{Topic: topic})
	if err != nil {
		if status.Code(err) != codes.AlreadyExists {
			return nil, err
		}

		sub = c.Subscription(subscriptionName)
		ok, err := topic.Exists(ctx)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New("topic not found")
		}
	}

	sub.ReceiveSettings.MaxExtension = 24 * time.Hour

	return &Subscription{
		sub: sub,
	}, nil
}

type Message struct {
	in *pubsub.Message
}

func (m *Message) Ack()         { m.in.Ack() }
func (m *Message) Nack()        { m.in.Nack() }
func (m *Message) Data() []byte { return m.in.Data }

func (s *Subscription) Subscribe(ctx context.Context, f func(sctx context.Context, msg *Message)) error {
	return s.sub.Receive(ctx,
		func(ctx context.Context, m *pubsub.Message) {
			f(ctx, &Message{in: m})
		},
	)
}
