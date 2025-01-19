package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var sleep = false

func main() {

	flag.BoolVar(&sleep, "is_sleep", false, "should we slow down consuming")
	flag.Parse()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9921, localhost:9922, localhost:9923",
		"group.id":          "kafka-go-getting-started6",
		"auto.offset.reset": "earliest",

		"enable.auto.offset.store": false,
		"max.poll.interval.ms":     20000, // session.timeout.ms по умолчанию = 10 secs
		"session.timeout.ms":       10000,
	})

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	topic := "example-nums"
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}

				if sleep {
					time.Sleep(120 * time.Second)
				}

				// _, err := c.StoreMessage(e)
				// if err != nil {
				// 	fmt.Fprintf(os.Stderr, "%% Error storing offset after message %s: %v\n", e.TopicPartition, err)
				// }
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			}
		}
	}

	c.Close()
}
