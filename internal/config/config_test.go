package config

import (
	"testing"
)

func TestConfigPath(t *testing.T) {
    actual, err := getConfigPath()
    if err != nil {
        t.Errorf("Error getting config path: %v\n", err)
    }

    expect:= "/home/aminkafri/.gatorconfig.json"
    if expect != actual {
        t.Errorf("Expect %v, Got %v\n", expect, actual)
    }
}
