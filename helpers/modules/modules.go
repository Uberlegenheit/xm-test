package modules

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"xm-task/log"
)

type Module interface {
	Run() error
	Stop() error
	Title() string
}

var gracefulTimeout = time.Second * 15
var makePanicIfError = true

func Stop(modules []Module) {
	wg := &sync.WaitGroup{}
	wg.Add(len(modules))
	for _, m := range modules {
		go func(m Module) {
			err := stopModule(m)
			if err != nil {
				log.Error("Module stopped with error", zap.String("module", m.Title()), zap.Error(err))
			}
			wg.Done()
		}(m)
	}
	wg.Wait()
	log.Info("All modules was stopped")
}

func stopModule(m Module) error {
	if m == nil {
		return nil
	}
	result := make(chan error)
	go func() {
		result <- m.Stop()
	}()
	select {
	case err := <-result:
		return err
	case <-time.After(gracefulTimeout):
		return fmt.Errorf("stoped by timeout")
	}
}

func Run(modules []Module) {
	type errResp struct {
		err    error
		module string
	}
	errors := make(chan errResp, len(modules))
	for _, m := range modules {
		go func(m Module) {
			err := m.Run()
			errR := errResp{
				err:    err,
				module: m.Title(),
			}
			errors <- errR
		}(m)
	}
	// handle errors
	go func() {
		for {
			err := <-errors
			if err.err != nil {
				log.Error("Module return error", zap.String("module", err.module), zap.Error(err.err))
				if makePanicIfError {
					Stop(modules)
					log.Fatal("Stopped via module error", zap.String("module", err.module), zap.Error(err.err))
				}
			}
			log.Info("Module finish work", zap.String("module", err.module))
		}
	}()
}
