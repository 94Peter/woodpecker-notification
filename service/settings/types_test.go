package settings

import (
	"os"
	"testing"
)

func TestGetSettings(t *testing.T) {
	// Set environment variables
	webhooksJSON := `[{"webhook":"https://example.com","provider":"example","configs":{}}]`
	os.Setenv("PLUGIN_WEBHOOKS", webhooksJSON)
	os.Setenv("PLUGIN_WOODPETER_ENVS", "[\"env1\",\"env2\"]")
	os.Setenv("PLUGIN_DEBUG", "true")

	// Get settings
	settings, err := GetSettings()
	if err != nil {
		t.Fatal(err)
	}

	// Check settings values
	if len(settings.Webhooks) != 1 {
		t.Errorf("expected 1 webhook, got %d", len(settings.Webhooks))
	}
	if settings.Webhooks[0].Webhook != "https://example.com" {
		t.Errorf("expected webhook URL to be https://example.com, got %s", settings.Webhooks[0].Webhook)
	}
	if settings.Webhooks[0].Provider != "example" {
		t.Errorf("expected provider to be example, got %s", settings.Webhooks[0].Provider)
	}
	if len(settings.UsedWoodEnvs) != 2 {
		t.Errorf("expected 2 used wood envs, got %d", len(settings.UsedWoodEnvs))
	}
	if settings.UsedWoodEnvs[0] != "env1" || settings.UsedWoodEnvs[1] != "env2" {
		t.Errorf("expected used wood envs to be [env1, env2], got %v", settings.UsedWoodEnvs)
	}
	if !settings.Debug {
		t.Errorf("expected debug to be true, got false")
	}
}

func TestGetSettings_InvalidJSON(t *testing.T) {
	// Set environment variable with invalid JSON
	os.Setenv("PLUGIN_WEBHOOKS", " invalid json ")

	// Get settings
	settings, err := GetSettings()
	if err == nil {
		t.Errorf("expected error parsing JSON, got nil")
	}
	if settings != nil {
		t.Errorf("expected settings to be nil, got %+v", settings)
	}
}

func TestGetSettings_TwoWebhooks(t *testing.T) {
	// Set environment variable with two webhooks JSON
	webhooksJSON := `[{"webhook":"https://example.com","provider":"example","configs":{}},{"webhook":"https://example2.com","provider":"example2","configs":{}}]`
	os.Setenv("PLUGIN_WEBHOOKS", webhooksJSON)
	os.Setenv("PLUGIN_WOODPETER_ENVS", "[\"env1\",\"env2\"]")
	os.Setenv("PLUGIN_DEBUG", "true")

	// Get settings
	settings, err := GetSettings()
	if err != nil {
		t.Fatal(err)
	}

	// Check settings values
	if len(settings.Webhooks) != 2 {
		t.Errorf("expected 2 webhooks, got %d", len(settings.Webhooks))
	}
	if settings.Webhooks[0].Webhook != "https://example.com" {
		t.Errorf("expected webhook URL to be https://example.com, got %s", settings.Webhooks[0].Webhook)
	}
	if settings.Webhooks[0].Provider != "example" {
		t.Errorf("expected provider to be example, got %s", settings.Webhooks[0].Provider)
	}
	if settings.Webhooks[1].Webhook != "https://example2.com" {
		t.Errorf("expected webhook URL to be https://example2.com, got %s", settings.Webhooks[1].Webhook)
	}
	if settings.Webhooks[1].Provider != "example2" {
		t.Errorf("expected provider to be example2, got %s", settings.Webhooks[1].Provider)
	}
	if len(settings.UsedWoodEnvs) != 2 {
		t.Errorf("expected 2 used wood envs, got %d", len(settings.UsedWoodEnvs))
	}
	if settings.UsedWoodEnvs[0] != "env1" || settings.UsedWoodEnvs[1] != "env2" {
		t.Errorf("expected used wood envs to be [env1, env2], got %v", settings.UsedWoodEnvs)
	}
	if !settings.Debug {
		t.Errorf("expected debug to be true, got false")
	}
}
