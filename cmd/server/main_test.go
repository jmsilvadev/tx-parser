package main

import (
	"testing"
	"time"

	"github.com/jmsilvadev/tx-parser/pkg/config"
)

func TestRun(t *testing.T) {
	c := config.GetDefaultConfig()
	go run(c)
	time.Sleep(time.Second)
}
