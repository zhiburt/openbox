package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/openbox/monitor/services/qcommunicator"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"

	log "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/openbox/monitor/services/monitor"
	"github.com/openbox/monitor/services/monitor/impl"
	"github.com/openbox/monitor/services/monitor/repositories"
	"github.com/openbox/monitor/services/monitor/transport/endpoints"
	"github.com/openbox/monitor/services/monitor/transport/httptransport"
	"github.com/openbox/monitor/services/qservice"
)

var net = "localhost"
var httpAddr = ":8082"

func init() {
	if n := os.Getenv("NETWORKNAME"); n != "" {
		net = n
	}
	if n := os.Getenv("HTTP_ADDR"); n != "" {
		httpAddr = n
	}
}

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"service", "monitor",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service finished")

	var db *mongo.Client
	{
		var err error
		url := "mongodb://mongo:27017"
		level.Debug(logger).Log("url", url)

		db, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}

		err = db.Ping(context.TODO(), nil)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	var qs qservice.QueueService
	{
		var err error
		qs, err = qservice.NewQueueService("guest", "guest", net, "task_exchange")
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	var qfs qcommunicator.QFileSystem = qcommunicator.NewQFileSystem(logger, qs)

	var mservice monitor.Service
	{
		repository := repositories.NewRepository(db, &logger)
		mservice = impl.NewService(repository, logger, qfs)
	}

	var h http.Handler
	{
		endpoints := endpoints.MakeEndpoints(mservice)
		h = httptransport.NewService(endpoints, logger)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", httpAddr)
		server := &http.Server{
			Addr:    httpAddr,
			Handler: h,
		}
		errs <- server.ListenAndServe()
	}()

	level.Error(logger).Log("exit with", <-errs)
}
