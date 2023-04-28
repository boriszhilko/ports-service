package tools

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// Stopper is an interface for objects that can be stopped gracefully.
type Stopper interface {
	Stop() error
}

type Stoppable struct {
	stoppables []Stopper
}

func NewStoppable(stoppables ...Stopper) *Stoppable {
	return &Stoppable{stoppables: stoppables}
}

func (s *Stoppable) ConfigureGracefulStop() chan struct{} {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	shutdownCompleteChan := make(chan struct{})

	go func() {
		<-sigChan
		log.Info("Graceful shutdown initiated")

		var err error
		for _, stoppable := range s.stoppables {
			if err = stoppable.Stop(); err != nil {
				log.Error("Error stopping stoppable: ", err)
			}
		}
		log.Info("Graceful shutdown completed")
		shutdownCompleteChan <- struct{}{}
	}()

	return shutdownCompleteChan
}
