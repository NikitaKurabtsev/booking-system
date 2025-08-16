package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_system_requests_total",
			Help: "Total number of HTTP requests by app_name, method and path.",
		},
		[]string{"method", "path"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "booking_system_requests_duration_seconds",
			Help:    "HTTP request latencies in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_system_responses_total",
			Help: "Total number of HTTP responses by status code, app_name and path.",
		},
		[]string{"status_code", "path"},
	)

	httpRequestsInProgress = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "booking_system_requests_in_progress",
			Help: "Number of HTTP requests currently being processed.",
		},
		[]string{"path"},
	)

	appExceptions = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_system_exceptions_total",
			Help: "Total number of exceptions or panics per application.",
		},
		[]string{},
	)
)

func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpResponses)
	prometheus.MustRegister(httpRequestsInProgress)
	prometheus.MustRegister(appExceptions)
}

// TODO: Authorization header - Bearer <token>
func (h *Handler) verifyUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeader)
		if authHeader == "" {
			h.errorResponse(c, http.StatusUnauthorized, "empty authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			h.errorResponse(c, http.StatusUnauthorized, "empty authorization header")
			return
		}

		if len(parts[1]) == 0 {
			h.errorResponse(c, http.StatusUnauthorized, "token is empty")
			return
		}

		token := parts[1]
		userID, err := h.service.User.ParseToken(token)
		if err != nil {
			h.errorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.Set(userCtx, userID)
	}
}

func (h *Handler) prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		method := c.Request.Method

		httpRequestsInProgress.WithLabelValues(path).Inc()
		defer httpRequestsInProgress.WithLabelValues(path).Dec()

		c.Next()

		statusCode := strconv.Itoa(c.Writer.Status())

		// Счётчик ответов по кодам
		httpResponses.WithLabelValues(statusCode, path).Inc()

		// Время выполнения
		elapsed := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(method, path).Observe(elapsed)

		// Если это ошибка — увеличиваем счётчик исключений
		if statusCode[0] == '5' || statusCode == "422" {
			appExceptions.WithLabelValues().Inc()
		}
	}
}
