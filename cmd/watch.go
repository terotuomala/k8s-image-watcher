package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/terotuomala/k8s-image-watcher/pkg/config"
	"github.com/terotuomala/k8s-image-watcher/pkg/controller"
	slack "github.com/terotuomala/k8s-image-watcher/pkg/handler"
	"github.com/terotuomala/k8s-image-watcher/pkg/logging"
	"github.com/terotuomala/k8s-image-watcher/pkg/utils"
)

func setConfigFromEnvVar(envVar string, setterFunc func(string), pkgLog string, sensitive bool) {
	value, err := utils.GetEnvStr(envVar)
	if err != nil {
		log.WithFields(log.Fields{"pkg": pkgLog}).Fatal(envVar, err)
	}

	if sensitive {
		log.WithFields(log.Fields{"pkg": pkgLog}).Infof("%s=[REDACTED]", envVar)
	} else {
		log.WithFields(log.Fields{"pkg": pkgLog}).Infof("%s=%s", envVar, value)
	}

	setterFunc(value)
}

func setBoolConfigFromEnvVar(envVar string, setterFunc func(bool), pkgLog string) {
	value, err := utils.GetEnvBool(envVar)
	if err != nil {
		log.WithFields(log.Fields{"pkg": pkgLog}).Fatal(envVar, err)
	}
	log.WithFields(log.Fields{"pkg": pkgLog}).Infof("%s=%v", envVar, value)
	setterFunc(value)
}

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Start watching image changes in Kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.New()
		pkgLog := "watch.go"

		setConfigFromEnvVar("LOG_LEVEL", func(value string) { conf.Logging.Level = value }, pkgLog, false)
		setConfigFromEnvVar("NAMESPACE", func(value string) { conf.Namespace = value }, pkgLog, false)
		setBoolConfigFromEnvVar("WATCH_DEPLOYMENT", func(value bool) { conf.Resource.Deployment = value }, pkgLog)
		setBoolConfigFromEnvVar("WATCH_DAEMONSET", func(value bool) { conf.Resource.DaemonSet = value }, pkgLog)
		setBoolConfigFromEnvVar("WATCH_STATEFULSET", func(value bool) { conf.Resource.StatefulSet = value }, pkgLog)
		setBoolConfigFromEnvVar("SLACK_ENABLED", func(value bool) { conf.SlackEnabled = value }, pkgLog)

		if conf.SlackEnabled {
			setConfigFromEnvVar("SLACK_CHANNEL", func(value string) { conf.Handler.Slack.Channel = value }, pkgLog, false)
			setConfigFromEnvVar("SLACK_MESSAGE_TITLE", func(value string) { conf.Handler.Slack.Title = value }, pkgLog, false)
			setConfigFromEnvVar("SLACK_TOKEN", func(value string) { conf.Handler.Slack.Token = value }, pkgLog, true)

			slackClient, err := slack.NewSlackNotifier(&conf.Handler.Slack)
			if err != nil {
				log.Fatalf("Failed to initialize Slack notifier: %v", err)
			}
			controller.Create(conf, slackClient)

		} else {
			controller.Create(conf, nil)
		}

		logging.InitLogging(&conf.Logging)
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
