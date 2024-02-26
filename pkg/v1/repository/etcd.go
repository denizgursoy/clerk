package repository

import (
	"fmt"
	"time"

	"github.com/denizgursoy/clerk/pkg/v1/config"
	"go.etcd.io/etcd/client/v3"
)

func CreateETCDClient(cfg config.Config) (*clientv3.Client, error) {
	// Set up an etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{cfg.ETCDEndpoint},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create etc client: %w", err)
	}

	return client, nil
}

func Stop(client *clientv3.Client) error {
	if err := client.Close(); err != nil {
		return fmt.Errorf("could not close the etcd client: %w", err)
	}

	return nil
}
