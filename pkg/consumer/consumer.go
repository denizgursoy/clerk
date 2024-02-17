package consumer

import (
	"fmt"

	"github.com/denizgursoy/clerk/internal/redis"
	"github.com/denizgursoy/clerk/internal/usecases"
)

type Consumer struct {
	group string
	c     usecases.Cache
	index int
}

func NewConsumer(group string) (Consumer, error) {
	consumer := Consumer{group: group, c: redis.NewRedisCache()}
	index, err := consumer.c.AddNewInstance()
	if err != nil {
		return Consumer{}, fmt.Errorf("could not initialize consumer: %w", err)
	}
	consumer.index = index

	return consumer, nil
}

func (c Consumer) Consume(d Action) {

}

type Action func(index, allPartition int)
