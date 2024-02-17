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

type Cre func(ctx context.Context, ordinal, total int64)
