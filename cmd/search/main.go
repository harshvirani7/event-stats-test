package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harshvirani7/event-stats-test/pkg/apis"
	"github.com/harshvirani7/event-stats-test/pkg/cache"
	config "github.com/harshvirani7/event-stats-test/pkg/config"
	"github.com/harshvirani7/event-stats-test/pkg/monitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	apiServerShutdownPeriod = 5 * time.Second
	dbReconnectRetries      = 5
	redisReconnectRetries   = 5
)

var version string

func init() {
	version = "1.0.1"
}

func main() {
	cfg, err := config.Load(
		config.Env(),
		[]string{
			"../../config/search",
			"/usr/src/app/config/search",
			"$GOPATH/src/github.com/eencloud/videosearch-api/config/search",
		},
	)
	exitOn(err)

	// Error handling channel to be passed to all services
	errs := make(chan error)

	logLevel, logger := setupLogger()

	cfgLogLevel := zapcore.Level(cfg.GetInt("log_level"))
	logLevel.SetLevel(cfgLogLevel)
	if cfgLogLevel >= 0 {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	}

	addr := cfg.GetString("redis_addr")
	var rdbClient *cache.Redis

	rdbClient = cache.NewRedis(addr, cfg.GetString("redis_password"), cfg.GetInt("redis_db"), logger, cfg.GetInt("redis_timeout"))

	if rdbClient == nil {
		exitOnNil(rdbClient, "Failed to setup redis connection")
	} else {
		logger.Infof("Redis connection established, addr: %v", addr)
	}

	PromCamRegistry := prometheus.NewRegistry()
	m := monitor.NewMetrics(PromCamRegistry)

	m.Info.With(prometheus.Labels{"version": version}).Set(1)

	// Set up API service routes and controller
	r := gin.Default()
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.GetInt("api.port")),
		Handler: r,
	}

	eventStatsApis := apis.EventStats{
		Logger:    logger,
		RdbClient: rdbClient,
		Cfg:       cfg,
		Metrics:   m,
	}

	r.Use(apis.MonitoringMiddleware(cfg, eventStatsApis))

	collectionEventStats := r.Group(cfg.GetString("api_path") + "eventStats")
	{
		collectionEventStats.POST("/storeEventData", eventStatsApis.StoreEventData())
		collectionEventStats.GET("/totalEventCountByEventType", eventStatsApis.TotalEventCountByType())
		collectionEventStats.GET("/totalEventCountByCameraId", eventStatsApis.TotalEventCountByCameraId())
		collectionEventStats.GET("/eventCountSummaryByCameraId", eventStatsApis.EventCountSummaryByCameraId())
		collectionEventStats.GET("/eventCountSummaryByEventType", eventStatsApis.EventCountSummaryByEventType())
		collectionEventStats.GET("/SummaryByCameraId", eventStatsApis.SummaryByCameraId())
		collectionEventStats.GET("/SummaryByEventType", eventStatsApis.SummaryByEventType())
	}

	r.GET(cfg.GetString("api_path")+"health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Ok",
		})
	})

	r.GET(cfg.GetString("api_path")+"metrics", prometheusHandler(PromCamRegistry))

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.Info("Http server closed")
		}
	}()

	logger.Info("exit", <-errs)

	// The context is used to inform the server it has to finish the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), apiServerShutdownPeriod)
	defer cancel()

	logger.Info("API Service Stop")
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown:", err)
	}

	logger.Info("API Server Exiting")
}

func setupLogger() (zap.AtomicLevel, *zap.SugaredLogger) {
	atom := zap.NewAtomicLevel()
	encoderCfg := zap.NewProductionEncoderConfig()
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	)).WithOptions(zap.AddCaller()).Sugar()
	defer logger.Sync()
	atom.SetLevel(zap.InfoLevel)
	return atom, logger
}

func exitOn(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
}

func exitOnNil(object interface{}, message string) {
	if object == nil {
		fmt.Fprintf(os.Stderr, "%+v", message)
		os.Exit(1)
	}
}

func prometheusHandler(reg *prometheus.Registry) gin.HandlerFunc {
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	return func(c *gin.Context) {
		promHandler.ServeHTTP(c.Writer, c.Request)
	}
}

// Metrics
// no. of events added
// count of total cameras

// total requests made
// average latency
// no of successfull and fail requests

// rate of request for all the endpoints in one graph - number and chart
