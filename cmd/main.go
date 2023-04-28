package main

import (
	"context"
	"github.com/boriszhilko/ports-service/internal/port/adapter/in/file"
	"github.com/boriszhilko/ports-service/internal/port/adapter/out/persistence/redis"
	ports "github.com/boriszhilko/ports-service/internal/port/service"
	"github.com/boriszhilko/ports-service/pkg/tools"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := configureLogging(); err != nil {
		log.Fatal(err)
	}

	repository := redis.NewRedisRepository(os.Getenv("REDIS_URL"))

	input, err := file.NewInput(os.Getenv("PORTS_FILE_ADDRESS"))
	if err != nil {
		log.Fatal(err)
	}

	service, err := ports.NewPortService(input, repository)
	if err != nil {
		log.Fatal(err)
	}

	stoppable := tools.NewStoppable(repository, input)
	shutdownCompleteChan := stoppable.ConfigureGracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = service.CreateOrUpdatePorts(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Ports inserted successfully")
	<-shutdownCompleteChan
}

func configureLogging() error {
	l, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return err
	}
	log.SetLevel(l)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat:   time.RFC3339Nano,
		DisableHTMLEscape: true,
	})
	return nil
}
