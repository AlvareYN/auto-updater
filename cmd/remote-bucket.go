package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func GetRemoteVersionMetadata() (map[string]interface{}, error) {
	enviroment := os.Getenv("ENV")

	if enviroment == "development" || enviroment == "" {
		file, err := os.ReadFile("./build/version.json")

		if err != nil {
			return nil, err
		}

		var versioning map[string]interface{}

		err = json.Unmarshal(file, &versioning)

		if err != nil {
			return nil, err
		}

		return versioning, nil
	}
	// return what bucket has
	return nil, nil
}

func FetchRemoteVersion(version string, fileName string) error {

	enviroment := os.Getenv("ENV")

	if enviroment == "development" || enviroment == "" {
		file, err := os.Open(fmt.Sprintf("./build/%s/%s", version, fileName))

		if err != nil {
			return fmt.Errorf("[%s] error opening file: %w", "FetchRemoteVersion", err)
		}

		defer file.Close()

		newFile, err := os.Create(fileName)

		if err != nil {
			return fmt.Errorf("[%s] error creating new file: %w", "FetchRemoteVersion", err)
		}

		defer newFile.Close()

		if _, err := io.Copy(newFile, file); err != nil {
			return fmt.Errorf("[%s] error copying file: %w", "FetchRemoteVersion", err)
		}

		return nil
	}

	// fetch from bucket
	return nil
}
