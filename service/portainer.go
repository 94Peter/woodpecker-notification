package service

import "woodpecker-webhook/service/settings"

type portainer struct {
	setting *settings.PortainerSettings
}

func newPortainerSendMessage(set *settings.PortainerSettings) SendMessage {
	chat := &portainer{
		setting: set,
	}
	return chat.SendMessage
}

func (port *portainer) SendMessage() error {
	return postWebhook(port.setting.GetWebhookURL(), nil)
}
