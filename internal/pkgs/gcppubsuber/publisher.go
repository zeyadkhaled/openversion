// Package gcppubsuber provides publisher and subscriber that sends/recieves
// messages from gcp.
package gcppubsuber

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Publisher struct {
	topic *pubsub.Topic
}

func NewPublisher(serviceAccountJSON []byte, projectID, topicName string) (*Publisher, error) {

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

	return &Publisher{
		topic: topic,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, data []byte) (id string, err error) {
	res := p.topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	return res.Get(ctx)
}

func (p *Publisher) PublishAsync(ctx context.Context, data []byte) {
	_ = p.topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})
}

func (p *Publisher) Close() error {
	p.topic.Stop()
	return nil
}

func client(ctx context.Context, serviceAccountJSON []byte, projectID string) (*pubsub.Client, error) {
	var opts []option.ClientOption
	if serviceAccountJSON != nil {

		creds, err := google.CredentialsFromJSON(ctx, serviceAccountJSON, pubsub.ScopePubSub)
		if err != nil {
			return nil, fmt.Errorf("failed to create credential from json: %w", err)
		}

		if projectID == "" {
			projectID = creds.ProjectID
		}

		opts = append(opts, option.WithCredentials(creds))
	} else {
		var err error
		projectID, err = metadata.ProjectID()
		if err != nil {
			return nil, fmt.Errorf("failed to get project id from gcp metadata server: %v", err)
		}
	}

	client, err := pubsub.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}
