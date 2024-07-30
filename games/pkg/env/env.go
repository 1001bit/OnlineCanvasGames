package env

import (
	"log"
	"os"
)

func GetEnvVal(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Println("couldn't find environment value:", key)
	}
	return value
}
