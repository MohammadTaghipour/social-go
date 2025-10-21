package env

import (
	"os"
	"strconv"
)

func GetString(key, defaul string) string {
	value, isOK := os.LookupEnv(key)
	if !isOK {
		return defaul
	}
	return value
}

func GetInt(key string, defaul int) int {
	value, isOK := os.LookupEnv(key)

	valAsInt, err := strconv.Atoi(value)

	if err != nil || !isOK {
		return defaul
	}

	return valAsInt
}

func GetBool(key string, defaul bool) bool {
	value, isOK := os.LookupEnv(key)

	valAsBool, err := strconv.ParseBool(value)

	if err != nil || !isOK {
		return defaul
	}

	return valAsBool
}
