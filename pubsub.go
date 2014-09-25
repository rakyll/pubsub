// Package pubsub is a Google Cloud Pub/Sub client.
package pubsub

import (
	"fmt"
	"net/http"
	"time"

	raw "code.google.com/p/google-api-go-client/pubsub/v1beta1"
)

type Client struct {
	proj string
	s    *raw.Service
}

type Subscription struct {
	proj string
	name string
	s    *raw.Service
}

type Topic struct {
	proj string
	name string
	s    *raw.Service
}

type Message struct{}

func New(projID string, tr http.RoundTripper) *Client {
	return NewWithClient(projID, &http.Client{Transport: tr})
}

func NewWithClient(projID string, c *http.Client) *Client {
	// TODO(jbd): Add user-agent.
	s, _ := raw.New(c)
	return &Client{proj: projID, s: s}
}

func (c *Client) Subscription(name string) *Subscription {
	return &Subscription{
		proj: c.proj,
		name: name,
		s:    c.s,
	}
}

func (s *Subscription) Create(topic string, deadline time.Duration, endpoint string) error {
	sub := &raw.Subscription{
		Topic: fullTopicName(s.proj, topic),
		Name:  fullSubName(s.proj, s.name),
	}
	if int64(deadline) > 0 {
		sub.AckDeadlineSeconds = int64(deadline)
	}
	if endpoint != "" {
		sub.PushConfig = &raw.PushConfig{PushEndpoint: endpoint}
	}
	_, err := s.s.Subscriptions.Create(sub).Do()
	return err
}

func (s *Subscription) Delete() error {
	panic("not yet implemented")
}

func (s *Subscription) ModifyAckDeadline(deadline time.Duration) error {
	panic("not yet implemented")
}

func (s *Subscription) ModifyPushEndpoint(endpoint string) error {
	panic("not yet implemented")
}

func (s *Subscription) IsExists() (bool, error) {
	panic("not yet implemented")
}

func (s *Subscription) Ack(id ...string) error {
	panic("not yet implemented")
}

func (s *Subscription) Listen() (<-chan *Message, error) {
	panic("not yet implemented")
}

func (c *Client) Topic(name string) *Topic {
	return &Topic{
		proj: c.proj,
		name: name,
		s:    c.s,
	}
}

func (t *Topic) Create() error {
	panic("not yet implemented")
}

func (t *Topic) Delete() error {
	panic("not yet implemented")
}

func (t *Topic) IsExists() (bool, error) {
	panic("not yet implemented")
}

func (t *Topic) Publish(msg *Message) error {
	panic("not yet implemented")
}

func fullSubName(proj, sub string) string {
	return fmt.Sprintf("/subscriptions/%s/%s", proj, sub)
}

func fullTopicName(proj, topic string) string {
	return fmt.Sprintf("/topics/%s/%s", proj, topic)
}
