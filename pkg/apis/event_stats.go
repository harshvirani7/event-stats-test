package apis

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harshvirani7/event-stats-test/model"
	"github.com/harshvirani7/event-stats-test/pkg/cache"
	"github.com/harshvirani7/event-stats-test/pkg/config"
	"github.com/harshvirani7/event-stats-test/pkg/monitor"
	"github.com/harshvirani7/event-stats-test/pkg/storage"
	"github.com/harshvirani7/event-stats-test/utils"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type RequestDetailCount struct {
	SuccessRequestCount                    int `json:"successRequestCount"`
	ErrorRequestCount                      int `json:"errorRequestCount"`
	StoreEventDataSuccessCnt               int `json:"storeEventDataSuccessCnt"`
	TotalEventCountByTypeSuccessCnt        int `json:"totalEventCountByTypeSuccessCnt"`
	TotalEventCountByCameraIdSuccessCnt    int `json:"totalEventCountByCameraIdSuccessCnt"`
	EventCountSummaryByCameraIdSuccessCnt  int `json:"eventCountSummaryByCameraIdSuccessCnt"`
	EventCountSummaryByEventTypeSuccessCnt int `json:"eventCountSummaryByEventTypeSuccessCnt"`
	SummaryByCameraIdSuccessCnt            int `json:"summaryByCameraIdSuccessCnt"`
	SummaryByEventTypeSuccessCnt           int `json:"summaryByEventTypeSuccessCnt"`
}

var ReqDetailCount RequestDetailCount

type EventStats struct {
	Logger    *zap.SugaredLogger
	RdbClient *cache.Redis
	Cfg       config.Config
	Metrics   *monitor.Metrics
}

func init() {

	ReqDetailCount = RequestDetailCount{
		SuccessRequestCount:                    0,
		ErrorRequestCount:                      0,
		StoreEventDataSuccessCnt:               0,
		TotalEventCountByTypeSuccessCnt:        0,
		TotalEventCountByCameraIdSuccessCnt:    0,
		EventCountSummaryByEventTypeSuccessCnt: 0,
		EventCountSummaryByCameraIdSuccessCnt:  0,
		SummaryByEventTypeSuccessCnt:           0,
		SummaryByCameraIdSuccessCnt:            0,
	}
}

func (es EventStats) StoreEventData() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var data []model.Data
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		sd := storage.StoreDataAPI{
			RdbClient: es.RdbClient,
		}
		// Store events into Redis
		// change this to call from a struct
		if err := sd.StoreEventData(data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			ReqDetailCount.ErrorRequestCount = ReqDetailCount.ErrorRequestCount + 1
			es.Metrics.ErrorRequest.Set(float64(ReqDetailCount.ErrorRequestCount))

			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Event data stored successfully"})

		ReqDetailCount.SuccessRequestCount += +1
		es.Metrics.SuccessRequest.Set(float64(ReqDetailCount.SuccessRequestCount))

		ReqDetailCount.StoreEventDataSuccessCnt += 1
		es.Metrics.StoreEventDataSuccess.Set(float64(ReqDetailCount.StoreEventDataSuccessCnt))

		count, err := storage.GetTotalEventCount(es.RdbClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			ReqDetailCount.ErrorRequestCount = ReqDetailCount.ErrorRequestCount + 1
			es.Metrics.ErrorRequest.Set(float64(ReqDetailCount.ErrorRequestCount))

			return
		}
		es.Metrics.EventCount.Set(float64(count))
	}
	return gin.HandlerFunc(fn)
}

