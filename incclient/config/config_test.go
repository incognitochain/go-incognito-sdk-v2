package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	filePath := "config_example.json"

	cfg, err := LoadConfig(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)

	jsb, _ := json.MarshalIndent(cfg, "", "\t")
	fmt.Println(string(jsb))
}

func TestSaveConfig(t *testing.T) {
	filePath := "config_example.json"
	filePath2 := "config_example2.json"

	cfg, err := LoadConfig(filePath)
	if err != nil {
		panic(err)
	}

	err = SaveConfig(*cfg, filePath2)
	if err != nil {
		panic(err)
	}
}
