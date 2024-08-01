package main

import (
	"flag"
	"os"
)

var ipAddr string = "localhost" //адрес сервера
var FlagServerPort string       //адрес сервера и порта var FlagBaseURL string
var FlagDatabaseDSN string      //наименование базы данных
var FlagLevelLogger string

func ParseFlags() {

	defaultRunAddr := ipAddr + ":8080"
	defaultDatabaseDSN := "" //"user=postgres password=postgres host=localhost port=5432 dbname=gophermart sslmode=disable"
	defaultLevelLogger := "local"

	flag.StringVar(&FlagServerPort, "a", defaultRunAddr, "address and port to run server")
	flag.StringVar(&FlagDatabaseDSN, "d", defaultDatabaseDSN, "name database Postgres")
	flag.StringVar(&FlagLevelLogger, "l", defaultLevelLogger, "log level")

	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		FlagServerPort = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		FlagLevelLogger = envLogLevel
	}
	envDatabaseDSN, ok := os.LookupEnv("DATABASE_URI")
	if ok && envDatabaseDSN != "" {
		FlagDatabaseDSN = envDatabaseDSN
	}
}
