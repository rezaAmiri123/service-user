package utils

import (
	"strconv"
)

// Get config path for local or docker
func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./cmd/config/config-local"
}

func UintToString(n uint) string{
	return strconv.FormatUint(uint64(n), 10)

}