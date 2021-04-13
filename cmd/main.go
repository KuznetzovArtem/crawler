package main

import (
	"crawler/actions/get_titles"
	"crawler/env"
	"crawler/http"
	"crawler/repositories/logger"
	"crawler/task"
	"fmt"
	httpnet "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	host := env.GetEnvString("API_ADDRESS", ":8004")
	errorChannel := make(chan error)

	config := map[logger.LogType]logger.BatchLogger{}

	var transport = httpnet.DefaultTransport

	var Actions = map[string]task.ActionRunner{
		"get_urls": get_titles.NewService(transport),
	}

	log := logger.NewLogService(config, time.Duration(4)*time.Second)

	taskService := task.NewService(log, Actions)
	// HTTP ClientService.
	go func() {
		httpServer := http.NewServer(
			host,
			taskService,
		)
		errorChannel <- httpServer.Open()
	}()
	go log.LaunchLogging()
	// Capture interrupts.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-time.After(10 * time.Second):
			errorChannel <- fmt.Errorf("got signal: %s", <-c)
		}
	}()
	// Wait (indefinitely) for any error.
	if err := <-errorChannel; err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
