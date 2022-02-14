package main

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/config"
	"log"
)

func main() {
	configData, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	runServer(configData)
}
