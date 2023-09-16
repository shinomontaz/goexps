package main

import (
	"context"
	"fmt"
	"math/rand"
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

	var mess []byte
	for {
		mess = []byte(randStringRunes(10))
		fmt.Println(mess)
		err := env.Kafka.WriteMessages(
			context.Background(),
			kafka.Message{Value: mess},
		)
		if err != nil {
			chErrors <- err
		}

	}
	fmt.Println("finished")
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
