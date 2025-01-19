package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/gmbyapa/kstream/v2/kafka"
	"github.com/gmbyapa/kstream/v2/streams"
	"github.com/gmbyapa/kstream/v2/streams/encoding"
)

const TopicNumbers = `example-nums`
const oddsStore = `odds-store5`
const evensStore = `evens-store5`

func main() {

	config := streams.NewStreamBuilderConfig()
	config.BootstrapServers = []string{"localhost:9921", "localhost:9922", "localhost:9923"}
	config.ApplicationId = `kstream-branching13`
	config.Consumer.Offsets.Initial = kafka.OffsetEarliest

	builder := streams.NewStreamBuilder(config)
	buildTopology(builder)

	topology, err := builder.Build()
	if err != nil {
		panic(err)
	}

	//	println("Topology - \n", topology.Describe())

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

	odds := splitted.New(`odd`, func(ctx context.Context, key interface{}, val interface{}) (bool, error) {
		return val.(int)%2 != 0, nil
	}).SelectKey(func(ctx context.Context, key, value interface{}) (kOut interface{}, err error) {
		return "odd", nil
	}).Repartition(`odds`, streams.RePartitionWithKeyEncoder(encoding.StringEncoder{}))
	odds.AddStateStore(oddsStore, encoding.StringEncoder{}, encoding.IntEncoder{})

	odds.NewProcessor(&Aggregate{storeName: oddsStore}).Each(func(ctx context.Context, key, value interface{}) {
		println(`Odd number:`, value.(int))
	}).Each(func(ctx context.Context, key, value interface{}) {
		println(`New max odd number:`, value.(int))
	})

	evens := splitted.New(`even`, func(ctx context.Context, key interface{}, val interface{}) (bool, error) {
		return val.(int)%2 == 0, nil
	}).SelectKey(func(ctx context.Context, key, value interface{}) (kOut interface{}, err error) {
		return "even", nil
	}).Repartition(`evens`, streams.RePartitionWithKeyEncoder(encoding.StringEncoder{}))
	evens.AddStateStore(evensStore, encoding.StringEncoder{}, encoding.IntEncoder{})

	evens.NewProcessor(&Aggregate{storeName: evensStore}).Each(func(ctx context.Context, key, value interface{}) {
		println(`Even number:`, value.(int))
	}).Each(func(ctx context.Context, key, value interface{}) {
		println(`New max even number:`, value.(int))
	})
}
