package usecases

import "time"

type Member struct {
	Group           string
	ID              string
	LastUpdatedTime *time.Time
	Partition
}

type Partition struct {
	Ordinal int64
	Total   int64
}
