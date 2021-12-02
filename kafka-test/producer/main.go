package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"kafka-test/producer/config"

	"github.com/segmentio/kafka-go"
)

var env *config.Env
var chErrors chan error

func init() {
	rand.Seed(time.Now().UnixNano())
	env = config.NewEnv("./config")
	env.InitKafka()
	env.InitLog()
	//	env.InitDb()

	chErrors = make(chan error, 1000)

	rand.Seed(time.Now().UnixNano())
}

func main() {

	go func() {
		for err := range chErrors {
			fmt.Println("Error", err)
		}
	}()

	chMessages := make(chan []byte, 1000)
	chOk := make(chan struct{})

	var wg, wg2 sync.WaitGroup
	wg.Add(1)
	go generateMessages(chMessages, chOk)

	go func() {
		for range chOk {
			fmt.Println("messages generated")
			wg.Done()
		}
	}()

	wg2.Add(1)
	go func() {
		for mess := range chMessages { // it is safe to do that due to Writer has internal bacth queue
			fmt.Println(string(mess))
			err := env.Kafka.WriteMessages(
				context.Background(),
				kafka.Message{Value: mess},
			)
			if err != nil {
				chErrors <- err
			}
		}
		wg2.Done()
	}()

	fmt.Println("finished 1")
	wg.Wait()

	close(chOk)
	close(chMessages)
	fmt.Println("finished 2")

	wg2.Wait()

	env.Kafka.Close()
	fmt.Println("finished")
}

func generateMessages(chMessages chan<- []byte, chOk chan<- struct{}) {
	defer func() { chOk <- struct{}{} }()

	for i := 0; i < 5; i++ {
		chMessages <- []byte(randStringRunes(10 - i))
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
