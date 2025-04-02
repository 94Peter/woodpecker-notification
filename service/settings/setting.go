package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

var defaultUsedWoodEnvs = []string{
	"CI_COMMIT_TAG",
	"CI_COMMIT_SHA",
	"CI_REPO",
	"CI_COMMIT_AUTHOR",
	"CI_COMMIT_AUTHOR_AVATAR",
	"CI_PREV_PIPELINE_STATUS",
	"CI_COMMIT_MESSAGE",
	"CI_COMMIT_REF",
	"CI_PREV_COMMIT_URL",
	"CI_PIPELINE_FORGE_URL",
	"CI_PIPELINE_URL",
}

type Settings struct {
	Webhooks     []*WebhookProvider `wp_env:"PLUGIN_WEBHOOKS"`
	UsedWoodEnvs []string           `wp_env:"PLUGIN_WOODPETER_ENVS"`
	Debug        bool               `wp_env:"PLUGIN_DEBUG"`
}

type WebhookProvider struct {
	Webhook  string           `json:"webhook"`
	Provider string           `json:"provider"`
	Configs  *json.RawMessage `json:"configs"`

	bindData map[string]string
}

func GetSettings() (*Settings, error) {
	settings := &Settings{}
	rv := reflect.ValueOf(settings).Elem()
	rt := rv.Type()
	for i := range rt.NumField() {
		field := rt.Field(i)
		envVarName := field.Tag.Get("wp_env")
		envVarValue := os.Getenv(envVarName)
		if envVarValue != "" {
			fv := rv.Field(i)
			if fv.CanSet() {
				switch fv.Kind() {
				case reflect.String:
					fv.SetString(envVarValue)
				case reflect.Slice:
					switch fv.Type().Elem().Kind() {
					case reflect.Ptr:
						// []*WebhookProvider
						fmt.Println("slice", envVarName)
						var webhooks []*WebhookProvider
						err := json.Unmarshal([]byte(envVarValue), &webhooks)
						if err != nil {
							return nil, err
						}
						fv.Set(reflect.ValueOf(webhooks))
					case reflect.String:
						// []string
						var values []string
						err := json.Unmarshal([]byte(envVarValue), &values)
						if err != nil {
							return nil, err
						}
						fv.Set(reflect.ValueOf(values))
					}
				case reflect.Bool:
					value, err := strconv.ParseBool(envVarValue)
					if err != nil {
						return nil, err
					}
					fv.SetBool(value)
				}
			}
		}
	}
	if len(settings.UsedWoodEnvs) == 0 {
		settings.UsedWoodEnvs = defaultUsedWoodEnvs
	}
	bindData := settings.genBindData()
	if settings.Debug {
		printBuildInfo(bindData)
	}
	if err := settings.validate(); err != nil {
		return nil, err
	}
	for _, hook := range settings.Webhooks {
		hook.bindData = bindData
	}
	return settings, nil
}

func (settings *Settings) validate() error {
	if len(settings.Webhooks) == 0 {
		return fmt.Errorf("need webhook")
	}
	return nil
}

func (settings *Settings) genBindData() map[string]string {
	result := map[string]string{}
	for _, key := range settings.UsedWoodEnvs {
		result[key] = os.Getenv(key)
	}
	var tag string
	if tag = result["CI_COMMIT_TAG"]; tag == "" {
		tag = result["CI_COMMIT_SHA"]
	}
	result["BUILD_TAG"] = tag
	return result
}

func printBuildInfo(bindData map[string]string) {
	fmt.Println("\nBuild Info:")
	fmt.Printf(" PROJECT: %s\n", bindData["CI_REPO"])
	fmt.Printf(" VERSION: %s\n", bindData["BUILD_TAG"])
	fmt.Printf(" STATUS:  %s\n", bindData["CI_PREV_PIPELINE_STATUS"])
	fmt.Printf(" DATE:    %s\n", time.Now().UTC().Format(time.RFC3339))
}
