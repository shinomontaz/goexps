package main

import (
	"context"

	"github.com/gmbyapa/kstream/v2/streams/topology"
)

type Aggregate struct {
	topology.DefaultNode

	maxName string

	store     topology.StateStore
	storeName string
}

// ReadsFrom is used when generating topology diagram
func (a *Aggregate) ReadsFrom() []string {
	return []string{a.storeName}
}

// WritesAt is used when generating topology diagram
func (a *Aggregate) WritesAt() []string {
	return []string{a.storeName}
}

func (a *Aggregate) Build(ctx topology.SubTopologyContext) (topology.Node, error) {
	return &Aggregate{storeName: a.storeName}, nil
}

func (a *Aggregate) Init(ctx topology.NodeContext) error {
	a.store = ctx.Store(a.storeName)
	a.maxName = "max"

	return nil
}

func (a *Aggregate) Type() topology.Type {
	return topology.Type{
		Name: "Custom Aggregator",
	}
}

func (a *Aggregate) Run(ctx context.Context, kIn, vIn interface{}) (kOut, vOut interface{}, cont bool, err error) {
	previous, err := a.store.Get(ctx, a.maxName)
	if err != nil {
		return a.IgnoreWithError(err)
	}

	if previous == nil || previous.(int) < vIn.(int) {
		err := a.store.Set(ctx, a.maxName, vIn, 0)
		if err != nil {
			return a.IgnoreWithError(err)
		}

		// Forward the new max number
		return a.Forward(ctx, kIn, vIn, true)
	}

	// No downstream nodes will be called
	return a.Ignore()
}
