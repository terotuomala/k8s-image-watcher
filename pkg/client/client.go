package client

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetClient() kubernetes.Interface {
	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		log.WithFields(log.Fields{"pkg": "client.go"}).Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		log.WithFields(log.Fields{"pkg": "client.go"}).Fatal(err)
	}

	return clientset
}
