package main

import (
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/config"
	gcpCustomAuth "github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/pkg/k8s/auth/gcp_custom"
	"log"
	"time"
)

func main() {
	configData, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	time.Local = time.UTC
	gcpCustomAuth.RegisterK8SGCPCustomAuthProvider()

	runService(configData)
}
