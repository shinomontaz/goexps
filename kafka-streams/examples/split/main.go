package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/gmbyapa/kstream/v2/kafka"
	"github.com/gmbyapa/kstream/v2/streams"
	"github.com/gmbyapa/kstream/v2/streams/encoding"
)

// установка vcpkg
// установка librdkafka(librdkafka:x64-windows@2.6.0): >vcpkg install librdkafka
// настройка cgo:        >go env -w "CGO_LDFLAGS=-fstack-protector"
// сборка:               >go build

//const TopicNumbers = `example-nums-int`

const TopicNumbers = `example-nums`

func main() {

	config := streams.NewStreamBuilderConfig()
	config.BootstrapServers = []string{"localhost:9921", "localhost:9922", "localhost:9923"}
	config.ApplicationId = `kstream-branching7`
	config.Consumer.Offsets.Initial = kafka.OffsetEarliest

	builder := streams.NewStreamBuilder(config)
	buildTopology(builder)

	topology, err := builder.Build()
	if err != nil {
		panic(err)
	}

	println("Topology - \n", topology.Describe())

	//	panic("!")

	runner := builder.NewRunner()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	go func() {
		<-sigs
		if err := runner.Stop(); err != nil {
			println(err)
		}
	}()

	if err := runner.Run(topology); err != nil {
		panic(err)
	}
}

func buildTopology(builder *streams.StreamBuilder) {
	stream := builder.KStream(TopicNumbers, encoding.NoopEncoder{}, encoding.IntEncoder{})

	splitted := stream.Split()

	splitted.New(`odd`, func(ctx context.Context, key interface{}, val interface{}) (bool, error) {
		return val.(int)%2 != 0, nil
	}).SelectKey(func(ctx context.Context, key, value interface{}) (kOut interface{}, err error) {
		return "odd", nil
	}).Repartition(`odds`, streams.RePartitionWithKeyEncoder(encoding.StringEncoder{})).Aggregate("max_odds", func(ctx context.Context, key, value, previous interface{}) (newAgg interface{}, err error) {
		if previous == nil || previous.(int) < value.(int) {
			return value, nil
		}
		return previous, nil
	}).ToStream().Each(func(ctx context.Context, key, value interface{}) {
		println(`New max odd number:`, value.(int))
	})

	evens := splitted.New(`even`, func(ctx context.Context, key interface{}, val interface{}) (bool, error) {
		return val.(int)%2 == 0, nil
	}).SelectKey(func(ctx context.Context, key, value interface{}) (kOut interface{}, err error) {
		return "even", nil
	}).Repartition(`evens`, streams.RePartitionWithKeyEncoder(encoding.StringEncoder{}))

	evens.Aggregate("max_evens", func(ctx context.Context, key, value, previous interface{}) (newAgg interface{}, err error) {
		if previous == nil || previous.(int) < value.(int) {
			return value, nil
		}
		return previous, nil
	}).Filter(func(ctx context.Context, key, value interface{}) (bool, error) {
		if previous == nil || previous.(int) < value.(int) {
			return value, nil
		}
		return previous, nil
	}).ToStream().Each(func(ctx context.Context, key, value interface{}) {
		println(`New max even number:`, value.(int))
	})
}
