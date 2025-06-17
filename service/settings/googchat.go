package settings

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"html/template"
)

const (
	Provider_googleChat = "google_chat"

	defaultGoogleChatMsg = `{"cardsV2":[{"cardId":"{{.CARD_ID}}","card":{"header":{"title":"{{.CI_REPO}}","subtitle":"@{{.CI_COMMIT_AUTHOR}} has updated","imageUrl":"{{.CI_COMMIT_AUTHOR_AVATAR}}","imageType":"CIRCLE"},"sections":[{"collapsible":false,"widgets":[{"decoratedText":{"icon":{"iconUrl":"{{.STATUS_ICON_URL}}"},"text":"{{.CI_PREV_PIPELINE_STATUS}}"}},{"decoratedText":{"icon":{"materialIcon":{"name":"comment"}},"topLabel":"Commit Message","text":"{{.CI_COMMIT_MESSAGE}}"}},{"decoratedText":{"icon":{"materialIcon":{"name":"linked_services"}},"topLabel":"Branch","text":"{{.CI_COMMIT_REF}}"}},{"decoratedText":{"icon":{"materialIcon":{"name":"label"}},"topLabel":"Tag","text":"{{.BUILD_TAG}}"}},{"buttonList":{"buttons":[{"text":"View Release","icon":{"iconUrl":"{{.VIEW_RELEASE_ICON_URL}}","altText":"view release"},"type":"OUTLINED","onClick":{"openLink":{"url":"{{.CI_PIPELINE_FORGE_URL}}"}}},{"text":"View Pipleline","icon":{"iconUrl":"{{.VIEW_PIPLELINE_ICON_URL}}","altText":"view pipleline"},"type":"OUTLINED","onClick":{"openLink":{"url":"{{.CI_PIPELINE_URL}}"}}}]}}]}]}}]}`

	releaseIcon  = "https://raw.githubusercontent.com/94Peter/woodpecker-notification/refs/heads/main/assets/git-icon.png"
	pipelineIcon = "https://raw.githubusercontent.com/94Peter/woodpecker-notification/refs/heads/main/assets/woodpecker-icon.png"
	successIcon  = "https://raw.githubusercontent.com/94Peter/woodpecker-notification/refs/heads/main/assets/success-icon.png"
	failIcon     = "https://raw.githubusercontent.com/94Peter/woodpecker-notification/refs/heads/main/assets/failure-icon.png"
)

func NewGoogleChatSettings(provider *WebhookProvider) (*GoogleChatSettings, error) {
	result := &GoogleChatSettings{
		webhook:  provider.Webhook,
		bindData: provider.bindData,
	}
	if provider.Configs != nil && len(*provider.Configs) > 2 {
		err := json.Unmarshal(*provider.Configs, result)
		if err != nil {
			return nil, err
		}
	} else {
		result.Message = defaultGoogleChatMsg
	}
	fmt.Println("google setting: ", result.MsgFile)
	tpl, err := template.New("googleMsg").Parse(result.Message)
	if err != nil {
		return nil, err
	}
	result.msgTpl = tpl
	return result, nil
}

type GoogleChatSettings struct {
	Message string `json:"message"`
	MsgFile string `json:"msg_file"`

	webhook  string
	msgTpl   *template.Template
	bindData map[string]string
}

func getUniqueCardID() string {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("gcard:%s", hex.EncodeToString(bytes))
}

func (google *GoogleChatSettings) getBindData() map[string]string {
	bindData := google.bindData
	bindData["VIEW_RELEASE_ICON_URL"] = releaseIcon
	bindData["VIEW_PIPLELINE_ICON_URL"] = pipelineIcon
	statusIcon := successIcon
	if bindData["DRONE_BUILD_STATUS"] == "failure" {
		statusIcon = failIcon
	}
	bindData["STATUS_ICON_URL"] = statusIcon
	bindData["CARD_ID"] = getUniqueCardID()
	return bindData
}

func (google *GoogleChatSettings) GetMsgReader() (io.Reader, error) {
	if google.MsgFile != "" {

		file, err := os.ReadFile(google.MsgFile)
		if err != nil {
			return nil, err
		}
		fmt.Println("google msg file: ", string(file))
		return bytes.NewReader(file), nil
	}
	var data bytes.Buffer
	err := google.msgTpl.Execute(&data, google.getBindData())
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (google *GoogleChatSettings) GetWebhookURL() string {
	return google.webhook
}
