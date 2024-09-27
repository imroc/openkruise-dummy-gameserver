package prom

import (
	"net/http"
	"os"

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

func StartServer(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, mux)
}
