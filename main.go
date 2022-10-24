package main

import (
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	api2 "luismatosgarcia.dev/video-sharing-go/internal/api"
	"luismatosgarcia.dev/video-sharing-go/internal/pkg/datastore"
	"luismatosgarcia.dev/video-sharing-go/internal/pkg/jsonlog"
	"luismatosgarcia.dev/video-sharing-go/internal/server/http"
	"luismatosgarcia.dev/video-sharing-go/internal/videos"
	"os"
	"strings"
)

var (
	version = "0.0.1"
)

func main() {
	var httpConfig http.Config

	// Environment flags ---------------------------------------------------------------------------

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

	flag.StringVar(&httpConfig.FileStore.AwsAccessKeyId, "filestore-access-key-id", "123", "S3 Bucket Key ID")
	flag.StringVar(&httpConfig.FileStore.AwsSecretKey, "filestore-secret-key", "xyz", "S3 Bucket Secret Key")
	flag.StringVar(&httpConfig.FileStore.AwsBucketName, "filestore-bucket-name", "video-sharing-app-bucket", "S3 Bucket Name")
	flag.StringVar(&httpConfig.FileStore.AwsRegion, "filestore-region", "eu-east-1", "S3 Region")
	flag.StringVar(&httpConfig.FileStore.AwsEndpoint, "filestore-endpoint", "http://localhost:4566", "S3 Endpoint")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version\t%s\n", version)
	}

	// Dependencies -------------------------------------------------------------------------------

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := datastore.NewService(&httpConfig)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	//TODO: Initialize FIlestore and pass it to videoserviec

	// Services ------------------------------------------------------------------------------------
	videoService, err := videos.NewService(db)
	if err != nil {
		logger.PrintFatal(err, nil)
		return
	}

	// API -----------------------------------------------------------------------------------------
	api, err := api2.NewService(logger, videoService)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	// Server
	h, err := http.NewService(&httpConfig, api)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	h.Start()

}
