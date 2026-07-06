package env

import (
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func Parse(config any) error {
	element, elementType, err := extractConfig(config)
	if err != nil {
		return err
	}

	for index := range element.NumField() {
		field, envKey, defaultVal, ok := extractFieldEnvInfo(element, elementType, index)
		if !ok {
			continue
		}

		setFieldFromEnv(*field, envKey, defaultVal)
	}

	return nil
}
