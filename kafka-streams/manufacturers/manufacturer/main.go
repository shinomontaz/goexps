package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	brokers := []string{"localhost:9921", "localhost:9922", "localhost:9923"}
	kreader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  "consumer-group-demands",
		Topic:    "demands",
		MaxBytes: 10e6, // 10MB
	})

	kwriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   "produced",
	})

	if kreader == nil {
		log.Fatal("cannot create reader")
	}
	if kwriter == nil {
		log.Fatal("cannot create writer")
	}

	var (
		m   kafka.Message
		p   Produced
		d   Demand
		jp  []byte
		err error
	)
	ctx := context.Background()
	for {
		m, err = kreader.FetchMessage(ctx)
		if err != nil {
			chErrors <- err
			break
		}

		err = json.Unmarshal(m.Value, &d)
		if err != nil {
			break
		}

		time.Sleep(time.Duration(rnd.Intn(20)) * time.Second)
		p = Produced{
			Number:     d.OrderNum,
			Product:    d.Id,
			Created_at: time.Now(),
		}

		jp, err = json.Marshal(p)
		if err != nil {
			chErrors <- err
			continue
		}

		err = kwriter.WriteMessages(
			ctx,
			kafka.Message{
				Value: jp,
				Key:   []byte(p.Number),
			},
		)
		if err != nil {
			chErrors <- err
		}

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

type Produced struct {
	Number     string
	Product    int
	Created_at time.Time
}
