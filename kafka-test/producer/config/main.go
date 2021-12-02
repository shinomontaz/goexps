package config

import (
	"fmt"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	_ "net/http/pprof"

	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
)

type Config struct {
	ListenPort int  `env:"CHUPD_LISTENPORT"`
	TestFlag   bool `env:"CHUPD_TESTFLAG"`

	DbHost string `env:"CHUPD_DBHOST"`
	DbName string `env:"CHUPD_DBNAME"`
	DbUser string `env:"CHUPD_DBUSER"`
	DbPass string `env:"CHUPD_DBPASS"`
	DbPort int    `env:"CHUPD_DBPORT"`

	KfConfig kafka.WriterConfig `json:"KfConfig"`
}

type Env struct {
	Db       *sqlx.DB
	Config   *Config
	loglevel log.Level
	Kafka    *kafka.Writer
}

func NewEnv(path string) *Env {
	viper.SetConfigType("json")
	viper.SetConfigName("conf")
	viper.AddConfigPath(".")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		checkErr(err)
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		checkErr(err)
	}

	loglevel := log.WarnLevel

	return &Env{
		Config:   &cfg,
		loglevel: loglevel,
	}
}

func (e *Env) InitLog() {
	if e.Config.TestFlag {
		e.loglevel = log.DebugLevel
	}

	log.SetLevel(e.loglevel)
	log.SetFormatter(&log.JSONFormatter{})
}

func (e *Env) InitDb() {
	dsn := initDbDsn(e.Config)
	log.Debug(dsn)
	db, err := sqlx.Connect("postgres", dsn)
	checkErr(err)
	e.Db = db
}

func initDbDsn(cfg *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName)
}

func (e *Env) InitKafka() {
	e.Kafka = kafka.NewWriter(e.Config.KfConfig)
}

func checkErr(err error) {
	if err != nil {
		_, filename, lineno, ok := runtime.Caller(1)
		message := ""
		if ok {
			message = fmt.Sprintf("%v:%v: %v\n", filename, lineno, err)
		}
		log.Panic(message, err)
	}
}
