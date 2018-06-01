package main

import (
	"encoding/json"
	"fmt"
	log2 "log"
	"os"

	"github.com/shinomontaz/goexps/FleetManager/FleetManager"
	"github.com/shinomontaz/goexps/FleetManager/model"
	"gitlab.ozon.ru/platform/scratch"
	"gitlab.ozon.ru/platform/tracer/log"
)

var appVersion, goVersion, buildTimestamp string

func main() {

	app, err := scratch.New(
		scratch.WithAppInfo("FleetManager", appVersion),
		scratch.WithBuildInfo(goVersion, buildTimestamp),
		scratch.WithPorts(4000, 4001, 4002),
		scratch.WithLogLevel(log.DEBUG),
		scratch.WithImpl(FleetManager.NewImpl(&FleetManager.ServerImpl{})),
	)

	if err != nil {
		log2.Fatalf("can't create app: %s", err)
	}
	app.Run()
}

type Config struct {
	DbHost   string
	DbUser   string
	DbPass   string
	DbName   string
	DbPort   int
	TestFlag bool
}

var cfg Config

func init() {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&cfg)
	if err != nil {
		log2.Fatal(err)
	}

	model.InitDb(initDsn())
}

func initDsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName)
}
