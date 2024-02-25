/*
Copyright © 2023 Tero Tuomala
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/terotuomala/k8s-image-watcher/pkg/config"
	"github.com/terotuomala/k8s-image-watcher/pkg/controller"
	"github.com/terotuomala/k8s-image-watcher/pkg/logging"
	"github.com/terotuomala/k8s-image-watcher/pkg/utils"
)

func setConfigFromEnvVar(envVar string, setterFunc func(string), pkgLog string) {
	value, err := utils.GetEnvStr(envVar)
	if err != nil {
		log.WithFields(log.Fields{"pkg": pkgLog}).Fatal(envVar, err)
	}
	log.WithFields(log.Fields{"pkg": pkgLog}).Infof("%s set to: %s", envVar, value)
	setterFunc(value)
}

func setBoolConfigFromEnvVar(envVar string, setterFunc func(bool), pkgLog string) {
	value, err := utils.GetEnvBool(envVar)
	if err != nil {
		log.WithFields(log.Fields{"pkg": pkgLog}).Fatal(envVar, err)
	}
	log.WithFields(log.Fields{"pkg": pkgLog}).Infof("%s set to: %v", envVar, value)
	setterFunc(value)
}

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Start watching image changes in Kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.New()
		pkgLog := "watch.go"

		setConfigFromEnvVar("LOG_LEVEL", func(value string) { conf.Logging.Level = value }, pkgLog)
		setConfigFromEnvVar("SLACK_CHANNEL", func(value string) { conf.Handler.Slack.Channel = value }, pkgLog)
		setConfigFromEnvVar("SLACK_MESSAGE_TITLE", func(value string) { conf.Handler.Slack.Title = value }, pkgLog)
		setConfigFromEnvVar("SLACK_TOKEN", func(value string) { conf.Handler.Slack.Token = value }, pkgLog)
		setConfigFromEnvVar("NAMESPACE", func(value string) { conf.Namespace = value }, pkgLog)

		setBoolConfigFromEnvVar("WATCH_DEPLOYMENT", func(value bool) { conf.Resource.Deployment = value }, pkgLog)
		setBoolConfigFromEnvVar("WATCH_DAEMONSET", func(value bool) { conf.Resource.DaemonSet = value }, pkgLog)
		setBoolConfigFromEnvVar("WATCH_STATEFULSET", func(value bool) { conf.Resource.StatefulSet = value }, pkgLog)

		logging.InitLogging(conf)
		controller.Create(conf)
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
