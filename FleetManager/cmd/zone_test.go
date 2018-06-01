package main

import (
	"fmt"
	"testing"

	"github.com/shinomontaz/goexps/FleetManager/model"
)

func TestGetZones(t *testing.T) {
	zones := model.GetZones()
	fmt.Println(zones)
}

func TestGetFleet(t *testing.T) {
	zones := model.GetFleet()
	fmt.Println(zones)
}
