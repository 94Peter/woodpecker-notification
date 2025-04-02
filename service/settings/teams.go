package settings

import (
	"encoding/json"
	"os"
	"strings"
)

const (
	Provider_Teams = "teams"
)

type TeamsSettings struct {
	webhookUrl  string
	bindData    map[string]string
	Facts       []string `json:"facts"`
	Buttons     []string `json:"buttons"`
	Variables   []string `json:"variables"`
	variableMap map[string]string
}

func NewTeamsSettings(settings *WebhookProvider) (*TeamsSettings, error) {
	teams := &TeamsSettings{
		webhookUrl: settings.Webhook,
		bindData:   settings.bindData,
	}
	if settings.Configs != nil {
		err := json.Unmarshal(*settings.Configs, teams)
		if err != nil {
			return nil, err
		}
	}
	for _, variable := range teams.Variables {
		teams.variableMap[variable] = os.Getenv(variable)
	}
	return teams, nil
}

func (teams *TeamsSettings) GetWebhookURL() string {
	return teams.webhookUrl
}

func (teams *TeamsSettings) GetVersion() string {
	return teams.bindData["BUILD_TAG"]
}

func (teams *TeamsSettings) GetPreBuildStatus() string {
	return teams.bindData["CI_PREV_PIPELINE_STATUS"]
}

func (teams *TeamsSettings) GetAuthorAvatar() string {
	return teams.bindData["CI_COMMIT_AUTHOR_AVATAR"]
}

func (teams *TeamsSettings) GetRepo() string {
	return teams.bindData["CI_REPO"]
}

// GetCommitMessage
func (teams *TeamsSettings) GetCommitMessage() string {
	return strings.Split(teams.bindData["CI_COMMIT_MESSAGE"], "\n")[0]
}

// GetCommitAuthor
func (teams *TeamsSettings) GetCommitAuthor() string {
	return teams.bindData["CI_COMMIT_AUTHOR"]
}

// GetVariableValue
func (teams *TeamsSettings) GetVariableValue(name string) string {
	return teams.variableMap[name]
}

// GetPipelineURL
func (teams *TeamsSettings) GetPipelineURL() string {
	return teams.bindData["CI_PIPELINE_URL"]
}

// GetPipelineForgeURL
func (teams *TeamsSettings) GetPipelineForgeURL() string {
	return teams.bindData["CI_PIPELINE_FORGE_URL"]
}

// GetCommitTag
func (teams *TeamsSettings) GetCommitTag() string {
	return teams.bindData["CI_COMMIT_TAG"]
}
