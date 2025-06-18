package test

import (
	"GoRedis/lib/logger"
	"testing"
)

func TestLogger(t *testing.T) {
	logger.Info("ceshi1")
	logger.Error("ceshi1")
}
