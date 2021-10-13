package appentrance

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)


type government struct {
	modules       []Department
	closeChannels map[string]chan struct{}
}

func NewGovernment() Government {
	return &government{
		closeChannels: make(map[string]chan struct{}),
	}
}

func (gmt *government) Register(department Department) error {
	if _, ok := gmt.closeChannels[department.Name()]; ok {
		return fmt.Errorf("department already registered")
	}

	gmt.modules = append(gmt.modules, department)
	gmt.closeChannels[department.Name()] = make(chan struct{}, 1)
	return nil
}

func (gmt *government) Start() []error {
	var (
		errors []error
	)

	for _, department := range gmt.modules {
		err := department.Start(gmt.closeChannels[department.Name()])
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (gmt *government) GracefulShutdown() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT)
	<-c

	for _, closeChannel := range gmt.closeChannels {
		closeChannel <- struct{}{}
	}

	return nil
}