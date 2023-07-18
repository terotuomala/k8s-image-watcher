package utils

import (
	"errors"
	"os"
	"strconv"
)

var ErrEnvVarEmpty = errors.New("environment variable empty")

func GetEnvStr(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return val, ErrEnvVarEmpty
	}
	return val, nil
}

func GetEnvBool(key string) (bool, error) {
	strVal, err := GetEnvStr(key)
	if err != nil {
		return false, err
	}
	boolVal, err := strconv.ParseBool(strVal)
	if err != nil {
		return false, err
	}
	return boolVal, nil
}
