package wg

import "sync"

var wg sync.WaitGroup

func Add() {
	wg.Add(1)
}

func Done() {
	wg.Done()
}

func Wait() {
	wg.Wait()
}
