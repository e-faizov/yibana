package main

import (
	"fmt"
	"time"

	"github.com/e-faizov/yibana/internal"
)

const pollInterval = 2
const reportInterval = 10

func main() {

	pollTicker := time.NewTicker(pollInterval * time.Second)
	reportTicker := time.NewTicker(reportInterval * time.Second)

	metrics := internal.Metrics{}
	metrics.Update()

	sender := internal.NewSender("localhost", 8080)

	go func() {
		for range pollTicker.C {
			metrics.Update()
		}
	}()

	for range reportTicker.C {
		for {
			next, ok := metrics.Front()
			if !ok {
				break
			}

			err := sender.SendMetric(next)
			if err != nil {
				fmt.Println("Ошибка отправки, попробуем в следующий раз", err.Error())
			} else {
				metrics.Pop()
			}

		}
	}
}
