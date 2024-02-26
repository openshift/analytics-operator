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

var InclusterAnomalyImageVersion string
var InclusterAnomalyImage string

func init() {
	InclusterAnomalyImageVersion = getEnv("InclusterAnomalyImageVersion", "latest")
	InclusterAnomalyImage = "quay.io/openshiftanalytics/incluster-anomaly:" + InclusterAnomalyImageVersion
}
