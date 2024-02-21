package usecases

import "time"

type Member struct {
	Group           string
	ID              string
	LastUpdatedTime *time.Time
}
