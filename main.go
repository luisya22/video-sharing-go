package main

import (
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	api2 "luismatosgarcia.dev/video-sharing-go/internal/api"
	"luismatosgarcia.dev/video-sharing-go/internal/pkg/datastore"
	"luismatosgarcia.dev/video-sharing-go/internal/pkg/jsonlog"
	"luismatosgarcia.dev/video-sharing-go/internal/server/http"
	"os"
	"strings"
)

var (
	version = "0.0.1"
)

func main() {
	var httpConfig http.Config

	flag.IntVar(&httpConfig.Port, "port", 4000, "API server port")
	flag.StringVar(&httpConfig.Env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&httpConfig.Db.Dsn, "db-dsn", "", "PostgreSQL DSN")
	flag.IntVar(&httpConfig.Db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&httpConfig.Db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&httpConfig.Db.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&httpConfig.Limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&httpConfig.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&httpConfig.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		httpConfig.Cors.TrustedOrigins = strings.Fields(val)
		return nil
	})

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version\t%s\n", version)
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := datastore.NewService(&httpConfig)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	api, err := api2.NewService(logger)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	h, err := http.NewService(&httpConfig, api)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	h.Start()

}