func (es EventStats) TotalEventCountByType() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		now := time.Now()
		// Get event type from query parameter
		eventType := c.Query("eventType")
		if eventType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "eventType query parameter is required"})

			ReqDetailCount.ErrorRequestCount = ReqDetailCount.ErrorRequestCount + 1
			es.Metrics.ErrorRequest.Set(float64(ReqDetailCount.ErrorRequestCount))

			return
		}

		// Get total event count for the given event type
		count, err := storage.GetTotalEventCountByType(eventType, es.RdbClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			ReqDetailCount.ErrorRequestCount = ReqDetailCount.ErrorRequestCount + 1
			es.Metrics.ErrorRequest.Set(float64(ReqDetailCount.ErrorRequestCount))

			return
		}

		utils.Sleep(200)

		es.Metrics.DurationCountByEventType.With(prometheus.Labels{"method": "GET", "status": "200"}).Observe(time.Since(now).Seconds())

		// Return total count and eventType in the response
		response := gin.H{
			"eventType":         eventType,
			"total_event_count": count,
		}
		c.JSON(http.StatusOK, response)

		ReqDetailCount.SuccessRequestCount = ReqDetailCount.SuccessRequestCount + 1
		es.Metrics.SuccessRequest.Set(float64(ReqDetailCount.SuccessRequestCount))

		ReqDetailCount.TotalEventCountByTypeSuccessCnt += 1
		es.Metrics.TotalEventCountByTypeSuccess.Set(float64(ReqDetailCount.TotalEventCountByTypeSuccessCnt))
	}
	return gin.HandlerFunc(fn)
}

func (es EventStats) TotalEventCountByCameraId() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		now := time.Now()
		// Get eventType from query parameter
		cameraId := c.Query("cameraId")
		if cameraId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cameraId query parameter is required"})

			ReqDetailCount.ErrorRequestCount = ReqDetailCount.ErrorRequestCount + 1
			es.Metrics.ErrorRequest.Set(float64(ReqDetailCount.ErrorRequestCount))

			return
		}

		// Get total event count for each cameraid for the given eventType
		count, err := storage.GetEventCountByCameraID(cameraId, es.RdbClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			ReqDetailCount.ErrorRequestCount = ReqDetailCount.ErrorRequestCount + 1
			es.Metrics.ErrorRequest.Set(float64(ReqDetailCount.ErrorRequestCount))

			return
		}

		utils.Sleep(200)

		es.Metrics.DurationCountByCameraId.With(prometheus.Labels{"method": "GET", "status": "200"}).Observe(time.Since(now).Seconds())

		// Return total count and eventType in the response
		response := gin.H{
			"cameraId":          cameraId,
			"total_event_count": count,
		}
		c.JSON(http.StatusOK, response)

		ReqDetailCount.SuccessRequestCount = ReqDetailCount.SuccessRequestCount + 1
		es.Metrics.SuccessRequest.Set(float64(ReqDetailCount.SuccessRequestCount))

		ReqDetailCount.TotalEventCountByCameraIdSuccessCnt += 1
		es.Metrics.TotalEventCountByCameraIdSuccess.Set(float64(ReqDetailCount.TotalEventCountByCameraIdSuccessCnt))
	}
	return gin.HandlerFunc(fn)
}

func (es EventStats) EventCountSummaryByCameraId() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Get cameraId from query parameter
		cameraId := c.Query("cameraId")
		if cameraId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cameraId query parameter is required"})

			ReqDetailCount.ErrorRequestCount = ReqDetailCount.ErrorRequestCount + 1
			es.Metrics.ErrorRequest.Set(float64(ReqDetailCount.ErrorRequestCount))

			return
		}

		// Get event counts for the given cameraId from Redis
		eventCounts, err := storage.GetEventCountSummaryByCameraID(cameraId, es.RdbClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			ReqDetailCount.ErrorRequestCount = ReqDetailCount.ErrorRequestCount + 1
			es.Metrics.ErrorRequest.Set(float64(ReqDetailCount.ErrorRequestCount))

			return
		}

		// Return each event counts for the given cameraId
		response := gin.H{
			"cameraId":      cameraId,
			"event_summary": eventCounts,
		}
		c.JSON(http.StatusOK, response)

		ReqDetailCount.SuccessRequestCount = ReqDetailCount.SuccessRequestCount + 1
		es.Metrics.SuccessRequest.Set(float64(ReqDetailCount.SuccessRequestCount))

		ReqDetailCount.EventCountSummaryByCameraIdSuccessCnt += 1
		es.Metrics.EventCountSummaryByCameraIdSuccess.Set(float64(ReqDetailCount.EventCountSummaryByCameraIdSuccessCnt))
	}
	return gin.HandlerFunc(fn)
}

