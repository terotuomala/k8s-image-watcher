/*
Copyright © 2023 Tero Tuomala
*/
package cmd

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/terotuomala/k8s-image-watcher/pkg/config"
	"github.com/terotuomala/k8s-image-watcher/pkg/controller"
	"github.com/terotuomala/k8s-image-watcher/pkg/logging"
	"github.com/terotuomala/k8s-image-watcher/pkg/utils"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Start watching image changes in Kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		envMsg := "set to:"

		conf := config.New()

		loggingLevel, err := utils.GetEnvStr("LOG_LEVEL")
		if err != nil {
			log.WithFields(log.Fields{"pkg": "watch.go"}).Fatal(loggingLevel, err)
		}
		log.WithFields(log.Fields{"pkg": "watch.go"}).Infof(`LOG_LEVEL %s %s`, envMsg, loggingLevel)
		conf.Logging.Level = loggingLevel

		logging.InitLogging(conf)

		slackChannel, err := utils.GetEnvStr("SLACK_CHANNEL")
		if err != nil {
			log.WithFields(log.Fields{"pkg": "watch.go"}).Fatal(slackChannel, err)
		}
		log.WithFields(log.Fields{"pkg": "watch.go"}).Infof(`SLACK_CHANNEL %s %s`, envMsg, slackChannel)
		conf.Handler.Slack.Channel = slackChannel

		slackMessageTitle, err := utils.GetEnvStr("SLACK_MESSAGE_TITLE")
		if err != nil {
			log.WithFields(log.Fields{"pkg": "watch.go"}).Fatal(slackMessageTitle, err)
		}
		log.WithFields(log.Fields{"pkg": "watch.go"}).Infof(`SLACK_MESSAGE_TITLE %s %s`, envMsg, slackMessageTitle)
		conf.Handler.Slack.Title = slackMessageTitle

		slackToken, err := utils.GetEnvStr("SLACK_TOKEN")
		if err != nil {
			log.WithFields(log.Fields{"pkg": "watch.go"}).Fatal(slackToken, err)
		}
		log.WithFields(log.Fields{"pkg": "watch.go"}).Infof(`SLACK_TOKEN %s <reducted>`, envMsg)
		conf.Handler.Slack.Token = slackToken

		namespace, err := utils.GetEnvStr("NAMESPACE")
		if err != nil {
			log.WithFields(log.Fields{"pkg": "watch.go"}).Fatal(namespace, err)
		}
		log.WithFields(log.Fields{"pkg": "watch.go"}).Infof(`NAMESPACE %s %s`, envMsg, namespace)
		conf.Namespace = namespace

		enableDeployment, err := utils.GetEnvBool("WATCH_DEPLOYMENT")
		if err != nil {
			log.WithFields(log.Fields{"pkg": "watch.go"}).Fatal(enableDeployment, err)
		}
		log.WithFields(log.Fields{"pkg": "watch.go"}).Infof(`WATCH_DEPLOYMENT %s %s`, envMsg, strconv.FormatBool(enableDeployment))
		conf.Resource.Deployment = enableDeployment

		enableDaemonset, err := utils.GetEnvBool("WATCH_DAEMONSET")
		if err != nil {
			log.WithFields(log.Fields{"pkg": "watch.go"}).Fatal(enableDaemonset, err)
		}
		log.WithFields(log.Fields{"pkg": "watch.go"}).Infof(`WATCH_DAEMONSET %s %s`, envMsg, strconv.FormatBool(enableDaemonset))
		conf.Resource.DaemonSet = enableDaemonset

		enableStatefulset, err := utils.GetEnvBool("WATCH_STATEFULSET")
		if err != nil {
			log.WithFields(log.Fields{"pkg": "watch.go"}).Fatal(enableStatefulset, err)
		}
		log.WithFields(log.Fields{"pkg": "watch.go"}).Infof(`WATCH_STATEFULSET %s %s`, envMsg, strconv.FormatBool(enableStatefulset))
		conf.Resource.StatefulSet = enableStatefulset

		controller.Create(conf)
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
