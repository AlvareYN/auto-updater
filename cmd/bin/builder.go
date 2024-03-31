package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/AlvareYN/auto-updater/internal/updater"
	"golang.org/x/mod/semver"
)

type Build map[string]string
type Releases map[string]*JSONVersioning

type JSONVersioning struct {
	Version string `json:"version"`
	Build   Build  `json:"build"`
	Date    string `json:"date"`
}

type JSONVersioningPackage struct {
	ApplicationName string   `json:"application_name"`
	Current         string   `json:"current"`
	Description     string   `json:"description"`
	Date            string   `json:"date"`
	Build           Build    `json:"build"`
	Releases        Releases `json:"releases"`
}

func main() {
	version := updater.Version

	versioning, err := loadJSONVersioning()

	if err != nil {
		log.Fatal(err)
		return
	}

	if semver.Compare(version, versioning.Current) < 0 {
		log.Fatal("Current version is lower than the last version")
		return
	}

	log.Println("Building current version:", version)

	err = os.MkdirAll(fmt.Sprintf("./build/%s", version), fs.ModePerm)

	if err != nil {
		log.Fatal(err)
		return
	}

	fileName := fmt.Sprintf("%s-%s-%s", version, runtime.GOOS, runtime.GOARCH)
	buildFolder := fmt.Sprintf(`./build/%s/%s`, version, fileName)

	log.Println("Building for:", buildFolder)

	err = exec.Command("go", "build", "-o", buildFolder, "./main.go").Run()

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Build successful")

	value := updateJSONVersioning(versioning, version, fileName)

	if value == nil {
		log.Println("Version already exists in the past releases or the build is the same as the current version")
		return
	}

	err = writeJSONVersioning(value)

	if err != nil {
		log.Fatal(err)
		return
	}

	//zip file using archive/zip
	originalFile, err := os.Open(buildFolder)
	if err != nil {
		log.Fatal(err)
	}

	zipFile, err := os.Create(fmt.Sprintf("%s.zip", buildFolder))

	if err != nil {
		log.Fatal(err)
		return
	}

	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	defer zipWriter.Close()

	fileContent, err := zipWriter.Create(fileName)

	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(fileContent, originalFile); err != nil {
		log.Fatal(err)
	}

	log.Println("Zip file created")

}

func loadJSONVersioning() (*JSONVersioningPackage, error) {
	enviroment := os.Getenv("ENV")

	var jsonVersioning JSONVersioningPackage

	if enviroment == "" || enviroment == "development" {
		enviroment = "development"
		log.Println("Loading development enviroment")
		file := "./build/version.json"

		jsonFile, err := os.Open(file)

		if err != nil {
			return &JSONVersioningPackage{
				ApplicationName: updater.AppName,
				Current:         updater.Version,
				Description:     "auto-updater",
				Date:            time.Now().Format(time.RFC3339),
				Build:           Build{},
				Releases:        Releases{},
			}, nil
		}

		defer jsonFile.Close()

		jsonParser := json.NewDecoder(jsonFile)

		err = jsonParser.Decode(&jsonVersioning)

		if err != nil {
			return nil, fmt.Errorf("error parsing json file [%s]: %w", enviroment, err)
		}

		return &jsonVersioning, nil

	}
	//read from bucket
	return nil, nil
}

func writeJSONVersioning(jsonVersioning *JSONVersioningPackage) error {
	enviroment := os.Getenv("ENV")

	if enviroment == "" || enviroment == "development" {
		enviroment = "development"
		log.Println("Writing development enviroment")
		file := "./build/version.json"

		bytes, err := json.Marshal(jsonVersioning)

		if err != nil {
			return fmt.Errorf("error marshalling json file [%s]: %w", enviroment, err)
		}

		err = os.WriteFile(file, bytes, fs.ModePerm)

		if err != nil {
			return fmt.Errorf("error writing json file [%s]: %w", enviroment, err)
		}

		return nil

	}
	//write to bucket
	return nil
}

func updateJSONVersioning(versionMetadata *JSONVersioningPackage, version string, fileName string) *JSONVersioningPackage {
	log.Println("Updating version metadata")

	if versionMetadata.Current == version {
		if versionMetadata.Build[runtime.GOOS] == fileName {
			return nil
		}

		versionMetadata.Build[runtime.GOOS] = fileName
		versionMetadata.Date = time.Now().Format(time.RFC3339)

		return versionMetadata
	}

	if versionMetadata.Releases[version] != nil {
		return nil
	}

	os.Open("./build/version.json")

	copy := JSONVersioning{
		Version: versionMetadata.Current,
		Date:    versionMetadata.Date,
		Build:   versionMetadata.Build,
	}

	versionMetadata.Releases[copy.Version] = &copy

	versionMetadata.Current = version
	versionMetadata.Date = time.Now().Format(time.RFC3339)

	versionMetadata.Build = Build{}
	versionMetadata.Build[runtime.GOOS] = fileName

	return versionMetadata
}
