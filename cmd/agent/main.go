package main

import (
	"github.com/e-faizov/yibana/internal"
	"time"
)

const pollInterval = 2
const reportInterval = 10

func main() {

	ticker := time.NewTicker(1 * time.Second)
	metrics := internal.NewMetrics()
	sender := internal.NewSender("127.0.0.1", 8080)
	metrics.Update()

	var tickCount int64
	for range ticker.C {
		tickCount++

		if tickCount%pollInterval == 0 {
			metrics.Update()
			if tickCount%reportInterval == 0 {
				sender.SendMetrics(metrics)
			}
		}
	}

}
