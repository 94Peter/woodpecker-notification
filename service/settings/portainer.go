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
	TagKey    *string `json:"tagkey,omitempty"`
	TagValue  *string `json:"tagvalue,omitempty"`
	PullImage *bool   `json:"pullimage,omitempty"`

	webhook  string
	bindData map[string]string
}

func (portainer *PortainerSettings) GetWebhookURL() string {
	myurl, err := url.Parse(portainer.webhook)
	if err != nil {
		return portainer.webhook
	}
	values := url.Values{}
	if portainer.TagKey != nil && portainer.TagValue != nil {
		values.Add(*portainer.TagKey, *portainer.TagValue)
	}
	if portainer.PullImage != nil {
		values.Add("pullimage", fmt.Sprintf("%t", *portainer.PullImage))
	}
	myurl.RawQuery = values.Encode()
	return myurl.String()
}
