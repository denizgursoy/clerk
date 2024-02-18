package client

import "context"

type ClerkServerConfig struct {
	Address                    string
	KeepAliveDurationInSeconds int64
}

type Member struct {
	Group string
	ID    string
}

type NotifyFunction func(ctx context.Context, ordinal, total int64) error
