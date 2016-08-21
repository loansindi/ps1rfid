package cfg

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestMissingConfig(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	settingsFile := fmt.Sprint(r.Int63())
	_, err := ReadConfig(settingsFile)
	if err == nil {
		t.Error("Expected error for missing config file")
	}
}

func TestValidConfig(t *testing.T) {
	_, err := ReadConfig("../config.toml-example")
	if err != nil {
		t.Error("Unexpected error for valid settings file: %v", err)
	}
}

func TestInValidConfig(t *testing.T) {
	_, err := ReadConfig("../README.md")
	if err == nil {
		t.Error("Expected error for missing config file")
	}
}
