package agent

import (
	"runtime"
	"time"

	"github.com/juju/loggo"
)

//
// Pulls together the bridge, a cached state configuration and the bus.
//
type Agent struct {
	conf     *Config
	bridge   *Bridge
	memstats *runtime.MemStats
	metrics  *MetricService
	eventCh  chan statusEvent
	log      loggo.Logger
}

func createAgent(conf *Config) *Agent {
	return &Agent{
		conf:     conf,
		bridge:   createBridge(conf),
		memstats: &runtime.MemStats{},
		metrics:  CreateMetricService(),
		log:      loggo.GetLogger("agent"),
	}
}

// TODO load the existing configuration on startup and start the bridge if needed
func (a *Agent) start() error {

	return nil
}

// stop all the things.
func (a *Agent) stop() error {

	return nil
}

func (a *Agent) startBridge(connect *connectRequest) error {
	return a.bridge.start(connect.Url, connect.Token)
}

// save the state of the bridge then disconnect it
func (a *Agent) stopBridge(disconnect *disconnectRequest) error {
	return a.bridge.stop()
}

func (a *Agent) getMetrics() *metricsEvent {
	return a.metrics.buildMetricsRequest()
}

func (a *Agent) getStatus() statsEvent {

	var lastError string

	if a.bridge.LastError != nil {
		lastError = a.bridge.LastError.Error()
	}

	runtime.ReadMemStats(a.memstats)

	return statsEvent{
		LastError:      lastError,
		Alloc:          a.memstats.Alloc,
		HeapAlloc:      a.memstats.HeapAlloc,
		TotalAlloc:     a.memstats.TotalAlloc,
		Connected:      a.bridge.IsConnected(),
		Configured:     a.bridge.Configured,
		Timestamp:      time.Now().Unix(),
		IngressCounter: a.bridge.IngressCounter,
		IngressBytes:   a.bridge.IngressBytes,
		EgressCounter:  a.bridge.EgressCounter,
		EgressBytes:    a.bridge.EgressBytes,
	}
}
