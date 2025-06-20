package store

import (
	"context"
	"log"

	qgrpc "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	Qdrant qgrpc.QdrantClient
}

func New(ctx context.Context, endpoint string) *Client {
	conn, err := grpc.DialContext(ctx, endpoint,
		grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("connect qdrant: %v", err)
	}
	return &Client{
		conn:   conn,
		Qdrant: qgrpc.NewQdrantClient(conn),
	}
}

func (c *Client) Close() error { return c.conn.Close() }
