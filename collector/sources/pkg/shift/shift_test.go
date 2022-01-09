package shift

import (
	"collector/pkg/config"
	"log"
	"testing"
)

func TestGetCurrentShift(t *testing.T) {
	log.Println(config.Cfg.Shifts)
	log.Println(GetCurrentShift())
}
