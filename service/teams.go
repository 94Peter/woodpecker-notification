package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"
	"woodpecker-webhook/service/settings"
)

type teams struct {
	setting *settings.TeamsSettings
}

func newTeamsSendMessage(set *settings.TeamsSettings) SendMessage {
	chat := &teams{
		setting: set,
	}
	return chat.SendMessage
}

func (teams *teams) SendMessage() error {
	fmt.Println("teams send message")
	card := teams.createTeamsCard()
	cardBytes, err := json.Marshal(card)
	if err != nil {
		return err
	}

	return postWebhook(teams.setting.GetWebhookURL(), strings.NewReader(string(cardBytes)))
}

func (teams *teams) createTeamsCard() map[string]any {
	return map[string]any{
		"type": "message",
		"attachments": []map[string]any{
			{
				"contentType": "application/vnd.microsoft.card.adaptive",
				"contentUrl":  nil,
				"content": map[string]any{
					"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
					"type":    "AdaptiveCard",
					"version": "1.5",
					"body":    teams.createCardBody(),
					"actions": teams.createCardActions(),
				},
			},
		},
	}
}

func (teams *teams) createCardBody() []map[string]any {
	projectVersion := teams.setting.GetVersion()
	status := teams.setting.GetPreBuildStatus()
	color := "good"
	title := "✔ Pipeline succeeded"
	if status == "failure" {
		color = "attention"
		title = "❌ Pipeline failed"
	}
	dateStr := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	avatarURL := teams.setting.GetAuthorAvatar()
	var avatarDataURI string
	if avatarURL != "" {
		if dataURI, err := getAvatarDataURI(avatarURL); err == nil {
			avatarDataURI = dataURI
		} else {
			fmt.Printf("Warning: Failed to process avatar image: %v\n", err)
			avatarDataURI = avatarURL
		}
	}

	body := teams.createBaseBody(color, title, avatarDataURI, dateStr, projectVersion)

	if len(teams.setting.Variables) > 0 {
		body = teams.appendVariablesTable(body, teams.setting.Variables)
	}

	return body
}

