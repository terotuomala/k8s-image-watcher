package slack

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/terotuomala/k8s-image-watcher/pkg/config"
)

type SlackNotifier struct {
	client  *slack.Client
	channel string
	title   string
}

func NewSlackNotifier(conf *config.Slack) (*SlackNotifier, error) {
	if err := validateSlackConfig(conf); err != nil {
		return nil, err
	}

	return &SlackNotifier{
		client:  slack.New(conf.Token),
		channel: conf.Channel,
		title:   conf.Title,
	}, nil
}

func validateSlackConfig(conf *config.Slack) error {
	if conf.Token == "" {
		log.WithFields(log.Fields{"pkg": "slack.go"}).Fatal("slack token is required")
	}
	if conf.Channel == "" {
		log.WithFields(log.Fields{"pkg": "slack.go"}).Fatal("slack channel is required")
	}

	return nil
}

func (sn *SlackNotifier) SendMessage(message string) error {
	_, _, err := sn.client.PostMessage(sn.channel,
		slack.MsgOptionText(fmt.Sprintf("*%s*\n%s", sn.title, message), false),
	)
	return err
}
