package main

import (
	"math/rand"
	"time"

	"github.com/RomanIkonnikov93/URLshortner/cmd/config"
	"github.com/RomanIkonnikov93/URLshortner/internal/repository"
	"github.com/RomanIkonnikov93/URLshortner/internal/server"
	"github.com/RomanIkonnikov93/URLshortner/logging"
)

func main() {

	rand.Seed(time.Now().UnixMicro())

	logger := logging.GetLogger()

	cfg, err := config.GetConfig()
	if err != nil {
		logger.Fatalf("GetConfig: %s", err)
	}

	rep, err := repository.NewReps(*cfg)
	if err != nil {
		logger.Fatalf("NewReps: %s", err)
	}

	err = server.StartServer(*rep, *cfg, *logger)
	if err != nil {
		logger.Fatalf("StartServer: %s", err)
	}

}
