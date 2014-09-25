// Package pubsub is a Google Cloud Pub/Sub client.
package pubsub

import (
	"net/http"
	"time"
)

type Client struct {
	c http.RoundTripper
}

type Subscription struct {
	name string
	c    http.RoundTripper
}

type Topic struct {
	name string
	c    http.RoundTripper
}

type Message struct {
}

func New(tr http.RoundTripper) (*Client, error) {
	panic("not yet implemented")
}

func NewWithClient(c *http.Client) (*Client, error) {
	panic("not yet implemented")
}

func (c *Client) Subscription(name) *Subscription {
	return &Subscription{name: name, c: c.c}
}

func (s *Subscription) Create() error {
	panic("not yet implemented")
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
	panic("not yet implemented")
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
