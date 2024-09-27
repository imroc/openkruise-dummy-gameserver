package main

import (
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/imroc/openkruise-dummy-gameserver/pkg/prom"
)

func main() {
	promAddr := os.Getenv("METRICS_ADDR")
	if promAddr == "" {
		promAddr = ":8099"
	}

	var err error
	roomTotalNumber := 1
	if roomTotal := os.Getenv("ROOM_TOTAL"); roomTotal != "" {
		roomTotalNumber, err = strconv.Atoi(roomTotal)
		if err != nil {
			panic(err)
		}
	}
	allocatedRatioNumber := 0
	if allocatedRatio := os.Getenv("ALLOCATED_ROOM_RATIO"); allocatedRatio != "" {
		allocatedRatioNumber, err = strconv.Atoi(allocatedRatio)
		if err != nil {
			panic(err)
		}
	}

	ratio := rand.Float64() + (float64(allocatedRatioNumber)-50)/100

	// 根据配置的占用比例，结合随机数计算出当前 Pod 房间分配数量
	allocatedNumber := int(math.Round(float64(roomTotalNumber) * ratio))

	prom.AddNewRoom(roomTotalNumber)
	prom.AllocateRoom(allocatedNumber)

	if err := prom.StartServer(promAddr); err != nil {
		panic(err)
	}
}
