package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

var chErrors chan error
var rnd *rand.Rand

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

	chErrors = make(chan error, 1000)
	go handleLog(chErrors)
}

func main() {
	cfg := kafka.WriterConfig{
		Brokers: []string{
			"localhost:9921", "localhost:9922", "localhost:9923",
		},
		Topic: "orders",
	}
	k := kafka.NewWriter(cfg)

	var (
		jo  []byte
		err error
	)
	ctx := context.Background()
	for {
		time.Sleep(time.Duration(rnd.Intn(10)) * time.Second)
		o := newOrder()
		jo, err = json.Marshal(o)
		if err != nil {
			chErrors <- err
			continue
		}
		err := k.WriteMessages(
			ctx,
			kafka.Message{
				Value: jo,
				Key:   []byte(o.Number),
			},
		)
		if err != nil {
			chErrors <- err
			continue
		}

	}
}

func handleLog(ch <-chan error) {
	for e := range ch {
		fmt.Println(e)
	}
}
