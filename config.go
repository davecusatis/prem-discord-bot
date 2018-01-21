package main

import (
	"os"
	"strings"
)

var (
	configs map[string]string
)

func parseConfig() {
	configs = make(map[string]string)
	for _, env := range os.Environ() {
		e := strings.Split(env, "=")
		configs[e[0]] = e[1]
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
	panic("Unknown config value")
}
