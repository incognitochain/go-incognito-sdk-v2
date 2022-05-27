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

	jsb, _ := json.MarshalIndent(cfg, "", "\t")
	fmt.Println(string(jsb))
}
