package config

import (
	"testing"
)

func TestShouldLoadConfigCorrectly(t *testing.T) {
	expectedServerMode := "staging"
	expectedStagingKey := "test_gshuPaZoeEG6ovbc8M79w0QyM"
	expectedLiveKey := "liveKey"

	configFilePath = "./testdata/config_test.json"
	loadConfig()

	actualServerMode := config["server"]["mode"]
	actualStagingKey := config["staging"]["accessKey"]
	actualLiveKey := config["live"]["accessKey"]

	if expectedServerMode != actualServerMode {
		t.Fatalf("Expected server mode to be %s but got %s", expectedServerMode, actualServerMode)
	}

	if expectedStagingKey != actualStagingKey {
		t.Fatalf("Expected staging apikey to be %s but got %s", expectedStagingKey, actualStagingKey)
	}

	if expectedLiveKey != actualLiveKey {
		t.Fatalf("Expected live aipkey to be %s but got %s", expectedLiveKey, actualLiveKey)
	}
}
