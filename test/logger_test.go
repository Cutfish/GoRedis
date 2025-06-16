package test

import (
	"go-redis/lib/logger"
	"testing"
)

func TestLogger(t *testing.T) {
	logger.Info("ceshi1")
	logger.Error("ceshi1")
}
