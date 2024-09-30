package prom

import (
	"net/http"
	"os"

	dto "github.com/prometheus/client_model/go"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var podName, podNamespace string

func init() {
	podName = os.Getenv("POD_NAME")
	podNamespace = os.Getenv("POD_NAMESPACE")
	if podName == "" || podNamespace == "" {
		panic("env POD_NAME or POD_NAMESPACE is required")
	}
}

var roomTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "game",
	Subsystem: "room",
	Name:      "total",
}, []string{
	"pod",
	"namespace",
})

var allocatedRoomTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "game",
	Subsystem: "room",
	Name:      "allocated_total",
}, []string{
	"pod",
	"namespace",
})

// 新建 N 个空闲房间
func AddNewRoom(n int) {
	roomTotal.WithLabelValues(podName, podNamespace).Add(float64(n))
}

// 分配 N 个房间
func AllocateRoom(n int) {
	allocatedRoomTotal.WithLabelValues(podName, podNamespace).Add(float64(n))
}

// 释放 N 个房间（reuse 表示房间是否可复用）
func ReleaseRoom(n int, reuse bool) {
	if !reuse {
		roomTotal.WithLabelValues(podName, podNamespace).Sub(float64(n))
	}
	allocatedRoomTotal.WithLabelValues(podName, podNamespace).Sub(float64(n))
}

func SetBusy() {
	currentTotalMetric := &dto.Metric{}
	currentAllocatedMetric := &dto.Metric{}
	if err := roomTotal.WithLabelValues(podName, podNamespace).Write(currentTotalMetric); err != nil {
		panic(err)
	}
	if err := allocatedRoomTotal.WithLabelValues(podName, podNamespace).Write(currentAllocatedMetric); err != nil {
		panic(err)
	}
	currentTotal := *currentTotalMetric.Gauge.Value
	currentAllocated := *currentAllocatedMetric.Gauge.Value
	if idle := currentTotal - currentAllocated; idle > 0.0 {
		AllocateRoom(int(idle))
	}
}

func SetIdle() {
	currentAllocatedMetric := &dto.Metric{}
	if err := allocatedRoomTotal.WithLabelValues(podName, podNamespace).Write(currentAllocatedMetric); err != nil {
		panic(err)
	}
	currentAllocated := *currentAllocatedMetric.Gauge.Value
	if currentAllocated > 0.0 {
		ReleaseRoom(int(currentAllocated), true)
	}
}

func IsAllIdle() bool {
	metric, err := allocatedRoomTotal.GetMetricWithLabelValues(podName, podNamespace)
	if err != nil {
		panic(err)
	}
	model := &dto.Metric{}
	if err := metric.Write(model); err != nil {
		panic(err)
	}
	return *model.Gauge.Value == 0
}

func StartServer(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, mux)
}