func (es EventStats) EventCountSummaryByEventType() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Get eventType from query parameter
		eventType := c.Query("eventType")
		if eventType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "eventType query parameter is required"})

			ReqDetailCount.ErrorRequestCount = ReqDetailCount.ErrorRequestCount + 1
			es.Metrics.ErrorRequest.Set(float64(ReqDetailCount.ErrorRequestCount))

			return
		}

		// Get total camera event counts for the given eventType from Redis
		eventCounts, err := storage.GetEventCountSummaryByEventType(eventType, es.RdbClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			ReqDetailCount.ErrorRequestCount = ReqDetailCount.ErrorRequestCount + 1
			es.Metrics.ErrorRequest.Set(float64(ReqDetailCount.ErrorRequestCount))

			return
		}

		// Return each event counts for the given cameraId
		response := gin.H{
			"eventType":      eventType,
			"camera_summary": eventCounts,
		}
		c.JSON(http.StatusOK, response)

		ReqDetailCount.SuccessRequestCount = ReqDetailCount.SuccessRequestCount + 1
		es.Metrics.SuccessRequest.Set(float64(ReqDetailCount.SuccessRequestCount))

		ReqDetailCount.EventCountSummaryByEventTypeSuccessCnt += 1
		es.Metrics.EventCountSummaryByEventTypeSuccess.Set(float64(ReqDetailCount.EventCountSummaryByEventTypeSuccessCnt))
	}
	return gin.HandlerFunc(fn)
}

func (es EventStats) SummaryByCameraId() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Get eventType from query parameter
		cameraId := c.Query("cameraId")
		if cameraId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "eventType query parameter is required"})

			ReqDetailCount.ErrorRequestCount = ReqDetailCount.ErrorRequestCount + 1
			es.Metrics.ErrorRequest.Set(float64(ReqDetailCount.ErrorRequestCount))

			return
		}

		// Get event summary for the given eventType from Redis
		eventSummary, err := storage.GetEventSummaryByCameraID(cameraId, es.RdbClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			ReqDetailCount.ErrorRequestCount = ReqDetailCount.ErrorRequestCount + 1
			es.Metrics.ErrorRequest.Set(float64(ReqDetailCount.ErrorRequestCount))

			return
		}

		// Return event summary for the given eventType
		c.JSON(http.StatusOK, gin.H{"eventType": cameraId, "event_summary": eventSummary})

		ReqDetailCount.SuccessRequestCount = ReqDetailCount.SuccessRequestCount + 1
		es.Metrics.SuccessRequest.Set(float64(ReqDetailCount.SuccessRequestCount))

		ReqDetailCount.SummaryByCameraIdSuccessCnt += 1
		es.Metrics.SummaryByCameraIdSuccess.Set(float64(ReqDetailCount.SummaryByCameraIdSuccessCnt))
	}
	return gin.HandlerFunc(fn)
}

func (es EventStats) SummaryByEventType() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Get eventType from query parameter
		eventType := c.Query("eventType")
		if eventType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "eventType query parameter is required"})

			ReqDetailCount.ErrorRequestCount = ReqDetailCount.ErrorRequestCount + 1
			es.Metrics.ErrorRequest.Set(float64(ReqDetailCount.ErrorRequestCount))

			return
		}

		// Get event summary for the given eventType from Redis
		eventSummary, err := storage.GetEventSummaryByEventType(eventType, es.RdbClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			ReqDetailCount.ErrorRequestCount = ReqDetailCount.ErrorRequestCount + 1
			es.Metrics.ErrorRequest.Set(float64(ReqDetailCount.ErrorRequestCount))

			return
		}

		// Return event summary for the given eventType
		c.JSON(http.StatusOK, gin.H{"eventType": eventType, "event_summary": eventSummary})

		ReqDetailCount.SuccessRequestCount = ReqDetailCount.SuccessRequestCount + 1
		es.Metrics.SuccessRequest.Set(float64(ReqDetailCount.SuccessRequestCount))

		ReqDetailCount.SummaryByEventTypeSuccessCnt += 1
		es.Metrics.SummaryByEventTypeSuccess.Set(float64(ReqDetailCount.SummaryByEventTypeSuccessCnt))
	}
	return gin.HandlerFunc(fn)
}
