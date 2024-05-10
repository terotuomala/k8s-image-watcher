package logging

import (
	log "github.com/sirupsen/logrus"
	"github.com/terotuomala/k8s-image-watcher/pkg/config"
)

func InitLogging(conf *config.Logging) {
	log.SetFormatter(&log.TextFormatter{})

	switch conf.Level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "":
		log.SetLevel(log.InfoLevel)
	default:
		log.Fatal("Invalid logging level: ", conf.Level)
	}
}
