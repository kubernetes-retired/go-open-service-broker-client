package v2

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var registerMetrics sync.Once

// Metrics are identified in Prometheus by concatinating Namespace, Subsystem and Name
// ex servicecatalog_osbclient_request_count
const (
	CatalogNamespace = "servicecatalog"
	OSBSubsystem     = "osbclient"
)

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: CatalogNamespace,
			Subsystem: OSBSubsystem,
			Name:      "request_count",
			Help:      "Cumulative number of requests made to the specified Service Broker.",
		},
		[]string{"broker"},
	)
	responses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: CatalogNamespace,
			Subsystem: OSBSubsystem,
			Name:      "response_by_status_count",
			Help:      "Cumulative number of responses from the specified Service Broker grouped by status.",
		},
		[]string{"broker", "status"},
	)
)

func init() {
	registerMetrics.Do(func() {
		prometheus.MustRegister(requests)
		prometheus.MustRegister(responses)
	})
}

func updateMetrics(c *client, response *http.Response, lastRequestError error) {
	requests.WithLabelValues(c.Name).Inc()

	if lastRequestError != nil {
		responses.WithLabelValues(c.Name, "client-error").Inc()
	} else {
		responses.WithLabelValues(c.Name, strconv.Itoa(response.StatusCode/100*100)).Inc()
	}
}
