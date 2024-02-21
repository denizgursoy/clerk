package usecases

import "time"

type Member struct {
	Group           string
	ID              string
	LastUpdatedTime *time.Time
	Ordinal         int64
	Total           int64
}
