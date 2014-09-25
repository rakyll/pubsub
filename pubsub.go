// Package pubsub is a Google Cloud Pub/Sub client.
package pubsub

import (
	"encoding/base64"
	"errors"
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

	open chan bool
}

type Topic struct {
	proj string
	name string
	s    *raw.Service
}

type Message struct {
	Data   []byte
	Labels map[string]interface{}
}

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
		open: make(chan bool),
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
	return s.s.Subscriptions.Delete(fullSubName(s.proj, s.name)).Do()
}

func (s *Subscription) ModifyAckDeadline(deadline time.Duration) error {
	return s.s.Subscriptions.ModifyAckDeadline(&raw.ModifyAckDeadlineRequest{
		Subscription:       fullSubName(s.proj, s.name),
		AckDeadlineSeconds: int64(deadline),
	}).Do()
}

func (s *Subscription) ModifyPushEndpoint(endpoint string) error {
	return s.s.Subscriptions.ModifyPushConfig(&raw.ModifyPushConfigRequest{
		Subscription: fullSubName(s.proj, s.name),
		PushConfig: &raw.PushConfig{
			PushEndpoint: endpoint,
		},
	}).Do()
}

func (s *Subscription) IsExists() (bool, error) {
	panic("not yet implemented")
}

func (s *Subscription) Ack(id ...string) error {
	return s.s.Subscriptions.Acknowledge(&raw.AcknowledgeRequest{
		Subscription: fullSubName(s.proj, s.name),
		AckId:        id,
	}).Do()
}

func (s *Subscription) pull() ([]*Message, error) {
	panic("not yet implemented")
}

func (s *Subscription) Listen() <-chan *Message {
	messages := make(chan *Message)
	go func() {
		select {
		case <-s.open:
			return
		default:
			msgs, err := s.pull()
			if err != nil {
				panic("not yet implemented")
				return
			}
			for _, m := range msgs {
				messages <- m
			}
		}
	}()
	return messages
}

func (s *Subscription) Stop() {
	close(s.open)
}

func (c *Client) Topic(name string) *Topic {
	return &Topic{
		proj: c.proj,
		name: name,
		s:    c.s,
	}
}

func (t *Topic) Create() error {
	_, err := t.s.Topics.Create(&raw.Topic{
		Name: fullTopicName(t.proj, t.name),
	}).Do()
	return err
}

func (t *Topic) Delete() error {
	return t.s.Topics.Delete(fullTopicName(t.proj, t.name)).Do()
}

func (t *Topic) IsExists() (bool, error) {
	panic("not yet implemented")
}

func (t *Topic) Publish(msg *Message) error {
	var labels []*raw.Label
	if msg.Labels != nil {
		labels := []*raw.Label{}
		for k, v := range msg.Labels {
			l := &raw.Label{Key: k}
			switch v.(type) {
			case int64:
				l.NumValue = v.(int64)
			case string:
				l.StrValue = v.(string)
			default:
				return errors.New("pubsub: label value could be either an int64 or a string")
			}
			labels = append(labels, l)
		}
	}
	return t.s.Topics.Publish(&raw.PublishRequest{
		Topic: fullTopicName(t.proj, t.name),
		Message: &raw.PubsubMessage{
			Data:  base64.StdEncoding.EncodeToString(msg.Data), // base64 encoded
			Label: labels,
		},
	}).Do()
}

func fullSubName(proj, name string) string {
	return fmt.Sprintf("/subscriptions/%s/%s", proj, name)
}

func fullTopicName(proj, name string) string {
	return fmt.Sprintf("/topics/%s/%s", proj, name)
}
