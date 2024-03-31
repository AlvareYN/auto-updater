package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigStructure struct {
	ServerPort string
	Bucket     string
}

var Config ConfigStructure

func ConfigInit() {
	godotenv.Load(".env")
	configString := os.Getenv("CONFIG")

	err := json.Unmarshal([]byte(configString), &Config)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Config loaded")
}
