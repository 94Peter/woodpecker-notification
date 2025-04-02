package service

import (
	"woodpecker-webhook/service/settings"
)

type googleChat struct {
	setting *settings.GoogleChatSettings
}

func newGoogleChantSendMessage(set *settings.GoogleChatSettings) SendMessage {
	chat := &googleChat{
		setting: set,
	}
	return chat.SendMessage
}

func (google *googleChat) SendMessage() error {
	msgReader, err := google.setting.GetMsgReader()
	if err != nil {
		return err
	}
	return postWebhook(google.setting.GetWebhookURL(), msgReader)
}
