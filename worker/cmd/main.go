package main

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sypheos/fogflow/worker"
)

func generateID(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func main() {
	config := worker.LoadConfig()

	// overwrite the configuration with environment variables
	if value, exist := os.LookupEnv("myip"); exist {
		config.MyIP = value
	}
	if value, exist := os.LookupEnv("discoveryURL"); exist {
		config.DiscoveryURL = value
	}
	if value, exist := os.LookupEnv("rabbitmq"); exist {
		config.MessageBus = value
	}

	// start the wrk to deal with tasks
	var wrk = &worker.Worker{}
	ok := wrk.Start(&config)
	if ok == false {
		worker.ERROR.Println("failed to start the wrk instance")
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	wrk.Quit()
}
