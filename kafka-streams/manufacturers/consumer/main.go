package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
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
	brokers := []string{"localhost:9921", "localhost:9922", "localhost:9923"}
	kreader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  "consumer-group-orders",
		Topic:    "orders",
		MaxBytes: 10e6, // 10MB
	})

	kwriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   "demands",
	})

	if kreader == nil {
		log.Fatal("cannot create reader")
	}
	if kwriter == nil {
		log.Fatal("cannot create writer")
	}

	var (
		o  Order
		d  Demand
		jd []byte
	)
	ctx := context.Background()
	for {
		m, err := kreader.FetchMessage(ctx)
		if err != nil {
			chErrors <- err
			break
		}

		err = json.Unmarshal(m.Value, &o)
		if err != nil {
			break
		}

		for _, it := range o.Items {
			d = Demand{
				Id:       it,
				OrderNum: o.Number,
			}
			jd, err = json.Marshal(d)
			if err != nil {
				chErrors <- err
				continue
			}

			err := kwriter.WriteMessages(
				ctx,
				kafka.Message{
					Value: jd,
					Key:   []byte(strconv.Itoa(it)),
				},
			)

			if err != nil {
				if err != nil {
					chErrors <- err
					break
				}
			}
		}

		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		if err := kreader.CommitMessages(ctx, m); err != nil {
			log.Fatal("failed to commit messages from input topic:", err)
		}
	}

	if err := kreader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	if err := kwriter.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func handleLog(ch <-chan error) {
	for e := range ch {
		fmt.Println(e)
	}
}

type Demand struct {
	Id       int
	OrderNum string
}

type Order struct {
	Number     string
	Items      []int
	Created_at time.Time
}
