package utils

import (
	"os"
)

// getEnv get key environment variable if exist, otherwise return defaultValue
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

var INCLUSTER_ANOMALY_IMAGE_VERSION = getEnv("INCLUSTER_ANOMALY_IMAGE_VERSION", "latest")
var INCLUSTER_ANOMALY_IMAGE = "quay.io/openshiftanalytics/incluster-anomaly:" + INCLUSTER_ANOMALY_IMAGE_VERSION
