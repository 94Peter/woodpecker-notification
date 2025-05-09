package settings

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	Provider_portainer = "portainer"
)

func NewPortainerSettings(provider *WebhookProvider) (*PortainerSettings, error) {
	result := &PortainerSettings{
		webhook:  provider.Webhook,
		bindData: provider.bindData,
	}
	if provider.Configs != nil && len(*provider.Configs) > 2 {
		err := json.Unmarshal(*provider.Configs, result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type PortainerSettings struct {
	TagKey    string `json:"tagkey"`
	TagValue  string `json:"tagvalue"`
	PullImage bool   `json:"pullimage"`

	webhook  string
	bindData map[string]string
}

func (portainer *PortainerSettings) GetWebhookURL() string {
	// url paramter generate
	myurl, err := url.Parse(portainer.webhook)
	if err != nil {
		fmt.Println("parse url error: ", err)
		return portainer.webhook
	}
	values := url.Values{}
	if portainer.TagKey != "" && portainer.TagValue != "" {
		values.Add(portainer.TagKey, portainer.TagValue)
	}
	values.Add("pullimage", fmt.Sprintf("%t", portainer.PullImage))
	myurl.RawQuery = values.Encode()
	url := myurl.String()
	fmt.Println(url)
	return url
}