func getAvatarDataURI(avatarURL string) (string, error) {
	resp, err := http.Get(avatarURL)
	if err != nil {
		return "", fmt.Errorf("failed to download avatar: %w", err)
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read avatar data: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	ext := path.Ext(avatarURL)
	switch {
	case contentType != "":
		// Use content type from header
	case ext != "":
		contentType = mime.TypeByExtension(ext)
	default:
		contentType = http.DetectContentType(imageData)
	}

	return fmt.Sprintf("data:%s;base64,%s",
		contentType,
		base64.StdEncoding.EncodeToString(imageData),
	), nil
}

func (teams *teams) createBaseBody(color, title, avatarDataURI, dateStr, projectVersion string) []map[string]any {
	body := []map[string]any{
		{
			"type":    "Container",
			"bleed":   true,
			"spacing": "None",
			"style":   color,
			"items": []map[string]any{
				{
					"type":   "TextBlock",
					"text":   title,
					"weight": "bolder",
					"size":   "medium",
					"color":  color,
				},
				teams.createAuthorSection(avatarDataURI, dateStr),
			},
		},
	}

	if facts := teams.createFactsSection(); facts != nil {
		body = append(body, facts)
	}

	return body
}

func (teams *teams) createAuthorSection(avatarDataURI, dateStr string) map[string]any {
	return map[string]any{
		"type": "ColumnSet",
		"columns": []map[string]any{
			{
				"type":  "Column",
				"width": "auto",
				"items": []map[string]any{
					{
						"type":  "Image",
						"url":   avatarDataURI,
						"size":  "small",
						"style": "Person",
					},
				},
			},
			{
				"type":  "Column",
				"width": "stretch",
				"items": []map[string]any{
					{
						"type":   "TextBlock",
						"text":   "@" + teams.setting.GetCommitAuthor(),
						"weight": "bolder",
						"wrap":   true,
					},
					{
						"type":     "TextBlock",
						"spacing":  "None",
						"text":     fmt.Sprintf("{{DATE(%s, SHORT)}} at {{TIME(%s)}}", dateStr, dateStr),
						"isSubtle": true,
						"wrap":     true,
					},
				},
			},
		},
	}
}

func (teams *teams) createFactsSection() map[string]any {
	// Define available facts
	allFacts := map[string]map[string]string{
		"project": {
			"title": "Project:",
			"value": teams.setting.GetRepo(),
		},
		"message": {
			"title": "Message:",
			"value": teams.setting.GetCommitMessage(),
		},
		"version": {
			"title": "Version:",
			"value": teams.setting.GetVersion(),
		},
	}

	// Get requested facts
	var facts []map[string]string
	requestedFacts := teams.setting.Facts
	if len(requestedFacts) == 0 {
		// If no facts specified, show all
		for _, fact := range allFacts {
			facts = append(facts, fact)
		}
	} else {
		// Show only requested facts
		for _, name := range requestedFacts {
			name = strings.TrimSpace(name)
			if fact, exists := allFacts[name]; exists {
				facts = append(facts, fact)
			}
		}
	}

	// Return nil if no facts to show
	if len(facts) == 0 {
		return nil
	}

	return map[string]any{
		"type": "Container",
		"items": []map[string]any{
			{
				"type":  "FactSet",
				"facts": facts,
			},
		},
	}
}

func (teams *teams) appendVariablesTable(body []map[string]any, variables []string) []map[string]any {
	body = append(body, map[string]any{
		"type":   "TextBlock",
		"text":   "Variables:",
		"weight": "bolder",
		"wrap":   true,
	})

	var rows []map[string]any
	for _, varName := range variables {
		varName = strings.TrimSpace(varName)
		rows = append(rows, createTableRow(varName, teams.setting.GetVariableValue(varName)))
	}

	body = append(body, map[string]any{
		"type": "Table",
		"columns": []map[string]any{
			{"width": 1},
			{"width": 2},
		},
		"spacing":           "Small",
		"showGridLines":     false,
		"firstRowAsHeaders": false,
		"rows":              rows,
	})

	return body
}

func createTableRow(name, value string) map[string]any {
	return map[string]any{
		"type": "TableRow",
		"cells": []map[string]any{
			createTableCell(name),
			createTableCell(value),
		},
		"style": "default",
	}
}

func createTableCell(text string) map[string]any {
	return map[string]any{
		"type": "TableCell",
		"items": []map[string]any{
			{
				"type":     "TextBlock",
				"text":     text,
				"wrap":     true,
				"weight":   "Default",
				"fontType": "Monospace",
			},
		},
	}
}

func (teams *teams) createCardActions() []map[string]any {
	// Define available actions
	allActions := map[string]any{
		"pipeline": map[string]any{
			"type":  "Action.OpenUrl",
			"title": "View Pipeline",
			"url":   teams.setting.GetPipelineURL(),
		},
	}

	// Add commit/release action
	actionURL := teams.setting.GetPipelineForgeURL()

	if tag := teams.setting.GetCommitTag(); tag != "" {
		actionURL = fmt.Sprintf("%s/releases/tag/%s", teams.setting.GetRepo(), tag)
		allActions["release"] = map[string]any{
			"type":  "Action.OpenUrl",
			"title": "View Release",
			"url":   actionURL,
		}
	} else {
		allActions["commit"] = map[string]any{
			"type":  "Action.OpenUrl",
			"title": "View Commit",
			"url":   actionURL,
		}
	}

	// Get requested buttons
	var actions []map[string]any
	requestedButtons := teams.setting.Buttons

	if len(requestedButtons) == 0 {
		// If no buttons specified, show all with pipeline first
		if pipeline, exists := allActions["pipeline"]; exists {
			actions = append(actions, pipeline.(map[string]any))
		}
		for name, action := range allActions {
			if name != "pipeline" {
				actions = append(actions, action.(map[string]any))
			}
		}
	} else {
		// Show buttons in the order specified in PLUGIN_BUTTONS
		for _, name := range requestedButtons {
			name = strings.TrimSpace(name)
			if action, exists := allActions[name]; exists {
				actions = append(actions, action.(map[string]any))
			}
		}
	}

	return actions
}
