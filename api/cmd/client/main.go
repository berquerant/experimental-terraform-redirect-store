package main

import (
	"context"
	"encoding/json"
	"errors"
	"experimental-terraform-redirect-store/api"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

const usage = `api-client

Usage:
  api-client status
  api-client scan
  api-client get NAME
  api-client put NAME TO
  api-clinet delete NAME

Flags:`

func Usage() {
	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
}

func main() {
	var (
		endpoint = flag.String("endpoint", "http://127.0.0.1:8030", "")
	)
	flag.Usage = Usage
	flag.Parse()

	client := api.NewClientImpl(
		*endpoint,
		&http.Client{
			Timeout: 3 * time.Second,
		},
	)

	args := flag.Args()
	r, err := send(context.Background(), client, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", b)
}

var (
	ErrInvalidArgument = errors.New("InvalidArgument")
)

func send(ctx context.Context, c api.Client, args []string) (any, error) {
	if len(args) == 0 {
		return nil, ErrInvalidArgument
	}
	switch args[0] {
	case "status":
		err := c.Status(ctx)
		return nil, err
	case "scan":
		return c.Scan(ctx)
	case "get":
		if len(args) < 2 {
			return nil, ErrInvalidArgument
		}
		return c.Get(ctx, args[1])
	case "put":
		if len(args) < 3 {
			return nil, ErrInvalidArgument
		}
		return c.Put(ctx, &api.Record{
			Name: args[1],
			To:   args[2],
		})
	case "delete":
		if len(args) < 2 {
			return nil, ErrInvalidArgument
		}
		err := c.Delete(ctx, args[1])
		return nil, err
	default:
		return nil, fmt.Errorf("%w, unknown command %s", ErrInvalidArgument, args[0])
	}
}
