package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var (
	_ Client = NewClientImpl("", nil)
)

type Client interface {
	Status(ctx context.Context) error
	Scan(ctx context.Context) ([]*Record, error)
	Get(ctx context.Context, name string) (*Record, error)
	Put(ctx context.Context, record *Record) (*Record, error)
	Delete(ctx context.Context, name string) error
}

func NewClientImpl(endpoint string, client *http.Client) *ClientImpl {
	return &ClientImpl{
		endpoint: endpoint,
		client:   client,
	}
}

type ClientImpl struct {
	endpoint string
	client   *http.Client
}

func (c *ClientImpl) api(pattern string) string {
	return fmt.Sprintf("%s%s", c.endpoint, pattern)
}

func (c *ClientImpl) Status(ctx context.Context) error {
	_, err := c.client.Get(c.api("/status"))
	return err
}

func (c *ClientImpl) Scan(ctx context.Context) ([]*Record, error) {
	r, err := Post[ScanRequest, ScanResponse](c.client, c.api("/scan"))(ctx, ScanRequest{})
	if err != nil {
		return nil, err
	}
	if r.Error != "" {
		return nil, errors.New(r.Error)
	}
	return r.Records, nil
}

func (c *ClientImpl) Get(ctx context.Context, name string) (*Record, error) {
	r, err := Post[GetRequest, GetResponse](c.client, c.api("/get"))(ctx, GetRequest{
		Name: name,
	})
	if err != nil {
		return nil, fmt.Errorf("%w, %s", err, name)
	}
	if r.Error != "" {
		return nil, fmt.Errorf("%s, %s", r.Error, name)
	}
	return r.Record, nil
}

func (c *ClientImpl) Put(ctx context.Context, record *Record) (*Record, error) {
	r, err := Post[PutRequest, PutResponse](c.client, c.api("/put"))(ctx, PutRequest{
		Record: record,
	})
	if err != nil {
		return nil, fmt.Errorf("%w, %v", err, record)
	}
	if r.Error != "" {
		return nil, fmt.Errorf("%s, %v", r.Error, record)
	}
	return r.Record, nil
}

func (c *ClientImpl) Delete(ctx context.Context, name string) error {
	r, err := Post[DeleteRequest, DeleteResponse](c.client, c.api("/delete"))(ctx, DeleteRequest{
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("%w, %s", err, name)
	}
	if r.Error != "" {
		return fmt.Errorf("%s, %s", r.Error, name)
	}
	return nil
}
