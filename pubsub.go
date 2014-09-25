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
}

type Topic struct {
	proj string
	name string
	s    *raw.Service
}

type Message struct {
	AckID  string
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

func (s *Subscription) Pull(retImmediately bool) (*Message, error) {
	resp, err := s.s.Subscriptions.Pull(&raw.PullRequest{
		Subscription:      fullSubName(s.proj, s.name),
		ReturnImmediately: retImmediately,
	}).Do()
	if err != nil {
		return nil, err
	}
	data, err := base64.StdEncoding.DecodeString(resp.PubsubEvent.Message.Data)
	if err != nil {
		return nil, err
	}

	labels := make(map[string]interface{})
	for _, l := range resp.PubsubEvent.Message.Label {
		if l.StrValue != "" {
			labels[l.Key] = l.StrValue
		} else {
			labels[l.Key] = l.NumValue
		}
	}
	return &Message{
		AckID:  resp.AckId,
		Data:   data,
		Labels: labels,
	}, nil
}

// TODO(jbd): Add (*Subscription).Listen and (*Subscription).Stop

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

func (t *Topic) Publish(data []byte, labels map[string]interface{}) error {
	var rawLabels []*raw.Label
	if labels != nil {
		rawLabels := []*raw.Label{}
		for k, v := range labels {
			l := &raw.Label{Key: k}
			switch v.(type) {
			case int64:
				l.NumValue = v.(int64)
			case string:
				l.StrValue = v.(string)
			default:
				return errors.New("pubsub: label value could be either an int64 or a string")
			}
			rawLabels = append(rawLabels, l)
		}
	}
	return t.s.Topics.Publish(&raw.PublishRequest{
		Topic: fullTopicName(t.proj, t.name),
		Message: &raw.PubsubMessage{
			Data:  base64.StdEncoding.EncodeToString(data),
			Label: rawLabels,
		},
	}).Do()
}

func fullSubName(proj, name string) string {
	return fmt.Sprintf("/subscriptions/%s/%s", proj, name)
}

func fullTopicName(proj, name string) string {
	return fmt.Sprintf("/topics/%s/%s", proj, name)
}
