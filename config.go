package main

import (
	"fmt"
	"os"
)

var (
	configs      map[string]string
	validConfigs map[string]bool
)

func init() {
	validConfigs = map[string]bool{
		"APP_SECRET":    true,
		"DISCORD_TOKEN": true,
		"CHANNEL_ID":    true,
		"BOT_PORT":      false,
	}
}

func parseConfig() {
	configs = make(map[string]string)
	for configName, required := range validConfigs {
		if configValue, ok := os.LookupEnv(configName); ok && configValue != "" {
			configs[configName] = configValue
		} else if required {
			panic(fmt.Sprintf("%s environment variable is not set", configName))
		}
	}
}

func getConfigValue(key, def string) string {
	if len(configs) == 0 {
		return def
	}
	if val, ok := configs[key]; ok {
		return val
	}
	return def
}

func mustGetConfigValue(key string) string {
	if len(configs) == 0 {
		panic("Config not parsed")
	}
	if val, ok := configs[key]; ok {
		return val
	}
	panic("Unknown config value " + key)
}
