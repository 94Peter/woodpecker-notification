package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"woodpecker-webhook/service/settings"
)

type SendMessage func() error

func GetSendMessageFunSlice() ([]SendMessage, error) {
	set, err := settings.GetSettings()
	if err != nil {
		return nil, err
	}

	result := make([]SendMessage, len(set.Webhooks))
	for i, hook := range set.Webhooks {
		switch hook.Provider {
		case settings.Provider_googleChat:
			googleChatSetting, err := settings.NewGoogleChatSettings(hook)
			if err != nil {
				return nil, err
			}
			result[i] = newGoogleChantSendMessage(googleChatSetting)
		case settings.Provider_Teams:
			teamSettings, err := settings.NewTeamsSettings(hook)
			if err != nil {
				return nil, err
			}
			result[i] = newTeamsSendMessage(teamSettings)
		case settings.Provider_portainer:
			portainerSettings, err := settings.NewPortainerSettings(hook)
			if err != nil {
				return nil, err
			}
			result[i] = newPortainerSendMessage(portainerSettings)
		default:
			return nil, errors.New("Webhooks containe not support provider")
		}
	}
	return result, nil
}

func postWebhook(url string, data io.Reader) error {
	resp, err := http.Post(url, "application/json", data)
	if err != nil {
		return fmt.Errorf("Error sending webhook: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Error status code [%d] response from webhook: %s", resp.StatusCode, string(body))
	}
	return nil
}
