package main

import (
	"testing"
	"time"
)

// sig возвращает канал, который закроется через after
func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func TestOrNoChannels(t *testing.T) {
	ch := or()
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("ожидался закрытый канал")
		}
	case <-time.After(50 * time.Millisecond):
		t.Error("чтение из or() без аргументов не должно блокировать")
	}
}

func TestOrSingleChannel(t *testing.T) {
	start := time.Now()
	<-or(
		sig(30 * time.Millisecond),
	)
	if time.Since(start) < 30*time.Millisecond {
		t.Error("канал должен закрыться после истечения времени sig()")
	}
}

func TestOrMultipleChannels(t *testing.T) {
	start := time.Now()
	<-or(
		sig(80*time.Millisecond),
		sig(40*time.Millisecond),
		sig(100*time.Millisecond),
	)
	elapsed := time.Since(start)

	if elapsed < 40*time.Millisecond || elapsed > 60*time.Millisecond {
		t.Errorf("ожидали закрытие примерно через 40ms, получили %v", elapsed)
	}
}

func TestOrImmediateClose(t *testing.T) {
	start := time.Now()
	<-or(
		sig(0*time.Millisecond),
		sig(50*time.Millisecond),
	)
	if time.Since(start) > 5*time.Millisecond {
		t.Error("выходной канал должен закрыться мгновенно, если один из входных уже закрыт")
	}
}
